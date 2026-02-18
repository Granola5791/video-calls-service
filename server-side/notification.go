package main

type Event int

const (
	ParticipantJoined Event = iota
	ParticipantLeft
	ParticipantKickedByHost
)

type ParticipantNotification struct {
	ParticipantID uint  `json:"participant_id"`
	Event         Event `json:"event"`
	Value         float64   `json:"value"`
}
