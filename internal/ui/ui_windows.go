//go:build windows

package ui

import (
	"os"
	"syscall"
	"unsafe"

	webview "github.com/jchv/go-webview2"
)

// Win32 plumbing for making the window frameless, rounded, and draggable.
var (
	user32               = syscall.NewLazyDLL("user32.dll")
	dwmapi               = syscall.NewLazyDLL("dwmapi.dll")
	procSetWindowLongPtr = user32.NewProc("SetWindowLongPtrW")
	procSetWindowPos     = user32.NewProc("SetWindowPos")
	procReleaseCapture   = user32.NewProc("ReleaseCapture")
	procSendMessageW     = user32.NewProc("SendMessageW")
	procShowWindow       = user32.NewProc("ShowWindow")
	procSetWindowAttr    = dwmapi.NewProc("DwmSetWindowAttribute")
)

const (
	gwlStyle      = ^uintptr(15) // -16
	gwlExStyle    = ^uintptr(19) // -20
	wsPopup       = 0x80000000
	wsVisible     = 0x10000000
	wsMinimizeBox = 0x00020000
	wsMaximizeBox = 0x00010000
	wsThickFrame  = 0x00040000 // resizable edges (no caption = frameless + resizable)
	wsExAppWindow = 0x00040000 // keep a frameless window in the taskbar

	swpNoMove      = 0x0002
	swpNoSize      = 0x0001
	swpNoZOrder    = 0x0004
	swpFrameChange = 0x0020

	wmNCLButtonDown = 0x00A1
	htCaption       = 2
	swMinimize      = 6

	dwmwaCornerPref = 33 // DWMWA_WINDOW_CORNER_PREFERENCE
	dwmwcpRound     = 2  // DWMWCP_ROUND
)

// Run creates the frameless WebView2 window and blocks until it closes.
func Run() {
	debug := os.Getenv("CSVCOLOR_DEBUG") != ""
	w := webview.NewWithOptions(webview.WebViewOptions{
		Debug:     debug,
		AutoFocus: true,
		WindowOptions: webview.WindowOptions{
			Title:  "CSV Color",
			Width:  460,
			Height: 620,
			Center: true,
		},
	})
	if w == nil {
		return
	}
	defer w.Destroy()

	h := &host{}
	hwnd := uintptr(w.Window())
	makeFrameless(hwnd)
	// Allow resizing but not below a size where the layout still reads well.
	w.SetSize(380, 480, webview.HintMin)

	w.Init("window.CSVCOLOR_FRAMELESS=true;") // page draws its own title bar
	w.Bind("hostLoadCSV", h.loadCSV)
	w.Bind("hostDataFor", h.dataFor)
	w.Bind("hostEditNote", h.editNote)
	w.Bind("hostCopy", h.copy)
	w.Bind("hostDrag", func() { startDrag(hwnd) })
	w.Bind("hostMinimize", func() { procShowWindow.Call(hwnd, swMinimize) })
	w.Bind("hostClose", func() { w.Terminate() })

	w.SetHtml(indexHTML)
	w.Run()
}

func makeFrameless(hwnd uintptr) {
	procSetWindowLongPtr.Call(hwnd, gwlExStyle, uintptr(wsExAppWindow))
	procSetWindowLongPtr.Call(hwnd, gwlStyle, uintptr(wsPopup|wsVisible|wsMinimizeBox|wsMaximizeBox|wsThickFrame))
	procSetWindowPos.Call(hwnd, 0, 0, 0, 0, 0, swpNoMove|swpNoSize|swpNoZOrder|swpFrameChange)
	// Win11 rounded window corners.
	pref := int32(dwmwcpRound)
	procSetWindowAttr.Call(hwnd, dwmwaCornerPref, uintptr(unsafe.Pointer(&pref)), unsafe.Sizeof(pref))
}

func startDrag(hwnd uintptr) {
	procReleaseCapture.Call()
	procSendMessageW.Call(hwnd, wmNCLButtonDown, htCaption, 0)
}
