// csvcolor extracts the Source In / Color / Notes columns from a CSV.
//
// Run with no arguments to open the GUI (drag a CSV onto the window).
// Run with arguments to use the CLI, e.g.:
//
//	csvcolor parse color "Red" myfile.csv
//
// See cmd/cli.go for the full command reference.
package main

import (
	"os"

	"csvcolor/cmd"
	"csvcolor/internal/ui"
)

func main() {
	if len(os.Args) > 1 {
		attachConsole() // Windows: reattach to the parent terminal for CLI output
		os.Exit(cmd.RunCLI(os.Args[1:]))
	}
	ui.Run()
}
