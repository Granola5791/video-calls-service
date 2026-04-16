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
	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))
	meetingParticipants, err := db.GetAllMeetingParticipantIDsFromDB(meetingID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, meetingParticipants)
}

func HandleTranscriptSummaryRequest(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))

	summary, err := db.GetSummaryFromDB(meetingID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, summary)
}

func HandleGetMeetingInfos(c *gin.Context) {
	from, err := time.Parse(
		time.RFC3339,
		c.Query(config.GetStringFromConfig("server.api.query_params.from_name")),
	)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	to, err := time.Parse(
		time.RFC3339,
		c.Query(config.GetStringFromConfig("server.api.query_params.to_name")),
	)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	hostName := c.Query(config.GetStringFromConfig("server.api.query_params.host_name"))
	meetingName := c.Query(config.GetStringFromConfig("server.api.query_params.meeting_name"))

	meetingIDs, err := db.GetMeetingInfosFromDB(from, to, hostName, meetingName)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, meetingIDs)
}

func HandleGetTranscript(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))
	participantID, err := strconv.Atoi(c.Param(config.GetStringFromConfig("server.api.params.participant_id_name")))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	username, err := db.GetUsernameFromDB(uint(participantID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	transcript, err := db.GetTranscriptFromDB(meetingID, uint(participantID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	usernameJsonName := config.GetStringFromConfig("json.username_name")
	transcriptJsonName := config.GetStringFromConfig("json.transcript_name")
	c.JSON(http.StatusOK, gin.H{usernameJsonName: username, transcriptJsonName: transcript})
}
