//go:build !windows && !darwin && !linux

package ui

import "fmt"

// Run is a placeholder on platforms without a webview host (anything other than
// Windows, macOS, and Linux). The CLI still works everywhere.
func Run() {
	fmt.Println("csvcolor: the GUI is not available on this platform. Use the CLI, e.g.:")
	fmt.Println(`  csvcolor parse color "Red" file.csv`)
	fmt.Println("Run 'csvcolor --help' for details.")
}
