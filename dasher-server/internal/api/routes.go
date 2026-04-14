package api

import (
	"os"

	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/Granola5791/video-calls-service/internal/stream"
)

func InitRouter() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.GetStringFromConfig("server.frontend_addr")},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
	}),
		RequireSameOrigin,
		RequireKeepAliveToken,
	)

	router.GET(config.GetStringFromConfig("server.api.stream_from_client_path"), RequireAuthentication, stream.HandleStream)
	router.GET(config.GetStringFromConfig("server.api.stream_to_client_path"), RequireAuthentication, stream.HandleStreamToClient)
	router.HEAD(config.GetStringFromConfig("server.api.stream_to_client_path"), RequireAuthentication, stream.HandleCheckStreamAvailable)
	
	router.POST(config.GetStringFromConfig("server.api.create_meeting_path"), RequireAuthentication, RequireAuthorizedMeeting, stream.HandleCreateMeeting)
	router.POST(config.GetStringFromConfig("server.api.join_meeting_path"), RequireAuthentication, stream.HandleJoinMeeting)

	router.RunTLS(
		config.GetStringFromConfig("server.listen_addr"),
		os.Getenv("TLS_CERT_PATH"),
		os.Getenv("TLS_KEY_PATH"),
	)
}
