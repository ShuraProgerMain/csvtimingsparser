//go:build windows

package main

import (
	"os"
	"syscall"
)

// attachConsole reattaches a GUI-subsystem binary (built with -H windowsgui) to
// the console of the parent process when run from a terminal, so CLI output is
// visible. When launched from Explorer there is no parent console and this is a
// no-op, which is exactly what keeps the GUI free of a flashing console window.
func attachConsole() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	attach := kernel32.NewProc("AttachConsole")
	const attachParentProcess = ^uintptr(0) // (DWORD)-1
	if r, _, _ := attach.Call(attachParentProcess); r == 0 {
		return // no parent console available
	}
	reopen := func(stdHandle int, target **os.File, name string) {
		h, err := syscall.GetStdHandle(stdHandle)
		if err != nil || h == 0 || h == syscall.InvalidHandle {
			return
		}
		*target = os.NewFile(uintptr(h), name)
	}
	reopen(syscall.STD_OUTPUT_HANDLE, &os.Stdout, "stdout")
	reopen(syscall.STD_ERROR_HANDLE, &os.Stderr, "stderr")
}
