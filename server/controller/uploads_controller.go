package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func getUploadsDir() (uploads string) {
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
		target := filepath.Join(getUploadsDir(), p)
		c.Header("Content-Disposition", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+p)
		c.Header("Content-Type", "application/octet-stream")
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}
