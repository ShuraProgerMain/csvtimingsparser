// Package clip copies text to the OS clipboard. It wraps atotto/clipboard,
// which uses native Win32 APIs on Windows (correct Unicode), and pbcopy /
// xsel / xclip on macOS and Linux. No cgo, keeping the binary lightweight.
package clip

import "github.com/atotto/clipboard"

// Copy writes text to the system clipboard.
func Copy(text string) error {
	return clipboard.WriteAll(text)
}
