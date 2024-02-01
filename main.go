package main

import (
	"github.com/808Mak1r/Beacon/config"
	"github.com/808Mak1r/Beacon/server"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
)

func startBrowser(port string) (cmd *exec.Cmd) {
	cmd = exec.Command(getBrowserPath(), "--app=http://127.0.0.1:"+port+"/static/index.html")
	cmd.Start()
	return
}

func listenToInterrupt() (chSignal chan os.Signal) {
	chSignal = make(chan os.Signal, 1)
	// 获取 Interrupt 信号到chan中
	signal.Notify(chSignal, os.Interrupt)
	return
}

func getBrowserPath() (browserPath string) {
	switch runtime.GOOS {
	case "darwin":
		browserPath = "/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge"
		if _, err := os.Stat(browserPath); err != nil {
			browserPath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
		}
	case "windows":
		browserPath = "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
		if _, err := os.Stat(browserPath); err != nil {
			browserPath = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe"
		}
	}
	return
}

func main() {
	port := config.GetPort()
	go server.Run(port)

	cmd := startBrowser(port)

	chSignal := listenToInterrupt()
	chBrowserDie := make(chan error)
	go func() {
		chBrowserDie <- cmd.Wait()
	}()

	select {
	// 阻塞等待 Interrupt 信号
	case <-chSignal:
		cmd.Process.Kill()
	case <-chBrowserDie:
		os.Exit(0)
	}
}
