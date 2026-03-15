package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	meetingNotifiers      = make(map[uuid.UUID]*MeetingNotifierStruct)
	meetingNotifiersMutex sync.Mutex
)

func AddMeetingNotifier(meetingID uuid.UUID) *MeetingNotifierStruct {
	var meeting MeetingNotifierStruct
	meeting.Init(meetingID)
	meetingNotifiersMutex.Lock()
	meetingNotifiers[meetingID] = &meeting
	meetingNotifiersMutex.Unlock()
	return &meeting
}

func RemoveMeetingNotifier(meetingID uuid.UUID) {
	meetingNotifiersMutex.Lock()
	defer meetingNotifiersMutex.Unlock()
	meeting, ok := meetingNotifiers[meetingID]
	if !ok {
		return
	}
	meeting.CloseAllParticipants()
	meeting.Close()
	delete(meetingNotifiers, meetingID)
}

func HandleCreateMeeting(c *gin.Context) {
	userID := c.GetInt(GetStringFromConfig("jwt.user_id_name"))
	isFaceDetectionRequired, err := strconv.ParseBool(c.Param(GetStringFromConfig("server.api.params.is_face_detection_required_name")))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.Println("user id:", userID)
	meetingID, err := CreateMeetingInDB(uint(userID), isFaceDetectionRequired)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	jwtKey := []byte(os.Getenv("MEETING_JWT_SECRET"))
	token, err := GenerateMeetingToken(meetingID, jwtKey, GetIntFromConfig("meeting.token_exp"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": GetStringFromConfig("errors.internal_server_error")})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     GetStringFromConfig("meeting.token_name"),
		Value:    token,
		Path:     "/", // visible to all paths
		Domain:   GetStringFromConfig("jwt.domain"),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   GetIntFromConfig("meeting.token_exp"),
	})

	meeting := AddMeetingNotifier(meetingID)
	go meeting.Run()

	AddMeetingKeepAlive(meetingID)

	c.String(http.StatusOK, meetingID.String())
}

func HandleJoinMeeting(c *gin.Context) {
	log.Println("started join call")
	userID := c.GetInt(GetStringFromConfig("jwt.user_id_name"))
	meetingID, err := uuid.Parse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	meetingParticipants, err := GetParticipantsInMeetingFromDB(meetingID, uint(userID))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Println("participants:", meetingParticipants)

	isHost, err := IsHostOfMeetingInDB(meetingID, uint(userID))
	if err != nil {
		log.Println(err)
	}

	err = AddParticipantToMeetingInDB(meetingID, uint(userID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Set keep alive cookie
	meetingKeepAlive := meetingKeepAliveMap[meetingID]
	meetingKeepAlive.AddParticipant(uint(userID))
	token := meetingKeepAlive.GetToken()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     GetStringFromConfig("keep_alive.token_cookie_name"),
		Value:    token,
		Path:     "/",
		Domain:   GetStringFromConfig("jwt.domain"),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   GetIntFromConfig("keep_alive.token_exp"),
	})

	err = LogEventToDB(meetingID, uint(userID), GetStringFromConfig("database.meeting_events.participant_joined"))
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		GetStringFromConfig("meeting.participants_name"): meetingParticipants,
		GetStringFromConfig("meeting.is_host_name"):      isHost,
	})
}

func HandleGetCallNotifications(c *gin.Context) {
	userID := c.GetInt(GetStringFromConfig("jwt.user_id_name"))
	userName := c.GetString(GetStringFromConfig("jwt.username_name"))
	meetingID := uuid.MustParse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))
	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	meeting, ok := meetingNotifiers[meetingID]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	participantNotifier := meeting.AddParticipant(uint(userID), userName)
	go participantNotifier.Run(ws)
}

func HandleLeaveMeeting(c *gin.Context) {
	userID := c.GetInt(GetStringFromConfig("jwt.user_id_name"))
	meetingID := uuid.MustParse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))
	err := LeaveMeeting(meetingID, uint(userID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = LogEventToDB(meetingID, uint(userID), GetStringFromConfig("database.meeting_events.participant_left"))
	if err != nil {
		log.Println(err)
	}

	c.Status(http.StatusOK)
}

func LeaveMeeting(meetingID uuid.UUID, participantID uint) error {
	isHost, _ := IsHostOfMeetingInDB(meetingID, participantID)
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

	err = RemoveParticipantFromMeetingInDB(meetingID, participantID)
	if err != nil {
		return err
	}

	err = RemoveParticipantKeepAlive(meetingID, participantID)
	if err != nil {
		return err
	}

	isEmpty, _ := IsMeetingEmptyInDB(meetingID)
	if isEmpty {
		err = RemoveMeeting(meetingID)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveParticipantNotifier(meetingID uuid.UUID, participantID uint) error {
	meeting, ok := meetingNotifiers[meetingID]
	if !ok {
		return errors.New(GetStringFromConfig("errors.meeting_not_found"))
	}
	meeting.RemoveParticipant(participantID)
	return nil
}

func RemoveParticipantKeepAlive(meetingID uuid.UUID, participantID uint) error {
	meeting, ok := meetingKeepAliveMap[meetingID]
	if !ok {
		return errors.New(GetStringFromConfig("errors.meeting_not_found"))
	}
	meeting.RemoveParticipant(participantID)
	return nil
}

func RemoveMeeting(meetingID uuid.UUID) error {
	RemoveMeetingNotifier(meetingID)

	RemoveMeetingKeepAlive(meetingID)

	err := RemoveAllMeetingParticipantsFromDB(meetingID)
	if err != nil {
		return err
	}

	go HandleTranscription(meetingID)

	return nil
}

func HandleKeepAlive(c *gin.Context) {
	userID := c.GetInt(GetStringFromConfig("jwt.user_id_name"))
	meetingID := uuid.MustParse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))

	meeting, ok := meetingKeepAliveMap[meetingID]
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
		Name:     GetStringFromConfig("keep_alive.token_cookie_name"),
		Value:    token,
		Path:     "/",
		Domain:   GetStringFromConfig("jwt.domain"),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   GetIntFromConfig("keep_alive.token_exp"),
	})

	c.Status(http.StatusOK)
}
