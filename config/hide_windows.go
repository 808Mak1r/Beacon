//go:build windows
// +build windows

package config

import (
	w32 "github.com/gonutz/w32/v2"
	"os"
)

func init() {
	showConsoleAsync()
}

func showConsoleAsync() {
	console := w32.GetConsoleWindow()
	if console != 0 {
		_, processID := w32.GetWindowThreadProcessId(console)
		myProcessID := os.Getpid()
		if int(processID) == myProcessID {
			w32.ShowWindow(console, w32.SW_HIDE)
		}
	}
}
