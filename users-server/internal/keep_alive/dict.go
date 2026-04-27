package keep_alive

import (
	"sync"

	"github.com/google/uuid"
)

var (
	MeetingKeepAliveMap   = make(map[uuid.UUID]*MeetingKeepAliveStruct)
	meetingKeepAliveMutex sync.Mutex
)

func RemoveMeetingKeepAlive(meetingID uuid.UUID) {
	meetingKeepAliveMutex.Lock()
	defer meetingKeepAliveMutex.Unlock()
	meeting := MeetingKeepAliveMap[meetingID]
	if meeting != nil {
		meeting.Close()
	}
	delete(MeetingKeepAliveMap, meetingID)
}

func AddMeetingKeepAlive(meetingID uuid.UUID) *MeetingKeepAliveStruct {
	meeting := &MeetingKeepAliveStruct{}
	meeting.Init(meetingID)
	meetingKeepAliveMutex.Lock()
	MeetingKeepAliveMap[meetingID] = meeting
	meetingKeepAliveMutex.Unlock()
	return meeting
}
