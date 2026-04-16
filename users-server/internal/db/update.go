package db

import "github.com/google/uuid"

// mark user video chunks in the range [minChunkNumber, maxChunkNumber] inclusive as visited
func MarkUserVideoChunksAsVisited(meetingID uuid.UUID, userID uint, minChunkNumber uint, maxChunkNumber uint) error {
	return db.
		Model(&UserVideoChunk{}).
		Where("meeting_id = ? AND user_id = ? AND chunk_number >= ? AND chunk_number <= ?", meetingID, userID, minChunkNumber, maxChunkNumber).
		Update("visited", true).Error
}

func UpdateSummary(meetingID uuid.UUID, summary string) error {
	return db.
		Model(&Meeting{}).
		Where("id = ?", meetingID).
		Update("summary", summary).Error
}

func UpdateMeetingName(meetingID uuid.UUID, meetingName string) error {
	return db.
		Model(&Meeting{}).
		Where("id = ?", meetingID).
		Update("name", meetingName).Error
}
