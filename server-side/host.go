package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequireHost(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))
	userID := c.GetInt(GetStringFromConfig("jwt.user_id_name"))
	isHost, err := IsHostOfMeetingInDB(meetingID, uint(userID))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !isHost {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

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