package server

import (
	"embed"
	c "github.com/808Mak1r/Beacon/server/controller"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed frontend/dist/*
var FS embed.FS

func Run(port string) {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	api := router.Group("/api/v1")
	{
		api.GET("/downloads/:path", c.UploadsController)
		api.GET("/qrcodes", c.QrcodesController)
		api.GET("/addresses", c.AddressesController)
		api.POST("/texts", c.TextsController)
		api.POST("/files", c.FilesController)
	}

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
	router.Run(":" + port)
}
