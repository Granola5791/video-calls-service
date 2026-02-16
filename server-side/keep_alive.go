package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	meetingKeepAliveMap   = make(map[uuid.UUID]*MeetingKeepAliveStruct)
	meetingKeepAliveMutex sync.Mutex
)

type MeetingKeepAliveStruct struct {
	ID           uuid.UUID
	mutex        sync.Mutex
	CurrToken    string
	TokenTimer   *time.Timer
	Participants map[uint]*time.Timer
}

func (m *MeetingKeepAliveStruct) Init(id uuid.UUID) {
	m.ID = id
	m.CurrToken = ""
	m.TokenTimer = time.AfterFunc(0, m.SetNewToken)
	m.Participants = make(map[uint]*time.Timer)
}

func (m *MeetingKeepAliveStruct) SetNewToken() {
	exp := GetIntFromConfig("keep_alive.token_exp")
	timerInterval := GetIntFromConfig("keep_alive.token_regen_interval")
	token, err := GenerateKeepAliveToken([]byte(os.Getenv("KEEP_ALIVE_JWT_SECRET")), m.ID, exp)
	if err != nil {
		log.Println(err)
		return
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.CurrToken = token
	if m.TokenTimer != nil {
		m.TokenTimer.Stop()
	}
	m.TokenTimer = time.AfterFunc(time.Duration(timerInterval)*time.Second, m.SetNewToken)
}

func (m *MeetingKeepAliveStruct) AddParticipant(participantID uint) {
	exp := GetIntFromConfig("keep_alive.token_exp")
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Participants[participantID] = time.AfterFunc(time.Duration(exp)*time.Second, func() {
		LogEventToDB(m.ID, participantID, GetStringFromConfig("database.meeting_events.participant_timeout"))
		LeaveMeeting(m.ID, participantID)
	})
}

func (m *MeetingKeepAliveStruct) RemoveParticipant(participantID uint) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	timer := m.Participants[participantID]
	if timer != nil {
		timer.Stop()
	}
	delete(m.Participants, participantID)
}

func (m *MeetingKeepAliveStruct) Close() {
	m.TokenTimer.Stop()
	for _, timer := range m.Participants {
		timer.Stop()
	}
}

func (m *MeetingKeepAliveStruct) RefreshParticipantTimer(participantID uint) (stillAlive bool) {
	exp := GetIntFromConfig("keep_alive.token_exp")
	participantTimer := m.Participants[participantID]
	if participantTimer == nil {
		return false
	}
	participantTimer.Reset(time.Duration(exp) * time.Second)
	return true
}

func (m *MeetingKeepAliveStruct) GetToken() string {
	return m.CurrToken
}

func RemoveMeetingKeepAlive(meetingID uuid.UUID) {
	meetingKeepAliveMutex.Lock()
	defer meetingKeepAliveMutex.Unlock()
	meeting := meetingKeepAliveMap[meetingID]
	if meeting != nil {
		meeting.Close()
	}
	delete(meetingKeepAliveMap, meetingID)
}

func AddMeetingKeepAlive(meetingID uuid.UUID) *MeetingKeepAliveStruct {
	meeting := &MeetingKeepAliveStruct{}
	meeting.Init(meetingID)
	meetingKeepAliveMutex.Lock()
	meetingKeepAliveMap[meetingID] = meeting
	meetingKeepAliveMutex.Unlock()
	return meeting
}
