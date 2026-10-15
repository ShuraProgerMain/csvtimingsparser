package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"csvcolor/internal/clip"
	"csvcolor/internal/csvparse"
)

// Usage text shown for -h/--help and on misuse.
const usage = `csvcolor — extract Source In / Color / Notes from a CSV.

GUI:
  csvcolor                 launch the window (drag a CSV onto it)

CLI:
  csvcolor parse color "<Color>" <file.csv>
                           print matching rows as "HH:MM:SS - Notes" and copy
                           them to the clipboard
  csvcolor colors <file.csv>
                           list the distinct colors found in the file
  csvcolor parse <file.csv>
                           print all rows (every color)

  csvcolor path status     show whether csvcolor's folder is on your PATH
  csvcolor path add        add csvcolor's folder to your user PATH (no admin),
                           so you can run it from any terminal
  csvcolor path remove     undo 'path add'

Flags:
  --no-copy                do not touch the clipboard (CLI only)
  -h, --help               show this help
`

// RunCLI executes the command-line interface. It returns a process exit code.
// args is os.Args[1:].
func RunCLI(args []string) int {
	noCopy := false
	var rest []string
	for _, a := range args {
		switch a {
		case "--no-copy":
			noCopy = true
		case "-h", "--help", "help":
			fmt.Print(usage)
			return 0
		default:
			rest = append(rest, a)
		}
	}

	if len(rest) == 0 {
		fmt.Fprint(os.Stderr, usage)
		return 2
	}

	switch strings.ToLower(rest[0]) {
	case "colors":
		return cmdColors(rest[1:])
	case "parse":
		return cmdParse(rest[1:], noCopy)
	case "path":
		return cmdPath(rest[1:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", rest[0], usage)
		return 2
	}
}

// cmdPath dispatches 'path status|add|remove'; the implementations are
// per-platform (path_windows.go / path_other.go).
func cmdPath(args []string) int {
	sub := "status"
	if len(args) > 0 {
		sub = strings.ToLower(args[0])
	}
	switch sub {
	case "status":
		return pathStatus()
	case "add", "install":
		return pathAdd()
	case "remove", "rm", "uninstall":
		return pathRemove()
	default:
		fmt.Fprintf(os.Stderr, "usage: csvcolor path [status|add|remove]\n")
		return 2
	}
}

// exeDir returns the directory containing the running executable, with symlinks
// resolved.
func exeDir() (string, error) {
	p, err := os.Executable()
	if err != nil {
		return "", err
	}
	if resolved, err := filepath.EvalSymlinks(p); err == nil {
		p = resolved
	}
	return filepath.Dir(p), nil
}

func cmdColors(args []string) int {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: csvcolor colors <file.csv>\n")
		return 2
	}
	res, err := csvparse.Parse(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	for _, c := range res.Colors {
		fmt.Println(c)
	}
	return 0
}

// cmdParse handles both:
//
//	parse color "<Color>" <file.csv>
//	parse <file.csv>            (all rows)
func cmdParse(args []string, noCopy bool) int {
	var path, color string
	switch {
	case len(args) == 1:
		path = args[0]
	case len(args) == 3 && strings.ToLower(args[0]) == "color":
		color = args[1]
		path = args[2]
	default:
		fmt.Fprintf(os.Stderr, "usage: csvcolor parse color \"<Color>\" <file.csv>\n")
		return 2
	}

	res, err := csvparse.Parse(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	rows := res.Rows
	if color != "" {
		rows = res.RowsForColor(color)
		if len(rows) == 0 {
			fmt.Fprintf(os.Stderr, "error: no rows with color %q (available: %s)\n",
				color, strings.Join(res.Colors, ", "))
			return 1
		}
	}

	out := csvparse.FormatList(rows)
	fmt.Print(out)

	if !noCopy {
		if err := clip.Copy(out); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not copy to clipboard: %v\n", err)
			// Not a hard failure — the data was still printed.
		}
	}
	return 0
}
