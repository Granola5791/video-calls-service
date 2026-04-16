package meeting

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/Granola5791/video-calls-service/internal/db"
	"github.com/Granola5791/video-calls-service/internal/keep_alive"
	"github.com/Granola5791/video-calls-service/internal/notifications"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleKickParticipant(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))
	userToKick := c.Param(config.GetStringFromConfig("server.api.params.user_to_kick_name"))
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

	db.LogEventToDB(meetingID, uint(userToKickInt), config.GetStringFromConfig("database.meeting_events.participant_kicked_by_host"))

	SendDangerPeriodNotification(meetingID, uint(userToKickInt))

	c.Status(http.StatusOK)
}

func KickParticipantFromMeeting(meetingID uuid.UUID, userToKick uint) error {
	err := db.BanUserFromMeetingInDB(meetingID, userToKick)
	if err != nil {
		return err
	}

	err = LeaveMeeting(meetingID, userToKick)
	return err
}

func SendDangerPeriodNotification(meetingID uuid.UUID, participantID uint) {
	meetingNotifier := notifications.MeetingNotifiers[meetingID]
	meetingKeepAlive := keep_alive.MeetingKeepAliveMap[meetingID]

	dangerPeriod := meetingKeepAlive.GetTokenRemainingTime() // this is the time left before the token expires and the player is kicked for certain.
	meetingNotifier.NotifyParticipants(notifications.ParticipantNotification{
		ParticipantID: participantID,
		Event:         notifications.ParticipantKickedByHost,
		Value:         dangerPeriod.Seconds(),
	})
}
