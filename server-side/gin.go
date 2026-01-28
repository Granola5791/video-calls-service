package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader websocket.Upgrader

func InitWsUpgrader() {
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  GetIntFromConfig("websocket.read_buffer_size"),
		WriteBufferSize: GetIntFromConfig("websocket.write_buffer_size"),
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin") == GetStringFromConfig("server.frontend_addr")
		},
	}
}

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
	router.GET(GetStringFromConfig("server.api.get_call_notifications_path"), RequireAuthentication, HandleGetCallNotifications)

	router.POST(GetStringFromConfig("server.api.signup_path"), HandleSignup)
	router.POST(GetStringFromConfig("server.api.login_path"), HandleLogin)
	router.POST(GetStringFromConfig("server.api.create_meeting_path"), RequireAuthentication, HandleCreateMeeting)
	router.POST(GetStringFromConfig("server.api.join_meeting_path"), RequireAuthentication, HandleJoinMeeting)
	router.POST(GetStringFromConfig("server.api.leave_meeting_path"), RequireAuthentication, HandleLeaveMeeting)
	router.POST(GetStringFromConfig("server.api.keep_alive_path"), RequireAuthentication, RequireKeepAliveToken, HandleKeepAlive)

	router.Run(GetStringFromConfig("server.listen_addr"))
}
