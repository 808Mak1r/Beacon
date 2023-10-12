package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.StaticFS("/static", http.FS(staticFiles))
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/static/") {
				reader, err := staticFiles.Open("index.html")
				if err != nil {
					log.Fatal(err)
				}
				defer reader.Close()
				stat, err := reader.Stat()
				if err != nil {
					log.Fatal(err)
				}
				c.DataFromReader(http.StatusOK, stat.Size(), "text/html;charset=utf-8", reader, nil)
			} else {
				c.Status(http.StatusNotFound)
			}
		})
		router.Run(":8080")
	}()

	macChromePath := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	macEdgePath := "/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge"
	_ = macChromePath
	cmd := exec.Command(macEdgePath, "--app=http://127.0.0.1:8080/static/index.html")
	cmd.Start()

	chSignal := make(chan os.Signal, 1)
	// 获取 Interrupt 信号到chan中
	signal.Notify(chSignal, os.Interrupt)

	// 阻塞等待 Interrupt 信号
	<-chSignal
	cmd.Process.Kill()
}
