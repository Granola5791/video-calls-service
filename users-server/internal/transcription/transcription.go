package transcription

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/Granola5791/video-calls-service/internal/db"
	"github.com/google/uuid"
)

type Offset struct {
	Offset float64   `json:"offset"`
	Time   time.Time `json:"time"`
}

func HandleTranscription(meetingID uuid.UUID) error {
	meetingParticipants, err := db.GetAllMeetingParticipantIDsFromDB(meetingID)
	if err != nil {
		return err
	}

	offsets, err := GetOffsetsOfUsers(meetingID, meetingParticipants)
	if err != nil {
		return err
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
				return err
			}
			var res []string
			json.Unmarshal(transcription, &res)
			fullTranscription = append(fullTranscription, res...)
		}
		standardizedText := StandardizeTranscriptionText(fullTranscription)
		err = db.InsertTranscriptionToDB(meetingID, participant, standardizedText)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetTranscription(meetingID uuid.UUID, userID uint, offset float64, minTime time.Time, maxTime time.Time) ([]byte, error) {
	reader, writer := io.Pipe()

	url := fmt.Sprintf("%s%s?%s=%f",
		config.GetStringFromConfig("ai_server.url"),
		config.GetStringFromConfig("ai_server.api.transcription_path"),
		config.GetStringFromConfig("ai_server.api.query_params.offset_name"),
		offset,
	)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return []byte{}, err
	}
	req.Header.Set("Content-Type", "video/webm")

	go func() {
		defer writer.Close()
		err := db.PipeUserVideoChunksBetweenFromDB(meetingID, userID, minTime, maxTime, writer)
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
