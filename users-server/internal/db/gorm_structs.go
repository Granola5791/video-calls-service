package db

import (
	"time"

	"github.com/google/uuid"
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
	Name                    string `gorm:"not null;default:'none'" json:"name"`
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
