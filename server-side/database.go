package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UuidModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
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
	HostID       uint `gorm:"not null"`
	BannedUsers []User `gorm:"many2many:meeting_banned_users;"`
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

	err = db.AutoMigrate(&User{}, &Meeting{}, &MeetingParticipant{}, &MeetingEvent{})
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

func GetUserIDAndRoleFromDB(username string) (int, string, error) {
	var user User
	err := db.Model(&User{}).
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return 0, "", err
	}
	return int(user.ID), user.Role, nil
}

func CreateMeetingInDB(hostID uint) (uuid.UUID, error) {
	meeting := Meeting{HostID: hostID}
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

func GetMeetingParticipantIDsFromDB(meetingID uuid.UUID, usersToIgnore ...uint) ([]uint, error) {
	var ids []uint

	query := db.
		Model(&MeetingParticipant{}).
		Where("meeting_id = ?", meetingID)

	if len(usersToIgnore) > 0 {
		query = query.Where("user_id NOT IN ?", usersToIgnore)
	}

	err := query.Pluck("user_id", &ids).Error
	if err != nil {
		return nil, err
	}

	return ids, nil
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

func DeleteMeetingFromDB(meetingID uuid.UUID) error {
	return db.Where("id = ?", meetingID).
		Delete(&Meeting{}).Error
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