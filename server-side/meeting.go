package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var meetingNotifiers = make(map[uuid.UUID]*MeetingNotifierStruct)

type Event int

const (
	ParticipantJoined Event = iota
	ParticipantLeft
)

type ParticipantNotification struct {
	ParticipantID uint  `json:"participant_id"`
	Event         Event `json:"event"`
}

type MeetingNotifierStruct struct {
	ID                 uuid.UUID
	participants       map[uint]chan ParticipantNotification
	notificationChanIn chan ParticipantNotification
}

func (m *MeetingNotifierStruct) Init(id uuid.UUID) {
	m.ID = id
	m.participants = make(map[uint]chan ParticipantNotification)
	m.notificationChanIn = make(chan ParticipantNotification, GetIntFromConfig("notifications.channel_buffer_size"))
}

func (m *MeetingNotifierStruct) NotifyParticipants(notification ParticipantNotification) {
	for _, participant := range m.participants {
		participant <- notification
	}
}

func HandleCreateMeeting(c *gin.Context) {
	userID := c.GetInt(GetStringFromConfig("jwt.user_id_name"))
	log.Println("user id:", userID)
	meetingID, err := CreateMeetingInDB(userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	jwtKey := []byte(os.Getenv("MEETING_JWT_SECRET"))
	token, err := GenerateMeetingToken(meetingID, jwtKey, GetIntFromConfig("meeting.token_exp"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": GetStringFromConfig("errors.internal_server_error")})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     GetStringFromConfig("meeting.token_name"),
		Value:    token,
		Path:     "/", // visible to all paths
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   GetIntFromConfig("meeting.token_exp"),
	})

	go MeetingNotifier(meetingID)

	c.String(http.StatusOK, meetingID.String())
}

func MeetingNotifier(meetingID uuid.UUID) {
	var meeting MeetingNotifierStruct
	meeting.Init(meetingID)
	meetingNotifiers[meetingID] = &meeting

	log.Println("waiting for notifications")
	for notification := range meeting.notificationChanIn {
		log.Println("got notification")
		meeting.NotifyParticipants(notification)
	}
}

func HandleJoinMeeting(c *gin.Context) {
	log.Println("started join call")
	userID := c.GetInt(GetStringFromConfig("jwt.user_id_name"))
	meetingID, err := uuid.Parse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	meetingParticipants, err := GetMeetingParticipantIDsFromDB(meetingID, uint(userID))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Println("participants:", meetingParticipants)

	exists, err := IsParticipantInMeetingInDB(meetingID, userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !exists {
		err = AddParticipantToMeetingInDB(meetingID, userID)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, meetingParticipants)
}

func HandleGetCallNotifications(c *gin.Context) {
	userID := c.GetInt(GetStringFromConfig("jwt.user_id_name"))
	meetingID := uuid.MustParse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))
	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	defer ws.Close()

	meeting, ok := meetingNotifiers[meetingID]
	if !ok {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	notificationChannel := make(chan ParticipantNotification, GetIntFromConfig("notifications.channel_buffer_size"))
	meeting.participants[uint(userID)] = notificationChannel
	meeting.notificationChanIn <- ParticipantNotification{ParticipantID: uint(userID), Event: ParticipantJoined}

	for notification := range notificationChannel {
		log.Println("notification:", notification)
		if notification.ParticipantID != uint(userID) {
			err := ws.WriteJSON(notification)
			if err != nil {
				log.Println(err)
			}
		}
	}

}
