package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ParseToken(tokenString string, jwtKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func RequireAuthorizedMeeting(c *gin.Context) {
	// get cookie
	tokenName := GetStringFromConfig("meeting.token_name")
	tokenString, err := c.Cookie(tokenName)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token
	jwtKey := []byte(os.Getenv("MEETING_JWT_SECRET"))
	token, err := ParseToken(tokenString, jwtKey)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims["exp"].(float64) < float64(jwt.NewNumericDate(time.Now()).Unix()) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(GetStringFromConfig("meeting.meeting_id_name"), int(claims[GetStringFromConfig("meeting.meeting_id_name")].(float64)))

	c.Next()
}
