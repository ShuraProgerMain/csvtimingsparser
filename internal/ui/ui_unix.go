//go:build darwin || linux

package ui

import (
	"os"
	"runtime"

	webview "github.com/webview/webview_go"
)

// Run creates the GUI on macOS and Linux using webview_go (WKWebView / WebKitGTK
// via cgo). Unlike the Windows build, the window keeps its native title bar and
// controls, so the HTML hides its own custom title bar (CSVCOLOR_FRAMELESS=false).
//
// Requires cgo and the platform WebKit (macOS: built in; Linux: a WebKitGTK dev
// package, e.g. libwebkit2gtk-4.1-dev). The CLI needs none of this.
func Run() {
	runtime.LockOSThread() // webview must run on the main OS thread

	debug := os.Getenv("CSVCOLOR_DEBUG") != ""
	w := webview.New(debug)
	if w == nil {
		return
	}
	defer w.Destroy()

	w.SetTitle("CSV Color")
	w.SetSize(460, 620, webview.HintNone)
	w.SetSize(380, 480, webview.HintMin)

	h := &host{}
	// On these platforms the OS draws the window chrome, so the page hides its
	// own title bar; the window-control bindings are harmless no-ops (close is
	// wired up in case a future UI needs it).
	w.Init("window.CSVCOLOR_FRAMELESS=false;")
	w.Bind("hostLoadCSV", h.loadCSV)
	w.Bind("hostDataFor", h.dataFor)
	w.Bind("hostEditNote", h.editNote)
	w.Bind("hostCopy", h.copy)
	w.Bind("hostDrag", func() {})
	w.Bind("hostMinimize", func() {})
	w.Bind("hostClose", func() { w.Terminate() })

	w.SetHtml(indexHTML)
	w.Run()
}
