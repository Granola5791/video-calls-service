package auth

import (
	"time"

	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateJwtToken(claims jwt.MapClaims, jwtKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateMeetingToken(meetingID uuid.UUID, jwtKey []byte, expTimeSec int) (string, error) {
	claims := jwt.MapClaims{
		config.GetStringFromConfig("meeting.meeting_id_name"): meetingID,
		config.GetStringFromConfig("jwt.exp_name"):            time.Now().Add(time.Second * time.Duration(expTimeSec)).Unix(),
	}
	return GenerateJwtToken(claims, jwtKey)
}

func GenerateLoginToken(userID uint, username string, role string, jwtKey []byte, expTimeSec int) (string, error) {
	claims := jwt.MapClaims{
		config.GetStringFromConfig("jwt.user_id_name"): userID,
		config.GetStringFromConfig("jwt.username_name"): username,
		config.GetStringFromConfig("jwt.role_name"):    role,
		config.GetStringFromConfig("jwt.exp_name"):     time.Now().Add(time.Second * time.Duration(expTimeSec)).Unix(),
	}
	return GenerateJwtToken(claims, jwtKey)
}

func GenerateKeepAliveToken(jwtKey []byte, meetingID uuid.UUID, expTimeSec int) (string, error) {
	claims := jwt.MapClaims{
		config.GetStringFromConfig("jwt.meeting_id_name"): meetingID.String(),
		config.GetStringFromConfig("jwt.exp_name"):        time.Now().Add(time.Second * time.Duration(expTimeSec)).Unix(),
	}
	return GenerateJwtToken(claims, jwtKey)
}

func ParseToken(tokenString string, jwtKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	return token, nil
}