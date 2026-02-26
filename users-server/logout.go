package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleLogout(c *gin.Context) {
	c.SetCookie(GetStringFromConfig("jwt.token_cookie_name"), "", -1, "/", "", false, true) // Delete the cookie by setting its MaxAge to -1
	c.Status(http.StatusOK)
}