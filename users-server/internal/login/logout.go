package login

import (
	"net/http"

	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/gin-gonic/gin"
)

func HandleLogout(c *gin.Context) {
	c.SetCookie(config.GetStringFromConfig("jwt.token_cookie_name"), "", -1, "/", "", false, true) // Delete the cookie by setting its MaxAge to -1
	c.Status(http.StatusOK)
}