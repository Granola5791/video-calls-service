package transcription

import (
	"time"

	"github.com/Granola5791/video-calls-service/internal/db"
	"github.com/google/uuid"
)

func GetOffsetsOfUsers(meetingID uuid.UUID, participants []uint) ([][]Offset, error) {
	offsets := make([][]Offset, len(participants))
	for i := range participants {
		cnt, err := db.CountStartChunks(meetingID, participants[i])
		if err != nil {
			return [][]Offset{}, err
		}
		offsets[i] = make([]Offset, int(cnt))
		for j := 0; j < int(cnt); j++ {
			firstChunk, err := db.GetKthStartChunk(meetingID, participants[i], j)
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
