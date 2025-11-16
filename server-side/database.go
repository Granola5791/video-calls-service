package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"uniqueIndex;not null"`
	HashedPassword string `gorm:"not null"`
	Salt           string `gorm:"not null"`
}

type UserAuth struct {
	HashedPassword string `gorm:"not null"`
	Salt           string `gorm:"not null"`
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

	err = db.AutoMigrate(&User{})
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
