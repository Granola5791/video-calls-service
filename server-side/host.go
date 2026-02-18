package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleKickParticipant(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))
	userToKick := c.Param(GetStringFromConfig("server.api.params.user_to_kick_name"))
	userToKickInt, err := strconv.Atoi(userToKick)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = KickParticipantFromMeeting(meetingID, uint(userToKickInt))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	LogEventToDB(meetingID, uint(userToKickInt), GetStringFromConfig("database.meeting_events.participant_kicked_by_host"))

	SendDangerPeriodNotification(meetingID, uint(userToKickInt))
	
	c.Status(http.StatusOK)
}

func KickParticipantFromMeeting(meetingID uuid.UUID, userToKick uint) error {
	err := BanUserFromMeetingInDB(meetingID, userToKick)
	if err != nil {
		return err
	}

	err = LeaveMeeting(meetingID, userToKick)
	return err
}

func SendDangerPeriodNotification(meetingID uuid.UUID, participantID uint) {
	meetingNotifier := meetingNotifiers[meetingID]
	meetingKeepAlive := meetingKeepAliveMap[meetingID]

	dangerPeriod := meetingKeepAlive.GetTokenRemainingTime() // this is the time left before the token expires and the player is kicked for certain.
	meetingNotifier.NotifyParticipants(ParticipantNotification{
		ParticipantID: participantID,
		Event:         ParticipantKickedByHost,
		Value:         dangerPeriod.Seconds(),
	})
}
