package auth

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/Granola5791/video-calls-service/internal/config"
	"golang.org/x/crypto/argon2"
)

func HashPassword(password string, salt string) string {
	hash := argon2.IDKey(
		[]byte(password),
		[]byte(salt),
		config.GetUint32("hash.time"),
		config.GetUint32("hash.memory"),
		config.GetUint8("hash.threads"),
		config.GetUint32("hash.keyLen"),
	)
	return base64.RawStdEncoding.EncodeToString(hash)
}

func GenerateSalt() (string, error) {
	salt := make([]byte, config.GetInt("hash.saltLen"))
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(salt), nil
}

func VerifyPassword(hashedPassword string, password string, salt string) bool {
	expectedHash := HashPassword(password, salt)
	return expectedHash == hashedPassword
}

func GenerateNewHashedPassword(password string) (string, string, error) {
	salt, err := GenerateSalt()
	if err != nil {
		return "", "", err
	}
	hashedPassword := HashPassword(password, salt)
	return hashedPassword, salt, nil
}
