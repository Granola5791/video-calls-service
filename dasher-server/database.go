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
	HostID      uint   `gorm:"not null"`
	BannedUsers []User `gorm:"many2many:meeting_banned_users;"`
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
	return nil
}

func SaveVideoChunkToDB(chunk []byte, meetingID uuid.UUID, userID uint, chunkNumber uint) error {
	userVideoChunk := UserVideoChunk{
		UserID:      userID,
		MeetingID:   meetingID,
		Chunk:       chunk,
		ChunkNumber: chunkNumber,
		Visited:     false,
	}
	err := db.Create(&userVideoChunk).Error
	if err != nil {
		return err
	}
	return nil
}
