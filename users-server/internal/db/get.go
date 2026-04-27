package db

import (
	"io"
	"strconv"
	"time"

	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParticipantInfo struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

type UserInfo struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type MeetingInfo struct {
	UuidModel
	HostID                  uint   `json:"host_id"`
	HostUsername            string `json:"host_username"`
	IsFaceDetectionRequired bool   `json:"is_face_detection_required"`
	Name                    string `json:"name"`
}

func UserExists(username string) (bool, error) {
	var count int64
	err := db.Model(&User{}).
		Where("username = ?", username).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUserAuth(username string) (string, string, error) {
	var userAuth UserAuth
	err := db.Model(&User{}).
		Where("username = ?", username).
		First(&userAuth).Error
	if err != nil {
		return "", "", err
	}
	return userAuth.HashedPassword, userAuth.Salt, nil
}

func GetUserInfo(username string) (UserInfo, error) {
	var user User
	err := db.Model(&User{}).
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return UserInfo{}, err
	}
	ret := UserInfo{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
	return ret, nil
}

func IsParticipantInMeeting(meetingID uuid.UUID, userID uint) (bool, error) {
	var count int64
	err := db.Model(&MeetingParticipant{}).
		Where("meeting_id = ?", meetingID).
		Where("user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetParticipantsInMeeting retrieves the participants currently in a meeting from the database.
// The function takes a meeting ID and an optional list of user IDs to ignore.
// It returns a slice of ParticipantInfo structs containing the user ID and username of each participant.
// If an error occurs while querying the database, it is returned along with a nil result.
func GetParticipantsInMeeting(meetingID uuid.UUID, usersToIgnore ...uint) ([]ParticipantInfo, error) {
	var results []ParticipantInfo

	// Start the query joining the users table to get the name
	query := db.Table("meeting_participants").
		Select("meeting_participants.user_id, users.username").
		Joins("join users on users.id = meeting_participants.user_id").
		Where("meeting_participants.meeting_id = ?", meetingID)

	if len(usersToIgnore) > 0 {
		query = query.Where("meeting_participants.user_id NOT IN ?", usersToIgnore)
	}

	// Use Scan instead of Pluck when dealing with structs/multiple columns
	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func GetParticipantCountInMeeting(meetingID uuid.UUID) (int64, error) {
	var count int64
	err := db.Model(&MeetingParticipant{}).
		Where("meeting_id = ?", meetingID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func IsMeetingEmpty(meetingID uuid.UUID) (bool, error) {
	participantCount, err := GetParticipantCountInMeeting(meetingID)
	if err != nil {
		return false, err
	}
	return participantCount == 0, nil
}

func IsHostOfMeeting(meetingID uuid.UUID, userID uint) (bool, error) {
	var count int64
	err := db.Model(&Meeting{}).
		Where("id = ?", meetingID).
		Where("host_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func BanUserFromMeeting(meetingID uuid.UUID, userID uint) error {
	meeting := Meeting{UuidModel: UuidModel{ID: meetingID}}
	user := User{Model: gorm.Model{ID: userID}}
	err := db.
		Model(&meeting).
		Association("BannedUsers").
		Append(&user)
	return err
}

func IsBannedFromMeeting(meetingID uuid.UUID, userID uint) (bool, error) {
	var count int64

	err := db.Table("meeting_banned_users").
		Where("meeting_id = ? AND user_id = ?", meetingID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetUserVideoChunks(meetingID uuid.UUID, userID uint) ([]UserVideoChunk, error) {
	var userVideoChunks []UserVideoChunk

	err := db.
		Where("meeting_id = ? AND user_id = ? AND visited = false", meetingID, userID).
		Order("chunk_number ASC").
		Find(&userVideoChunks).Error

	return userVideoChunks, err
}

func GetLatestStartChunk(meetingID uuid.UUID, userID uint) (*UserVideoChunk, error) {
	var userVideoChunk UserVideoChunk
	err := db.
		Where("meeting_id = ? AND user_id = ? AND chunk_number = 0", meetingID, userID).
		Order("created_at DESC").
		First(&userVideoChunk).Error

	return &userVideoChunk, err
}

func GetKthStartChunk(meetingID uuid.UUID, userID uint, k int) (*UserVideoChunk, error) {
	var userVideoChunk UserVideoChunk
	err := db.
		Where("meeting_id = ? AND user_id = ? AND chunk_number = 0", meetingID, userID).
		Order("created_at ASC").
		Offset(k).
		First(&userVideoChunk).Error

	return &userVideoChunk, err
}

func CountStartChunks(meetingID uuid.UUID, userID uint) (int64, error) {
	var count int64
	err := db.Model(&UserVideoChunk{}).
		Where("meeting_id = ? AND user_id = ? AND chunk_number = 0", meetingID, userID).
		Count(&count).Error
	return count, err
}

func IsFaceDetectionRequired(meetingID uuid.UUID) (bool, error) {
	var count int64
	err := db.Model(&Meeting{}).
		Where("id = ?", meetingID).
		Where("is_face_detection_required = ?", true).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// pipe in all user video chunks created between maxTime and minTime from the database,
// including minTime and excluding maxTime.
// if minTime or maxTime are zero, they are ignored.
func PipeUserVideoChunksBetween(meetingID uuid.UUID, userID uint, minTime time.Time, maxTime time.Time, pipeIn io.Writer) error {
	query := db.Model(&UserVideoChunk{}).
		Where("meeting_id = ? AND user_id = ?", meetingID, userID)
	if !minTime.IsZero() {
		query = query.Where("created_at >= ?", minTime)
	}
	if !maxTime.IsZero() {
		query = query.Where("created_at < ?", maxTime)
	}

	rows, err := query.
		Order("chunk_number ASC").
		Rows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var userVideoChunk UserVideoChunk
		db.ScanRows(rows, &userVideoChunk)
		_, err = pipeIn.Write(userVideoChunk.Chunk)
		if err != nil {
			return err
		}
	}
	return rows.Err()
}

// get ids of all participants who joined the meeting
// including participants who left the meeting
func GetAllMeetingParticipantIDs(meetingID uuid.UUID) ([]uint, error) {
	var meetingParticipants []uint
	participantJoinedEvent := config.GetString("database.meeting_events.participant_joined")
	err := db.
		Model(&MeetingEvent{}).
		Where("meeting_id = ? AND event = ?", meetingID, participantJoinedEvent).
		Distinct().
		Pluck("user_id", &meetingParticipants).Error
	return meetingParticipants, err
}

// get meetings infos from the database, filtering by the given parameters.
// hostName and meetingName can be either a name or an id.
func GetMeetingsInfo(from time.Time, to time.Time, hostName string, meetingName string) ([]MeetingInfo, error) {
	var meetings []MeetingInfo
	hostId, err := strconv.ParseUint(hostName, 10, 32)
	if err != nil {
		hostId = 0
	}
	meetingID, err := uuid.Parse(meetingName)
	if err != nil {
		meetingID = uuid.Nil
	}
	err = db.Model(&Meeting{}).
		Select("meetings.id, meetings.created_at, meetings.deleted_at, meetings.updated_at, meetings.name, meetings.is_face_detection_required, meetings.host_id, users.username as host_username").
		Joins("left join users on users.id = meetings.host_id").
		Where("meetings.created_at BETWEEN ? AND ?", from, to).
		Where("users.username ILIKE ? OR users.id = ?", "%"+hostName+"%", hostId).
		Where("meetings.name ILIKE ? OR meetings.id = ?", "%"+meetingName+"%", meetingID).
		Scan(&meetings).Error
	return meetings, err
}

func GetTranscript(meetingID uuid.UUID, userID uint) (string, error) {
	var transcript string
	err := db.
		Model(&ParticipantTranscription{}).
		Select("transcript").
		Where("meeting_id = ? AND user_id = ?", meetingID, userID).
		Take(&transcript).Error
	return transcript, err
}

func GetUsername(userID uint) (string, error) {
	var username string
	err := db.
		Model(&User{}).
		Select("username").
		Where("id = ?", userID).
		Take(&username).Error
	return username, err
}

func GetSummary(meetingID uuid.UUID) (string, error) {
	var summary string
	err := db.
		Model(&Meeting{}).
		Select("summary").
		Where("id = ?", meetingID).
		Take(&summary).Error
	return summary, err
}

func GetMeetingTranscripts(meetingID uuid.UUID) ([]ParticipantTranscription, error) {
	var transcripts []ParticipantTranscription
	err := db.
		Model(&ParticipantTranscription{}).
		Preload("User").
		Where("meeting_id = ?", meetingID).
		Find(&transcripts).Error
	return transcripts, err
}
