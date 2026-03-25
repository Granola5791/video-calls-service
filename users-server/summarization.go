package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleTranscriptSummary(meetingID uuid.UUID) {

	transcripts, err := GetMeetingTranscriptsFromDB(meetingID)
	if err != nil {
		log.Println(err)
		return
	}

	summary, err := GetSummary(transcripts)
	if err != nil {
		log.Println(err)
		return
	}

	err = UpdateSummaryToDB(meetingID, summary)
	if err != nil {
		log.Println(err)
		return
	}
}

func GetSummary(transcriptions []ParticipantTranscription) (string, error) {
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		for _, tr := range transcriptions {
			fmt.Fprintf(writer, "User \"{%s}\":\n%s\n", tr.User.Username, tr.Transcript)
		}
	}()

	url := fmt.Sprintf("%s%s",
		GetStringFromConfig("ai_server.url"),
		GetStringFromConfig("ai_server.api.summary_path"),
	)

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func HandleTranscriptSummaryRequest(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))

	summary, err := GetSummaryFromDB(meetingID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, summary)
}
