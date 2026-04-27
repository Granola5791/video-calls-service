package notifications

import (
	"sync"

	"github.com/google/uuid"
)

var (
	MeetingNotifiers      = make(map[uuid.UUID]*MeetingNotifierStruct)
	meetingNotifiersMutex sync.Mutex
)

func AddMeetingNotifier(meetingID uuid.UUID) *MeetingNotifierStruct {
	var meeting MeetingNotifierStruct
	meeting.Init(meetingID)
	meetingNotifiersMutex.Lock()
	MeetingNotifiers[meetingID] = &meeting
	meetingNotifiersMutex.Unlock()
	return &meeting
}

func RemoveMeetingNotifier(meetingID uuid.UUID) {
	meetingNotifiersMutex.Lock()
	defer meetingNotifiersMutex.Unlock()
	meeting, ok := MeetingNotifiers[meetingID]
	if !ok {
		return
	}
	meeting.CloseAllParticipants()
	meeting.Close()
	delete(MeetingNotifiers, meetingID)
}