package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleTranscription(meetingID uuid.UUID) {
	meetingParticipants, err := GetAllMeetingParticipantIDsFromDB(meetingID)
	if err != nil {
		log.Println(err)
		return
	}

	offsets, err := GetOffsetsOfUsers(meetingID, meetingParticipants)
	if err != nil {
		log.Println(err)
		return
	}

	for i, participant := range meetingParticipants {

		transcription, err := GetTranscription(meetingID, participant, offsets[i])
		if err != nil {
			log.Println(err)
			return
		}
		var res []string
		json.Unmarshal(transcription, &res)
		standardizedText := StandardizeTranscriptionText(res)
		err = InsertTranscriptionToDB(meetingID, participant, standardizedText)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func GetTranscription(meetingID uuid.UUID, userID uint, offset float64) ([]byte, error) {
	reader, writer := io.Pipe()

	url := fmt.Sprintf("%s%s?%s=%f",
		GetStringFromConfig("ai_server.url"),
		GetStringFromConfig("ai_server.api.transcription_path"),
		GetStringFromConfig("ai_server.api.query_params.offset_name"),
		offset,
	)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "video/webm")

	go func() {
		defer writer.Close()
		err := PipeAllUserVideoChunksFromDB(meetingID, userID, writer)
		if err != nil {
			log.Println(err)
		}
	}()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode == http.StatusOK {
		return bodyBytes, nil
	}
	return bodyBytes, nil
}

// return a text with the following format:
//
//	startTime1 endtime1 text1\n
//	startTime2 endtime2 text2\n
//	...
//
// an example of a standardize transcription text
// can be found in the transcription_test.txt file
func StandardizeTranscriptionText(segments []string) string {
	return strings.Join(segments, "\n")
}

func GetOffsetsOfUsers(meetingID uuid.UUID, participants []uint) ([]float64, error) {
	starts := make([]time.Time, len(participants))
	offsets := make([]float64, len(participants))
	for i := range participants {
		firstChunk, err := GetFirstUserVideoChunkFromDB(meetingID, participants[i])
		if err != nil {
			return []float64{}, err
		}
		starts[i] = firstChunk.CreatedAt
	}
	first := MinTime(starts)

	for i := range participants {
		offsets[i] = starts[i].Sub(first).Seconds()
	}
	return offsets, nil
}

func MinTime(times []time.Time) time.Time {
	if len(times) == 0 {
		return time.Time{}
	}
	minTime := times[0]
	for _, time := range times {
		if time.Before(minTime) {
			minTime = time
		}
	}
	return minTime
}

func HandleGetTranscriptionMeetings(c *gin.Context) {
	meetingIDs, err := GetTranscriptionMeetingsFromDB()
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, meetingIDs)
}

func HandleGetTranscript(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))
	participantID, err := strconv.Atoi(c.Param(GetStringFromConfig("server.api.params.participant_id_name")))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	username, err := GetUsernameFromDB(uint(participantID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	transcript, err := GetTranscriptFromDB(meetingID, uint(participantID))
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	usernameJsonName := GetStringFromConfig("json.username_name")
	transcriptJsonName := GetStringFromConfig("json.transcript_name")
	c.JSON(http.StatusOK, gin.H{usernameJsonName: username, transcriptJsonName: transcript})
}
