package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{GetStringFromConfig("server.frontend_addr")},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}))

	router.GET(GetStringFromConfig("server.api.check_login_path"), RequireAuthentication)
	router.GET(GetStringFromConfig("server.api.check_admin_path"), RequireAuthentication, RequireAdmin)

	router.POST(GetStringFromConfig("server.api.signup_path"), HandleSignup)
	router.POST(GetStringFromConfig("server.api.login_path"), HandleLogin)
	router.POST(GetStringFromConfig("server.api.create_meeting_path"), RequireAuthentication, HandleCreateMeeting)

	router.Run(GetStringFromConfig("server.listen_addr"))
}