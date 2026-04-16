package summarization

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/Granola5791/video-calls-service/internal/db"
	"github.com/google/uuid"
)

type SummaryResponse struct {
	Summary string `json:"summary"`
	MeetingName string `json:"meeting_name"`
}

func HandleTranscriptSummary(meetingID uuid.UUID) {

	transcripts, err := db.GetMeetingTranscriptsFromDB(meetingID)
	if err != nil {
		log.Println(err)
		return
	}

	response, err := GetSummary(transcripts)
	if err != nil {
		log.Println(err)
		return
	}

	err = db.UpdateMeetingNameToDB(meetingID, response.MeetingName)
	if err != nil {
		log.Println(err)
		return
	}

	err = db.UpdateSummaryToDB(meetingID, response.Summary)
	if err != nil {
		log.Println(err)
		return
	}
}

func GetSummary(transcriptions []db.ParticipantTranscription) (SummaryResponse, error) {
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()
		for _, tr := range transcriptions {
			fmt.Fprintf(writer, "User \"{%s}\":\n%s\n", tr.User.Username, tr.Transcript)
		}
	}()

	url := fmt.Sprintf("%s%s",
		config.GetStringFromConfig("ai_server.url"),
		config.GetStringFromConfig("ai_server.api.summary_path"),
	)

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return SummaryResponse{}, err
	}
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return SummaryResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SummaryResponse{}, err
	}

	ret := SummaryResponse{}
	err = json.Unmarshal(body, &ret)
	if err != nil {
		return SummaryResponse{}, err
	}

	return ret, nil
}