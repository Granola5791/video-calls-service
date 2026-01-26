package main

type Event int

const (
	ParticipantJoined Event = iota
	ParticipantLeft
)

type ParticipantNotification struct {
	ParticipantID uint  `json:"participant_id"`
	Event         Event `json:"event"`
}