package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UuidModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type User struct {
	gorm.Model
	Username       string `gorm:"uniqueIndex;not null"`
	Role           string `gorm:"not null;default:user"`
	HashedPassword string `gorm:"not null"`
	Salt           string `gorm:"not null"`
}

type Meeting struct {
	UuidModel
	HostID                  uint   `gorm:"not null" json:"host_id"`
	IsFaceDetectionRequired bool   `gorm:"not null;default:false" json:"is_face_detection_required"`
	BannedUsers             []User `gorm:"many2many:meeting_banned_users;" json:"banned_users"`
	Summary                 string `gorm:"not null;default:''" json:"summary"`
}

type MeetingParticipant struct {
	gorm.Model
	UserID    uint      `gorm:"not null; uniqueIndex:idx_user_meeting"`
	User      User      `gorm:"foreignKey:UserID"`
	MeetingID uuid.UUID `gorm:"not null; uniqueIndex:idx_user_meeting"`
	Meeting   Meeting   `gorm:"foreignKey:MeetingID"`
}

type MeetingEvent struct {
	gorm.Model
	MeetingID uuid.UUID `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	Event     string    `gorm:"not null"`
}

type UserAuth struct {
	HashedPassword string `gorm:"not null"`
	Salt           string `gorm:"not null"`
}

type UserRole struct {
	ID   uint
	Role string `gorm:"not null"`
}

type UserVideoChunk struct {
	gorm.Model
	UserID      uint      `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID"`
	MeetingID   uuid.UUID `gorm:"not null"`
	Meeting     Meeting   `gorm:"foreignKey:MeetingID"`
	Chunk       []byte    `gorm:"not null"`
	ChunkNumber uint      `gorm:"not null"`
	Visited     bool      `gorm:"not null"`
}

type ParticipantTranscription struct {
	gorm.Model
	UserID     uint      `gorm:"not null; uniqueIndex:idx_user_meeting"`
	User       User      `gorm:"foreignKey:UserID"`
	MeetingID  uuid.UUID `gorm:"not null; uniqueIndex:idx_user_meeting"`
	Meeting    Meeting   `gorm:"foreignKey:MeetingID"`
	Transcript string    `gorm:"not null"`
}

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
}

var db *gorm.DB

func InitDatabaseConnection() error {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		GetStringFromConfig("database.sslmode"),
		GetStringFromConfig("database.timezone"),
	)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&User{},
		&Meeting{},
		&MeetingParticipant{},
		&MeetingEvent{},
		&UserAuth{},
		&UserRole{},
		&UserVideoChunk{},
		&ParticipantTranscription{},
	)
	if err != nil {
		return err
	}

	return nil
}

func UserExistsInDB(username string) (bool, error) {
	var count int64
	err := db.Model(&User{}).
		Where("username = ?", username).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func InsertUserToDB(useranme string, hashedPassword string, salt string) error {
	user := User{
		Username:       useranme,
		HashedPassword: hashedPassword,
		Salt:           salt,
	}
	return db.Create(&user).Error
}

func GetUserAuthFromDB(username string) (string, string, error) {
	var userAuth UserAuth
	err := db.Model(&User{}).
		Where("username = ?", username).
		First(&userAuth).Error
	if err != nil {
		return "", "", err
	}
	return userAuth.HashedPassword, userAuth.Salt, nil
}

func GetUserInfoFromDB(username string) (UserInfo, error) {
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

func CreateMeetingInDB(hostID uint, isFaceDetectionRequired bool) (uuid.UUID, error) {
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

func AddParticipantToMeetingInDB(meetingID uuid.UUID, userID uint) error {
	meetingParticipant := MeetingParticipant{
		UserID:    userID,
		MeetingID: meetingID,
	}
	return db.Create(&meetingParticipant).Error
}

func RemoveParticipantFromMeetingInDB(meetingID uuid.UUID, userID uint) error {
	return db.
		Unscoped().
		Where("meeting_id = ? AND user_id = ?", meetingID, userID).
		Delete(&MeetingParticipant{}).Error
}

func IsParticipantInMeetingInDB(meetingID uuid.UUID, userID uint) (bool, error) {
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

// GetParticipantsInMeetingFromDB retrieves the participants currently in a meeting from the database.
// The function takes a meeting ID and an optional list of user IDs to ignore.
// It returns a slice of ParticipantInfo structs containing the user ID and username of each participant.
// If an error occurs while querying the database, it is returned along with a nil result.
func GetParticipantsInMeetingFromDB(meetingID uuid.UUID, usersToIgnore ...uint) ([]ParticipantInfo, error) {
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

func MeetingExistsInDB(meetingID uuid.UUID) (bool, error) {
	var count int64
	err := db.Model(&Meeting{}).
		Where("id = ?", meetingID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetParticipantCountInMeetingInDB(meetingID uuid.UUID) (int64, error) {
	var count int64
	err := db.Model(&MeetingParticipant{}).
		Where("meeting_id = ?", meetingID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func IsMeetingEmptyInDB(meetingID uuid.UUID) (bool, error) {
	participantCount, err := GetParticipantCountInMeetingInDB(meetingID)
	if err != nil {
		return false, err
	}
	return participantCount == 0, nil
}

func RemoveAllMeetingParticipantsFromDB(meetingID uuid.UUID) error {
	return db.
		Unscoped().
		Where("meeting_id = ?", meetingID).
		Delete(&MeetingParticipant{}).Error
}

func LogEventToDB(meetingID uuid.UUID, userID uint, event string) error {
	meetingEvent := MeetingEvent{
		MeetingID: meetingID,
		UserID:    userID,
		Event:     event,
	}
	return db.Create(&meetingEvent).Error
}

func IsHostOfMeetingInDB(meetingID uuid.UUID, userID uint) (bool, error) {
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

func BanUserFromMeetingInDB(meetingID uuid.UUID, userID uint) error {
	meeting := Meeting{UuidModel: UuidModel{ID: meetingID}}
	user := User{Model: gorm.Model{ID: userID}}
	err := db.
		Model(&meeting).
		Association("BannedUsers").
		Append(&user)
	return err
}

func IsBannedFromMeetingInDB(meetingID uuid.UUID, userID uint) (bool, error) {
	var count int64

	err := db.Table("meeting_banned_users").
		Where("meeting_id = ? AND user_id = ?", meetingID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func GetUserVideoChunksFromDB(meetingID uuid.UUID, userID uint) ([]UserVideoChunk, error) {
	var userVideoChunks []UserVideoChunk

	err := db.
		Where("meeting_id = ? AND user_id = ? AND visited = false", meetingID, userID).
		Order("chunk_number ASC").
		Find(&userVideoChunks).Error

	return userVideoChunks, err
}

// mark user video chunks in the range [minChunkNumber, maxChunkNumber] inclusive as visited
func MarkUserVideoChunksAsVisitedInDB(meetingID uuid.UUID, userID uint, minChunkNumber uint, maxChunkNumber uint) error {
	return db.
		Model(&UserVideoChunk{}).
		Where("meeting_id = ? AND user_id = ? AND chunk_number >= ? AND chunk_number <= ?", meetingID, userID, minChunkNumber, maxChunkNumber).
		Update("visited", true).Error
}

func GetLatestStartChunkFromDB(meetingID uuid.UUID, userID uint) (*UserVideoChunk, error) {
	var userVideoChunk UserVideoChunk
	err := db.
		Where("meeting_id = ? AND user_id = ? AND chunk_number = 0", meetingID, userID).
		Order("created_at DESC").
		First(&userVideoChunk).Error

	return &userVideoChunk, err
}

func GetKthStartChunkFromDB(meetingID uuid.UUID, userID uint, k int) (*UserVideoChunk, error) {
	var userVideoChunk UserVideoChunk
	err := db.
		Where("meeting_id = ? AND user_id = ? AND chunk_number = 0", meetingID, userID).
		Order("created_at ASC").
		Offset(k).
		First(&userVideoChunk).Error

	return &userVideoChunk, err
}

func CountStartChunksFromDB(meetingID uuid.UUID, userID uint) (int64, error) {
	var count int64
	err := db.Model(&UserVideoChunk{}).
		Where("meeting_id = ? AND user_id = ? AND chunk_number = 0", meetingID, userID).
		Count(&count).Error
	return count, err
}

func isFaceDetectionRequiredInDB(meetingID uuid.UUID) (bool, error) {
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
func PipeUserVideoChunksBetweenFromDB(meetingID uuid.UUID, userID uint, minTime time.Time, maxTime time.Time, pipeIn io.Writer) error {
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
func GetAllMeetingParticipantIDsFromDB(meetingID uuid.UUID) ([]uint, error) {
	var meetingParticipants []uint
	participantJoinedEvent := GetStringFromConfig("database.meeting_events.participant_joined")
	err := db.
		Model(&MeetingEvent{}).
		Where("meeting_id = ? AND event = ?", meetingID, participantJoinedEvent).
		Distinct().
		Pluck("user_id", &meetingParticipants).Error
	return meetingParticipants, err
}

func InsertTranscriptionToDB(meetingID uuid.UUID, userID uint, transcription string) error {
	return db.
		Create(&ParticipantTranscription{
			MeetingID:  meetingID,
			UserID:     userID,
			Transcript: transcription,
		}).Error
}

func GetAllMeetingInfosFromDB(from time.Time, to time.Time, hostName string) ([]MeetingInfo, error) {
	var meetings []MeetingInfo
	err := db.Model(&Meeting{}).
		Select("meetings.id, meetings.created_at, meetings.deleted_at, meetings.updated_at, meetings.is_face_detection_required, meetings.host_id, users.username as host_username").
		Joins("left join users on users.id = meetings.host_id").
		Where("meetings.created_at BETWEEN ? AND ?", from, to).
		Where("users.username ILIKE ?", "%"+hostName+"%").
		Scan(&meetings).Error
	return meetings, err
}

func GetTranscriptFromDB(meetingID uuid.UUID, userID uint) (string, error) {
	var transcript string
	err := db.
		Model(&ParticipantTranscription{}).
		Select("transcript").
		Where("meeting_id = ? AND user_id = ?", meetingID, userID).
		Take(&transcript).Error
	return transcript, err
}

func GetUsernameFromDB(userID uint) (string, error) {
	var username string
	err := db.
		Model(&User{}).
		Select("username").
		Where("id = ?", userID).
		Take(&username).Error
	return username, err
}

func UpdateSummaryToDB(meetingID uuid.UUID, summary string) error {
	return db.
		Model(&Meeting{}).
		Where("id = ?", meetingID).
		Update("summary", summary).Error
}

func GetSummaryFromDB(meetingID uuid.UUID) (string, error) {
	var summary string
	err := db.
		Model(&Meeting{}).
		Select("summary").
		Where("id = ?", meetingID).
		Take(&summary).Error
	return summary, err
}
