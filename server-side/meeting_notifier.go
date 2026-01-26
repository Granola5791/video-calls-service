package main

import (
	"sync"

	"context"
	"github.com/google/uuid"
)


type MeetingNotifierStruct struct {
	ID                 uuid.UUID
	participants       map[uint]*ParticipantNotifierStruct
	notificationChanIn chan ParticipantNotification
	mutex              sync.Mutex
	ctx                context.Context
	cancel             context.CancelFunc
}

func (m *MeetingNotifierStruct) Init(id uuid.UUID) {
	m.ID = id
	m.participants = make(map[uint]*ParticipantNotifierStruct)
	m.notificationChanIn = make(chan ParticipantNotification, GetIntFromConfig("notifications.channel_buffer_size"))
	m.ctx, m.cancel = context.WithCancel(context.Background())
}

func (m *MeetingNotifierStruct) lock() {
	m.mutex.Lock()
}

func (m *MeetingNotifierStruct) unlock() {
	m.mutex.Unlock()
}

func (m *MeetingNotifierStruct) AddParticipant(participantID uint) *ParticipantNotifierStruct {
	participant := &ParticipantNotifierStruct{}
	participant.Init(participantID, m.ctx)
	m.lock()
	defer m.unlock()
	m.participants[participantID] = participant
	m.notificationChanIn <- ParticipantNotification{ParticipantID: participantID, Event: ParticipantJoined}
	return participant
}

func (m *MeetingNotifierStruct) RemoveParticipant(participantID uint) {
	m.lock()
	defer m.unlock()
	m.notificationChanIn <- ParticipantNotification{ParticipantID: participantID, Event: ParticipantLeft}
	m.participants[participantID].Close()
	delete(m.participants, participantID)
}

func (m *MeetingNotifierStruct) NotifyParticipants(notification ParticipantNotification) {
	for _, participant := range m.participants {
		participant.channel <- notification
	}
}

func (m *MeetingNotifierStruct) Close() {
	m.cancel()
}

func (m *MeetingNotifierStruct) Run() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case notification := <-m.notificationChanIn:
			m.NotifyParticipants(notification)
		}
	}
}
