package main

import (
	"github.com/808Mak1r/Beacon/server"
	"os"
	"os/exec"
	"os/signal"
)

func startBrowser(port string) (cmd *exec.Cmd) {
	macChromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	macEdgePath := "/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge"
	_ = macChromePath
	cmd = exec.Command(macEdgePath, "--app=http://127.0.0.1:"+port+"/static/index.html")
	cmd.Start()
	return
}

func listenToInterrupt() (chSignal chan os.Signal) {
	chSignal = make(chan os.Signal, 1)
	// 获取 Interrupt 信号到chan中
	signal.Notify(chSignal, os.Interrupt)
	return
}

func main() {
	port := "27149"
	go server.Run(port)

	cmd := startBrowser(port)

	chSignal := listenToInterrupt()

	// 阻塞等待 Interrupt 信号
	<-chSignal
	cmd.Process.Kill()
}
