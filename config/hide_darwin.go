//go:build darwin
// +build darwin

package config

import (
	"os"
	"os/exec"
	"strconv"
)

func init() {
	showConsoleAsync()
}

func showConsoleAsync() {
	cmd := exec.Command("osascript", "-e", `tell application "System Events" to set visible of processes whose unix id is `+strconv.Itoa(os.Getpid())+` to false`)
	cmd.Start()
}
