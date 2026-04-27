package stream

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/gin-gonic/gin"
)

func HandleStreamToClient(c *gin.Context) {
	path := filepath.Join(
		config.GetString("meeting.dir_path"),
		c.Params[0].Value,
		c.Params[1].Value,
		c.Params[2].Value,
	)
	c.File(path)
}

func HandleCheckStreamAvailable(c *gin.Context) {
	path := filepath.Join(
		config.GetString("meeting.dir_path"),
		c.Params[0].Value,
		c.Params[1].Value,
		c.Params[2].Value,
	)
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		c.Status(http.StatusNotFound)
		return
	}
	c.Status(http.StatusOK)
}
