//go:build windows
// +build windows

package config

import (
	"github.com/AllenDang/w32"
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
		if processID == myProcessID {
			w32.ShowWindow(console, w32.SW_HIDE)
		}
	}
}
