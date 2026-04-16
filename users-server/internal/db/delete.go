package db

import "github.com/google/uuid"

func RemoveParticipantFromMeeting(meetingID uuid.UUID, userID uint) error {
	return db.
		Unscoped().
		Where("meeting_id = ? AND user_id = ?", meetingID, userID).
		Delete(&MeetingParticipant{}).Error
}

func RemoveAllMeetingParticipants(meetingID uuid.UUID) error {
	return db.
		Unscoped().
		Where("meeting_id = ?", meetingID).
		Delete(&MeetingParticipant{}).Error
}