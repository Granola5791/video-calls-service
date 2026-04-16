package meeting

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Granola5791/video-calls-service/internal/auth"
	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/Granola5791/video-calls-service/internal/db"
	"github.com/Granola5791/video-calls-service/internal/keep_alive"
	"github.com/Granola5791/video-calls-service/internal/mywebsocket"
	"github.com/Granola5791/video-calls-service/internal/notifications"
	"github.com/Granola5791/video-calls-service/internal/summarization"
	"github.com/Granola5791/video-calls-service/internal/transcription"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleCreateMeeting(c *gin.Context) {
	userID := c.GetInt(config.GetString("jwt.user_id_name"))
	isFaceDetectionRequired, err := strconv.ParseBool(c.Param(config.GetString("server.api.params.is_face_detection_required_name")))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.Println("user id:", userID)
	meetingID, err := db.CreateMeeting(uint(userID), isFaceDetectionRequired)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	jwtKey := []byte(os.Getenv("MEETING_JWT_SECRET"))
	token, err := auth.GenerateMeetingToken(meetingID, jwtKey, config.GetInt("meeting.token_exp"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.GetString("errors.internal_server_error")})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.GetString("meeting.token_name"),
		Value:    token,
		Path:     "/", // visible to all paths
		Domain:   config.GetString("jwt.domain"),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   config.GetInt("meeting.token_exp"),
	})

	meeting := notifications.AddMeetingNotifier(meetingID)
	go meeting.Run()

	keep_alive.AddMeetingKeepAlive(meetingID)

	c.String(http.StatusOK, meetingID.String())
}

func HandleJoinMeeting(c *gin.Context) {
	log.Println("started join call")
	userID := c.GetInt(config.GetString("jwt.user_id_name"))
	meetingID, err := uuid.Parse(c.Param(config.GetString("server.api.params.meeting_id_name")))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	meetingParticipants, err := db.GetParticipantsInMeeting(meetingID, uint(userID))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Println("participants:", meetingParticipants)

	isHost, err := db.IsHostOfMeeting(meetingID, uint(userID))
	if err != nil {
		log.Println(err)
	}

	err = db.AddParticipantToMeeting(meetingID, uint(userID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Set keep alive cookie
	meetingKeepAlive := keep_alive.MeetingKeepAliveMap[meetingID]
	meetingKeepAlive.AddParticipant(uint(userID), func() {
		db.LogEvent(meetingID, uint(userID), config.GetString("database.meeting_events.participant_timeout"))
		LeaveMeeting(meetingID, uint(userID))
	})
	token := meetingKeepAlive.GetToken()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.GetString("keep_alive.token_cookie_name"),
		Value:    token,
		Path:     "/",
		Domain:   config.GetString("jwt.domain"),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   config.GetInt("keep_alive.token_exp"),
	})

	err = db.LogEvent(meetingID, uint(userID), config.GetString("database.meeting_events.participant_joined"))
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		config.GetString("meeting.participants_name"): meetingParticipants,
		config.GetString("meeting.is_host_name"):      isHost,
	})
}

func HandleGetCallNotifications(c *gin.Context) {
	userID := c.GetInt(config.GetString("jwt.user_id_name"))
	userName := c.GetString(config.GetString("jwt.username_name"))
	meetingID := uuid.MustParse(c.Param(config.GetString("server.api.params.meeting_id_name")))
	ws, err := mywebsocket.UpgradeToWebsocket(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	meeting, ok := notifications.MeetingNotifiers[meetingID]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	participantNotifier := meeting.AddParticipant(uint(userID), userName)
	go participantNotifier.Run(ws)
}

func HandleLeaveMeeting(c *gin.Context) {
	userID := c.GetInt(config.GetString("jwt.user_id_name"))
	meetingID := uuid.MustParse(c.Param(config.GetString("server.api.params.meeting_id_name")))
	err := LeaveMeeting(meetingID, uint(userID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = db.LogEvent(meetingID, uint(userID), config.GetString("database.meeting_events.participant_left"))
	if err != nil {
		log.Println(err)
	}

	c.Status(http.StatusOK)
}

func LeaveMeeting(meetingID uuid.UUID, participantID uint) error {
	isHost, _ := db.IsHostOfMeeting(meetingID, participantID)
	if isHost {
		err := RemoveMeeting(meetingID)
		if err != nil {
			return err
		}
		return nil
	}

	err := RemoveParticipantNotifier(meetingID, participantID)
	if err != nil {
		return err
	}

	err = db.RemoveParticipantFromMeeting(meetingID, participantID)
	if err != nil {
		return err
	}

	err = RemoveParticipantKeepAlive(meetingID, participantID)
	if err != nil {
		return err
	}

	isEmpty, _ := db.IsMeetingEmpty(meetingID)
	if isEmpty {
		err = RemoveMeeting(meetingID)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveParticipantNotifier(meetingID uuid.UUID, participantID uint) error {
	meeting, ok := notifications.MeetingNotifiers[meetingID]
	if !ok {
		return errors.New(config.GetString("errors.meeting_not_found"))
	}
	meeting.RemoveParticipant(participantID)
	return nil
}

func RemoveParticipantKeepAlive(meetingID uuid.UUID, participantID uint) error {
	meeting, ok := keep_alive.MeetingKeepAliveMap[meetingID]
	if !ok {
		return errors.New(config.GetString("errors.meeting_not_found"))
	}
	meeting.RemoveParticipant(participantID)
	return nil
}

func RemoveMeeting(meetingID uuid.UUID) error {
	notifications.RemoveMeetingNotifier(meetingID)

	keep_alive.RemoveMeetingKeepAlive(meetingID)

	err := db.RemoveAllMeetingParticipants(meetingID)
	if err != nil {
		return err
	}

	go func() {
		err := transcription.MakeMeetingTranscription(meetingID)
		if err != nil {
			log.Println(err)
			return
		}
		summarization.MakeTranscriptSummary(meetingID)
	}()

	return nil
}

func HandleKeepAlive(c *gin.Context) {
	userID := c.GetInt(config.GetString("jwt.user_id_name"))
	meetingID := uuid.MustParse(c.Param(config.GetString("server.api.params.meeting_id_name")))

	meeting, ok := keep_alive.MeetingKeepAliveMap[meetingID]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	stillAlive := meeting.RefreshParticipantTimer(uint(userID))
	if !stillAlive {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Set cookie
	token := meeting.GetToken()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.GetString("keep_alive.token_cookie_name"),
		Value:    token,
		Path:     "/",
		Domain:   config.GetString("jwt.domain"),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   config.GetInt("keep_alive.token_exp"),
	})

	c.Status(http.StatusOK)
}
