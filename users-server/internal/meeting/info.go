package meeting

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/Granola5791/video-calls-service/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleGetAllMeetingParticipants(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(config.GetString("server.api.params.meeting_id_name")))
	meetingParticipants, err := db.GetAllMeetingParticipantIDs(meetingID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, meetingParticipants)
}

func HandleTranscriptSummaryRequest(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(config.GetString("server.api.params.meeting_id_name")))

	summary, err := db.GetSummary(meetingID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, summary)
}

func HandleGetMeetingsInfo(c *gin.Context) {
	from, err := time.Parse(
		time.RFC3339,
		c.Query(config.GetString("server.api.query_params.from_name")),
	)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	to, err := time.Parse(
		time.RFC3339,
		c.Query(config.GetString("server.api.query_params.to_name")),
	)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	hostName := c.Query(config.GetString("server.api.query_params.host_name"))
	meetingName := c.Query(config.GetString("server.api.query_params.meeting_name"))

	meetingIDs, err := db.GetMeetingsInfo(from, to, hostName, meetingName)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, meetingIDs)
}

func HandleGetTranscript(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(config.GetString("server.api.params.meeting_id_name")))
	participantID, err := strconv.Atoi(c.Param(config.GetString("server.api.params.participant_id_name")))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	username, err := db.GetUsername(uint(participantID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	transcript, err := db.GetTranscript(meetingID, uint(participantID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	usernameJsonName := config.GetString("json.username_name")
	transcriptJsonName := config.GetString("json.transcript_name")
	c.JSON(http.StatusOK, gin.H{usernameJsonName: username, transcriptJsonName: transcript})
}
