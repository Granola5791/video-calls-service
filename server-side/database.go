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
}

type MeetingParticipant struct {
	gorm.Model
	UserID    uint      `gorm:"not null; uniqueIndex:idx_user_meeting"`
	User      User      `gorm:"foreignKey:UserID"`
	MeetingID uuid.UUID `gorm:"not null; uniqueIndex:idx_user_meeting"`
	Meeting   Meeting   `gorm:"foreignKey:MeetingID"`
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

	err = db.AutoMigrate(&User{}, &Meeting{}, &MeetingParticipant{})
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

func CreateMeetingInDB(hostID int) (uuid.UUID, error) {
	meeting := Meeting{}
	err := db.Create(&meeting).Error
	if err != nil {
		return uuid.Nil, err
	}
	return meeting.ID, nil
}

func AddParticipantToMeetingInDB(meetingID uuid.UUID, userID int) error {
	meetingParticipant := MeetingParticipant{
		UserID:    uint(userID),
		MeetingID: meetingID,
	}
	return db.Create(&meetingParticipant).Error
}

func IsParticipantInMeetingInDB(meetingID uuid.UUID, userID int) (bool, error) {
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
