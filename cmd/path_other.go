//go:build !windows

package cmd

import (
	"fmt"
	"os"
)

// On non-Windows, PATH lives in shell profiles, which vary too much to edit
// safely for the user. Print the directory and manual instructions instead.

func pathStatus() int {
	dir, _ := exeDir()
	fmt.Printf("csvcolor directory:\n%s\n", dir)
	fmt.Println("Add it to PATH in your shell profile (see 'csvcolor path add').")
	return 0
}

func pathAdd() int {
	dir, err := exeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	fmt.Printf("Automatic PATH setup is Windows-only. Add this line to your\n"+
		"shell profile (~/.bashrc, ~/.zshrc, or ~/.config/fish/config.fish):\n\n"+
		"  export PATH=\"%s:$PATH\"\n\nThen restart your shell.\n", dir)
	return 0
}

func pathRemove() int {
	fmt.Println("Automatic PATH setup is Windows-only. Remove the csvcolor line")
	fmt.Println("you added from your shell profile.")
	return 0
}
