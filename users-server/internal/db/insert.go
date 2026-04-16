package db

import "github.com/google/uuid"

func InsertUser(username string, hashedPassword string, salt string) error {
	user := User{
		Username:       username,
		HashedPassword: hashedPassword,
		Salt:           salt,
	}
	return db.Create(&user).Error
}

func CreateMeeting(hostID uint, isFaceDetectionRequired bool) (uuid.UUID, error) {
	meeting := Meeting{
		HostID:                  hostID,
		IsFaceDetectionRequired: isFaceDetectionRequired,
	}
	err := db.Create(&meeting).Error
	if err != nil {
		return uuid.Nil, err
	}
	return meeting.ID, nil
}

func AddParticipantToMeeting(meetingID uuid.UUID, userID uint) error {
	meetingParticipant := MeetingParticipant{
		UserID:    userID,
		MeetingID: meetingID,
	}
	return db.Create(&meetingParticipant).Error
}

func LogEvent(meetingID uuid.UUID, userID uint, event string) error {
	meetingEvent := MeetingEvent{
		MeetingID: meetingID,
		UserID:    userID,
		Event:     event,
	}
	return db.Create(&meetingEvent).Error
}

func InsertTranscription(meetingID uuid.UUID, userID uint, transcription string) error {
	return db.
		Create(&ParticipantTranscription{
			MeetingID:  meetingID,
			UserID:     userID,
			Transcript: transcription,
		}).Error
}
