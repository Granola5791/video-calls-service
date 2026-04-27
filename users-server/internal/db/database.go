package db

import (
	"fmt"
	"os"

	"github.com/Granola5791/video-calls-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabaseConnection() error {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		config.GetString("database.sslmode"),
		config.GetString("database.timezone"),
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
		&UserVideoChunk{},
		&ParticipantTranscription{},
	)
	if err != nil {
		return err
	}

	return nil
}
