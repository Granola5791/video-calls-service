package main

type Event int

const (
	ParticipantJoined Event = iota
	ParticipantLeft
	ParticipantKickedByHost
	MeetingEnded
)

type ParticipantNotification struct {
	ParticipantID   uint    `json:"participant_id"`
	ParticipantName string  `json:"participant_name"`
	Event           Event   `json:"event"`
	Value           float64 `json:"value"`
}
