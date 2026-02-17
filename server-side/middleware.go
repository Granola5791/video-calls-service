package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func RequireAuthentication(c *gin.Context) {

	// get cookie
	tokenName := GetStringFromConfig("jwt.token_cookie_name")
	tokenString, err := c.Cookie(tokenName)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
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

	// check if token is expired
	if claims[GetStringFromConfig("jwt.exp_name")].(float64) < float64(time.Now().Unix()) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(GetStringFromConfig("jwt.user_id_name"), int(claims[GetStringFromConfig("jwt.user_id_name")].(float64)))
	c.Set(GetStringFromConfig("jwt.role_name"), claims[GetStringFromConfig("jwt.role_name")].(string))

	c.Next()
}

// can only be called after RequireAuthentication
func RequireAdmin(c *gin.Context) {
	role, _ := c.Get(GetStringFromConfig("jwt.role_name"))
	roleString, ok := role.(string)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if !strings.Contains(roleString, GetStringFromConfig("jwt.admin_role")) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}

func RequireKeepAliveToken(c *gin.Context) {

	// get cookie
	tokenName := GetStringFromConfig("keep_alive.token_cookie_name")
	tokenString, err := c.Cookie(tokenName)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token
	jwtKey := []byte(os.Getenv("KEEP_ALIVE_JWT_SECRET"))
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

	// check if token is expired
	if claims[GetStringFromConfig("jwt.exp_name")].(float64) < float64(time.Now().Unix()) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}

// can only be called after RequireAuthentication
func RequireHost(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))
	userID := c.GetInt(GetStringFromConfig("jwt.user_id_name"))
	isHost, err := IsHostOfMeetingInDB(meetingID, uint(userID))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !isHost {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func RequireNotBanned(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(GetStringFromConfig("server.api.params.meeting_id_name")))
	userID := c.GetInt(GetStringFromConfig("jwt.user_id_name"))
	isBanned, err := IsBannedFromMeetingInDB(meetingID, uint(userID))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if isBanned {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}