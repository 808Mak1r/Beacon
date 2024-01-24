package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
)

//go:embed frontend/dist/*
var FS embed.FS

// TextsController 文本上传
func TextsController(c *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		exe, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exe)
		if err != nil {
			log.Fatal(err)
		}
		filename := uuid.New().String()
		uploads := filepath.Join(dir, "uploads")
		err = os.MkdirAll(uploads, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		fullpath := path.Join("uploads", filename+".txt")
		err = os.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
	}
}

func GetUploadsDir() (uploads string) {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dir := filepath.Dir(exe)
	uploads = filepath.Join(dir, "uploads")

	return
}

func UploadsController(c *gin.Context) {
	if p := c.Param("path"); p != "" {
		target := filepath.Join(GetUploadsDir(), p)
		c.Header("Content-Disposition", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+p)
		c.Header("Content-Type", "application/octet-stream")
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func AddressesController(c *gin.Context) {
	addrs, _ := net.InterfaceAddrs()
	var result []string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}
	c.JSON(http.StatusOK, map[string]interface{}{"addresses": result})
}

func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.GET("/api/v1/downloads/:path", UploadsController)
		router.POST("/api/v1/texts", TextsController)
		router.GET("/api/v1/addresses", AddressesController)
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
