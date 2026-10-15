//go:build !windows

package main

// attachConsole is a no-op on non-Windows platforms, where a single binary
// serves both the terminal and a windowed launch without console tricks.
func attachConsole() {}
