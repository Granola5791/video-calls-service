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
		ReadBufferSize:  GetIntFromConfig("stream.read_buffer_size"),
		WriteBufferSize: GetIntFromConfig("stream.write_buffer_size"),
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin") == GetStringFromConfig("server.frontend_addr")
		},
	}
}

func InitRouter() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{GetStringFromConfig("server.frontend_addr")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}),
		RequireKeepAliveToken,
	)

	router.GET(GetStringFromConfig("server.api.stream_from_client_path"), RequireAuthentication, HandleStream)

	router.POST(GetStringFromConfig("server.api.create_meeting_path"), RequireAuthentication, RequireAuthorizedMeeting, HandleCreateMeeting)
	router.POST(GetStringFromConfig("server.api.join_meeting_path"), RequireAuthentication, HandleJoinMeeting)

	router.StaticFS(GetStringFromConfig("server.api.stream_to_client_path"), gin.Dir("./meetings", false))

	router.Run(GetStringFromConfig("server.listen_addr"))
}
