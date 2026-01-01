package main

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string, salt string) string {
	hash := argon2.IDKey(
		[]byte(password),
		[]byte(salt),
		GetUint32FromConfig("hash.time"),
		GetUint32FromConfig("hash.memory"),
		GetUint8FromConfig("hash.threads"),
		GetUint32FromConfig("hash.keyLen"),
	)
	return base64.RawStdEncoding.EncodeToString(hash)
}

func GenerateSalt() (string, error) {
	salt := make([]byte, GetIntFromConfig("hash.saltLen"))
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
