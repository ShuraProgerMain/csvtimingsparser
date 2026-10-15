//go:build windows

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

// userPath reads the current user's PATH from the registry, unexpanded, along
// with its value type (so we can write it back in the same form).
func userPath() (raw string, valtype uint32, err error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	if err != nil {
		return "", registry.EXPAND_SZ, err
	}
	defer k.Close()
	raw, valtype, err = k.GetStringValue("Path")
	if err == registry.ErrNotExist {
		return "", registry.EXPAND_SZ, nil
	}
	return raw, valtype, err
}

func pathEntries(raw string) []string {
	var out []string
	for _, e := range strings.Split(raw, ";") {
		if strings.TrimSpace(e) != "" {
			out = append(out, e)
		}
	}
	return out
}

func sameDir(a, b string) bool {
	norm := func(s string) string {
		return strings.ToLower(strings.TrimRight(filepath.Clean(strings.TrimSpace(s)), `\`))
	}
	return norm(a) == norm(b)
}

func onUserPath(dir, raw string) bool {
	for _, e := range pathEntries(raw) {
		if sameDir(e, dir) {
			return true
		}
	}
	return false
}

func writeUserPath(value string, valtype uint32) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()
	if valtype == registry.SZ {
		return k.SetStringValue("Path", value)
	}
	return k.SetExpandStringValue("Path", value)
}

// broadcastEnv notifies running apps that the environment changed, so Explorer
// (and shells it spawns afterwards) pick up the new PATH. Already-open terminals
// still need to be reopened.
func broadcastEnv() {
	user32 := syscall.NewLazyDLL("user32.dll")
	send := user32.NewProc("SendMessageTimeoutW")
	const hwndBroadcast = 0xffff
	const wmSettingChange = 0x001A
	const smtoAbortIfHung = 0x0002
	env, _ := syscall.UTF16PtrFromString("Environment")
	var res uintptr
	send.Call(hwndBroadcast, wmSettingChange, 0,
		uintptr(unsafe.Pointer(env)), smtoAbortIfHung, 5000, uintptr(unsafe.Pointer(&res)))
}

func pathStatus() int {
	dir, err := exeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	raw, _, err := userPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading PATH: %v\n", err)
		return 1
	}
	if onUserPath(dir, raw) {
		fmt.Printf("on PATH: yes\n%s\n", dir)
	} else {
		fmt.Printf("on PATH: no\n%s\nRun 'csvcolor path add' to add it.\n", dir)
	}
	return 0
}

func pathAdd() int {
	dir, err := exeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	raw, valtype, err := userPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading PATH: %v\n", err)
		return 1
	}
	if onUserPath(dir, raw) {
		fmt.Printf("Already on PATH:\n%s\n", dir)
		return 0
	}
	newValue := dir
	if raw != "" {
		newValue = strings.TrimRight(raw, ";") + ";" + dir
	}
	if err := writeUserPath(newValue, valtype); err != nil {
		fmt.Fprintf(os.Stderr, "error updating PATH: %v\n", err)
		return 1
	}
	broadcastEnv()
	fmt.Printf("Added to your user PATH:\n%s\n\nOpen a NEW terminal, then run 'csvcolor' from anywhere.\n", dir)
	return 0
}

func pathRemove() int {
	dir, err := exeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	raw, valtype, err := userPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading PATH: %v\n", err)
		return 1
	}
	if !onUserPath(dir, raw) {
		fmt.Printf("Not on PATH; nothing to remove:\n%s\n", dir)
		return 0
	}
	var kept []string
	for _, e := range pathEntries(raw) {
		if !sameDir(e, dir) {
			kept = append(kept, e)
		}
	}
	if err := writeUserPath(strings.Join(kept, ";"), valtype); err != nil {
		fmt.Fprintf(os.Stderr, "error updating PATH: %v\n", err)
		return 1
	}
	broadcastEnv()
	fmt.Printf("Removed from your user PATH:\n%s\n\nOpen a new terminal for the change to take effect.\n", dir)
	return 0
}
