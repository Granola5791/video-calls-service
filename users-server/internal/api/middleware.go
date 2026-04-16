package api

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Granola5791/video-calls-service/internal/auth"
	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/Granola5791/video-calls-service/internal/db"
	"github.com/Granola5791/video-calls-service/internal/face_detection"
	"github.com/Granola5791/video-calls-service/internal/meeting"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func RequireAuthentication(c *gin.Context) {

	// get cookie
	tokenName := config.GetStringFromConfig("jwt.token_cookie_name")
	tokenString, err := c.Cookie(tokenName)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := auth.ParseToken(tokenString, jwtKey)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// check if token is expired
	if claims[config.GetStringFromConfig("jwt.exp_name")].(float64) < float64(time.Now().Unix()) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(config.GetStringFromConfig("jwt.user_id_name"), int(claims[config.GetStringFromConfig("jwt.user_id_name")].(float64)))
	c.Set(config.GetStringFromConfig("jwt.username_name"), claims[config.GetStringFromConfig("jwt.username_name")].(string))
	c.Set(config.GetStringFromConfig("jwt.role_name"), claims[config.GetStringFromConfig("jwt.role_name")].(string))

	c.Next()
}

// can only be called after RequireAuthentication
func RequireAdmin(c *gin.Context) {
	role, _ := c.Get(config.GetStringFromConfig("jwt.role_name"))
	roleString, ok := role.(string)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if !strings.Contains(roleString, config.GetStringFromConfig("jwt.admin_role")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}

func RequireKeepAliveToken(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))

	// get cookie
	tokenName := config.GetStringFromConfig("keep_alive.token_cookie_name")
	tokenString, err := c.Cookie(tokenName)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token
	jwtKey := []byte(os.Getenv("KEEP_ALIVE_JWT_SECRET"))
	token, err := auth.ParseToken(tokenString, jwtKey)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// check if token is for the correct meeting
	if claims[config.GetStringFromConfig("jwt.meeting_id_name")].(string) != meetingID.String() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// check if token is expired
	if claims[config.GetStringFromConfig("jwt.exp_name")].(float64) < float64(time.Now().Unix()) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}

// can only be called after RequireAuthentication
func RequireHost(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))
	userID := c.GetInt(config.GetStringFromConfig("jwt.user_id_name"))
	isHost, err := db.IsHostOfMeetingInDB(meetingID, uint(userID))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !isHost {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func RequireNotBanned(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))
	userID := c.GetInt(config.GetStringFromConfig("jwt.user_id_name"))
	isBanned, err := db.IsBannedFromMeetingInDB(meetingID, uint(userID))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if isBanned {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func RequireMeetingExists(c *gin.Context) {
	meetingID, err := uuid.Parse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	exists, err := db.MeetingExistsInDB(meetingID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !exists {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
}

func RequireSameOrigin(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	if origin != config.GetStringFromConfig("server.frontend_addr") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func RequireFaceDetection(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))
	userID := c.GetInt(config.GetStringFromConfig("jwt.user_id_name"))

	// check if face detection is required
	required, err := db.IsFaceDetectionRequiredInDB(meetingID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !required {
		return
	}

	videoChunks, err := db.GetUserVideoChunksFromDB(meetingID, uint(userID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if len(videoChunks) == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	go db.MarkUserVideoChunksAsVisitedInDB(
		meetingID,
		uint(userID),
		videoChunks[0].ChunkNumber,
		videoChunks[len(videoChunks)-1].ChunkNumber,
	)

	outputPipeRead, err := face_detection.ConcatenateChunks(videoChunks)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	defer outputPipeRead.Close()

	url := config.GetStringFromConfig("ai_server.url") + config.GetStringFromConfig("ai_server.api.face_detection_path")
	framesWithFace, totalFrames, err := face_detection.SendvideoToFaceDetector(url, outputPipeRead)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !face_detection.PassedFaceDetectionThreshold(framesWithFace, totalFrames) {
		meeting.KickParticipantFromMeeting(meetingID, uint(userID))
		db.LogEventToDB(meetingID, uint(userID), config.GetStringFromConfig("database.meeting_events.participant_kicked_by_face_detection"))
		meeting.SendDangerPeriodNotification(meetingID, uint(userID))
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
