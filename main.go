package main

import (
	"os"
	"os/exec"
	"os/signal"

	"github.com/gin-gonic/gin"
)

func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		router.GET("/", func(c *gin.Context) {
			c.Writer.Write([]byte("Hello World"))
		})
		router.Run(":8080")
	}()

	_ = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	macEdgePath := "/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge"
	cmd := exec.Command(macEdgePath, "--app=http://127.0.0.1:8080/")
	cmd.Start()

	chSignal := make(chan os.Signal, 1)
	// 获取 Interrupt 信号到chan中
	signal.Notify(chSignal, os.Interrupt)

	// 阻塞等待 Interrupt 信号
	<-chSignal
	cmd.Process.Kill()
}
