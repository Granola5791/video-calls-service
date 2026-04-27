package keep_alive

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/Granola5791/video-calls-service/internal/auth"
	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/google/uuid"
)

type MeetingKeepAliveStruct struct {
	ID                  uuid.UUID
	mutex               sync.Mutex
	CurrToken           string
	TokenTimer          *time.Timer
	TokenTimerStartTime time.Time
	Participants        map[uint]*time.Timer
}

func (m *MeetingKeepAliveStruct) Init(id uuid.UUID) {
	m.ID = id
	m.CurrToken = ""
	m.TokenTimer = time.AfterFunc(0, m.SetNewToken)
	m.Participants = make(map[uint]*time.Timer)
}

func (m *MeetingKeepAliveStruct) SetNewToken() {
	exp := config.GetInt("keep_alive.token_exp")
	timerInterval := config.GetInt("keep_alive.token_regen_interval")
	token, err := auth.GenerateKeepAliveToken([]byte(os.Getenv("KEEP_ALIVE_JWT_SECRET")), m.ID, exp)
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
	m.TokenTimerStartTime = time.Now()
}

func (m *MeetingKeepAliveStruct) AddParticipant(participantID uint, onEnd func()) {
	exp := config.GetInt("keep_alive.token_exp")
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Participants[participantID] = time.AfterFunc(time.Duration(exp)*time.Second, onEnd)
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
	exp := config.GetInt("keep_alive.token_exp")
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

func (m *MeetingKeepAliveStruct) GetTokenStartTime() time.Time {
	return m.TokenTimerStartTime
}

func (m *MeetingKeepAliveStruct) GetTokenExpTime() time.Time {
	return m.TokenTimerStartTime.Add(time.Duration(config.GetInt("keep_alive.token_exp")) * time.Second)
}

func (m *MeetingKeepAliveStruct) GetTokenRemainingTime() time.Duration {
	return time.Until(m.GetTokenExpTime())
}

func (m *MeetingKeepAliveStruct) CloseAllParticipants() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, timer := range m.Participants {
		timer.Stop()
	}
}
