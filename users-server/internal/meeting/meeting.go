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
	userID := c.GetInt(config.GetStringFromConfig("jwt.user_id_name"))
	isFaceDetectionRequired, err := strconv.ParseBool(c.Param(config.GetStringFromConfig("server.api.params.is_face_detection_required_name")))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.Println("user id:", userID)
	meetingID, err := db.CreateMeetingInDB(uint(userID), isFaceDetectionRequired)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	jwtKey := []byte(os.Getenv("MEETING_JWT_SECRET"))
	token, err := auth.GenerateMeetingToken(meetingID, jwtKey, config.GetIntFromConfig("meeting.token_exp"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.GetStringFromConfig("errors.internal_server_error")})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.GetStringFromConfig("meeting.token_name"),
		Value:    token,
		Path:     "/", // visible to all paths
		Domain:   config.GetStringFromConfig("jwt.domain"),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   config.GetIntFromConfig("meeting.token_exp"),
	})

	meeting := notifications.AddMeetingNotifier(meetingID)
	go meeting.Run()

	keep_alive.AddMeetingKeepAlive(meetingID)

	c.String(http.StatusOK, meetingID.String())
}

func HandleJoinMeeting(c *gin.Context) {
	log.Println("started join call")
	userID := c.GetInt(config.GetStringFromConfig("jwt.user_id_name"))
	meetingID, err := uuid.Parse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	meetingParticipants, err := db.GetParticipantsInMeetingFromDB(meetingID, uint(userID))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Println("participants:", meetingParticipants)

	isHost, err := db.IsHostOfMeetingInDB(meetingID, uint(userID))
	if err != nil {
		log.Println(err)
	}

	err = db.AddParticipantToMeetingInDB(meetingID, uint(userID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Set keep alive cookie
	meetingKeepAlive := keep_alive.MeetingKeepAliveMap[meetingID]
	meetingKeepAlive.AddParticipant(uint(userID), func() {
		db.LogEventToDB(meetingID, uint(userID), config.GetStringFromConfig("database.meeting_events.participant_timeout"))
		LeaveMeeting(meetingID, uint(userID))
	})
	token := meetingKeepAlive.GetToken()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.GetStringFromConfig("keep_alive.token_cookie_name"),
		Value:    token,
		Path:     "/",
		Domain:   config.GetStringFromConfig("jwt.domain"),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   config.GetIntFromConfig("keep_alive.token_exp"),
	})

	err = db.LogEventToDB(meetingID, uint(userID), config.GetStringFromConfig("database.meeting_events.participant_joined"))
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		config.GetStringFromConfig("meeting.participants_name"): meetingParticipants,
		config.GetStringFromConfig("meeting.is_host_name"):      isHost,
	})
}

func HandleGetCallNotifications(c *gin.Context) {
	userID := c.GetInt(config.GetStringFromConfig("jwt.user_id_name"))
	userName := c.GetString(config.GetStringFromConfig("jwt.username_name"))
	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))
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
	userID := c.GetInt(config.GetStringFromConfig("jwt.user_id_name"))
	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))
	err := LeaveMeeting(meetingID, uint(userID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = db.LogEventToDB(meetingID, uint(userID), config.GetStringFromConfig("database.meeting_events.participant_left"))
	if err != nil {
		log.Println(err)
	}

	c.Status(http.StatusOK)
}

func LeaveMeeting(meetingID uuid.UUID, participantID uint) error {
	isHost, _ := db.IsHostOfMeetingInDB(meetingID, participantID)
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

	err = db.RemoveParticipantFromMeetingInDB(meetingID, participantID)
	if err != nil {
		return err
	}

	err = RemoveParticipantKeepAlive(meetingID, participantID)
	if err != nil {
		return err
	}

	isEmpty, _ := db.IsMeetingEmptyInDB(meetingID)
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
		return errors.New(config.GetStringFromConfig("errors.meeting_not_found"))
	}
	meeting.RemoveParticipant(participantID)
	return nil
}

func RemoveParticipantKeepAlive(meetingID uuid.UUID, participantID uint) error {
	meeting, ok := keep_alive.MeetingKeepAliveMap[meetingID]
	if !ok {
		return errors.New(config.GetStringFromConfig("errors.meeting_not_found"))
	}
	meeting.RemoveParticipant(participantID)
	return nil
}

func RemoveMeeting(meetingID uuid.UUID) error {
	notifications.RemoveMeetingNotifier(meetingID)

	keep_alive.RemoveMeetingKeepAlive(meetingID)

	err := db.RemoveAllMeetingParticipantsFromDB(meetingID)
	if err != nil {
		return err
	}

	go func() {
		err := transcription.HandleTranscription(meetingID)
		if err != nil {
			log.Println(err)
			return
		}
		summarization.HandleTranscriptSummary(meetingID)
	}()

	return nil
}

func HandleKeepAlive(c *gin.Context) {
	userID := c.GetInt(config.GetStringFromConfig("jwt.user_id_name"))
	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))

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
		Name:     config.GetStringFromConfig("keep_alive.token_cookie_name"),
		Value:    token,
		Path:     "/",
		Domain:   config.GetStringFromConfig("jwt.domain"),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   config.GetIntFromConfig("keep_alive.token_exp"),
	})

	c.Status(http.StatusOK)
}