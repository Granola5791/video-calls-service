package main

import (
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func HandleTranscription(meetingID uuid.UUID) {
	log.Println("started transcrioption")
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
		log.Println(t)
	}
}

func GetTranscription(meetingID uuid.UUID, userID uint) (string, error) {
	log.Println("started transcrioption for user", userID)
	reader, writer := io.Pipe()

	url := GetStringFromConfig("ai_server.url") + GetStringFromConfig("ai_server.api.transcription_path")
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return "", err
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
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK {
		return string(bodyBytes), nil
	}
	return string(bodyBytes), nil
}
