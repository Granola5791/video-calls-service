package api

import (
	"net/http"
	"os"

	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/Granola5791/video-calls-service/internal/login"
	"github.com/Granola5791/video-calls-service/internal/meeting"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader websocket.Upgrader

func InitWsUpgrader() {
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  config.GetIntFromConfig("websocket.read_buffer_size"),
		WriteBufferSize: config.GetIntFromConfig("websocket.write_buffer_size"),
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin") == config.GetStringFromConfig("server.frontend_addr")
		},
	}
}

func InitRouter() {
	router := gin.Default()

	router.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{config.GetStringFromConfig("server.frontend_addr")},
			AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
			AllowHeaders:     []string{"Content-Type"},
			AllowCredentials: true,
		}),
		RequireSameOrigin,
	)

	router.GET(config.GetStringFromConfig("server.api.check_login_path"), RequireAuthentication)
	router.GET(config.GetStringFromConfig("server.api.check_admin_path"), RequireAuthentication, RequireAdmin)
	router.GET(config.GetStringFromConfig("server.api.get_call_notifications_path"), RequireAuthentication, meeting.HandleGetCallNotifications)
	router.GET(config.GetStringFromConfig("server.api.is_able_to_join_meeting_path"), RequireAuthentication, RequireMeetingExists, RequireNotBanned)
	router.GET(config.GetStringFromConfig("server.api.get_meeting_infos_path"), RequireAuthentication, RequireAdmin, meeting.HandleGetMeetingInfos)
	router.GET(config.GetStringFromConfig("server.api.get_transcript_path"), RequireAuthentication, RequireAdmin, meeting.HandleGetTranscript)
	router.GET(config.GetStringFromConfig("server.api.get_all_meeting_participants_path"), RequireAuthentication, RequireAdmin, meeting.HandleGetAllMeetingParticipants)
	router.GET(config.GetStringFromConfig("server.api.get_summary_path"), RequireAuthentication, RequireAdmin, meeting.HandleTranscriptSummaryRequest)

	router.POST(config.GetStringFromConfig("server.api.signup_path"), login.HandleSignup)
	router.POST(config.GetStringFromConfig("server.api.login_path"), login.HandleLogin)
	router.POST(config.GetStringFromConfig("server.api.logout_path"), RequireAuthentication, login.HandleLogout)
	router.POST(config.GetStringFromConfig("server.api.create_meeting_path"), RequireAuthentication, meeting.HandleCreateMeeting)
	router.POST(config.GetStringFromConfig("server.api.join_meeting_path"), RequireAuthentication, RequireMeetingExists, RequireNotBanned, meeting.HandleJoinMeeting)
	router.POST(config.GetStringFromConfig("server.api.leave_meeting_path"), RequireAuthentication, meeting.HandleLeaveMeeting)
	router.POST(config.GetStringFromConfig("server.api.keep_alive_path"), RequireAuthentication, RequireKeepAliveToken, RequireFaceDetection, meeting.HandleKeepAlive)
	router.POST(config.GetStringFromConfig("server.api.kick_participant_path"), RequireAuthentication, RequireKeepAliveToken, RequireHost, meeting.HandleKickParticipant)

	router.RunTLS(
		config.GetStringFromConfig("server.listen_addr"),
		os.Getenv("TLS_CERT_PATH"),
		os.Getenv("TLS_KEY_PATH"),
	)
}
