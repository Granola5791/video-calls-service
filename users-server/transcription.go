package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func HandleTranscription(meetingID uuid.UUID) {
	meetingParticipants, err := GetAllMeetingParticipantIDsFromDB(meetingID)
	if err != nil {
		log.Println(err)
		return
	}
	for _, participant := range meetingParticipants {
		t, err := GetTranscription(meetingID, participant)
		if err != nil {
			log.Println(err)
			return
		}
		var s []string
		json.Unmarshal(t, &s)
		standardizedText := StandardizeTranscriptionText(s)
		fmt.Println(standardizedText)
	}
}

func GetTranscription(meetingID uuid.UUID, userID uint) ([]byte, error) {
	reader, writer := io.Pipe()

	url := GetStringFromConfig("ai_server.url") + GetStringFromConfig("ai_server.api.transcription_path")
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
// 	startTime1 endtime1 text1\n
// 	startTime2 endtime2 text2\n
//	...
//
// an example of a standardize transcription text
// can be found in the transcription_test.txt file
func StandardizeTranscriptionText(segments []string) string {
	return strings.Join(segments, "\n")
}