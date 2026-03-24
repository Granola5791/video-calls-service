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

type Offset struct {
	Offset float64   `json:"offset"`
	Time   time.Time `json:"time"`
}

func HandleTranscription(meetingID uuid.UUID) {
	meetingParticipants, err := GetAllMeetingParticipantIDsFromDB(meetingID)
	if err != nil {
		log.Println(err)
		return
	}

	summaryCh := make(chan string)
	go HandleTranscriptSummary(meetingID, len(meetingParticipants), summaryCh)

	offsets, err := GetOffsetsOfUsers(meetingID, meetingParticipants)
	if err != nil {
		log.Println(err)
		return
	}

	for i, participant := range meetingParticipants {
		fullTranscription := []string{}

		for j := range offsets[i] {
			minTime := time.Time{}
			maxTime := time.Time{}
			if j != 0 {
				minTime = offsets[i][j].Time
			}
			if j != len(offsets[i])-1 {
				maxTime = offsets[i][j+1].Time
			}

			transcription, err := GetTranscription(meetingID, participant, offsets[i][j].Offset, minTime, maxTime)
			if err != nil {
				log.Println(err)
				return
			}
			var res []string
			json.Unmarshal(transcription, &res)
			fullTranscription = append(fullTranscription, res...)
		}
		standardizedText := StandardizeTranscriptionText(fullTranscription)
		summaryCh <- standardizedText
		err = InsertTranscriptionToDB(meetingID, participant, standardizedText)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func GetTranscription(meetingID uuid.UUID, userID uint, offset float64, minTime time.Time, maxTime time.Time) ([]byte, error) {
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
		err := PipeUserVideoChunksBetweenFromDB(meetingID, userID, minTime, maxTime, writer)
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
	if resp.StatusCode != http.StatusOK {
		log.Println(resp.StatusCode, resp.Status)
		return []byte{}, err
	}

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

func GetOffsetsOfUsers(meetingID uuid.UUID, participants []uint) ([][]Offset, error) {
	offsets := make([][]Offset, len(participants))
	for i := range participants {
		cnt, err := CountStartChunksFromDB(meetingID, participants[i])
		if err != nil {
			return [][]Offset{}, err
		}
		offsets[i] = make([]Offset, int(cnt))
		for j := 0; j < int(cnt); j++ {
			firstChunk, err := GetKthStartChunkFromDB(meetingID, participants[i], j)
			if err != nil {
				return [][]Offset{}, err
			}
			offsets[i][j].Time = firstChunk.CreatedAt
		}
	}
	first := MinTimeInOffsets(offsets)

	for i := range participants {
		for j := range offsets[i] {
			offsets[i][j].Offset = offsets[i][j].Time.Sub(first).Seconds()
		}
	}
	return offsets, nil
}

func MinTimeInOffsets(offsets [][]Offset) time.Time {
	if len(offsets) == 0 {
		return time.Time{}
	}
	minTime := offsets[0][0].Time
	for i := range offsets {
		for j := range offsets[i] {
			if offsets[i][j].Time.Before(minTime) {
				minTime = offsets[i][j].Time
			}
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
