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

func GenerateJwtToken(claims jwt.MapClaims, jwtKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func GenerateMeetingToken(meetingID uuid.UUID, jwtKey []byte, expTimeSec int) (string, error) {
	claims := jwt.MapClaims{
		GetStringFromConfig("meeting.meeting_id_name"): meetingID,
		GetStringFromConfig("jwt.exp_name"):            time.Now().Add(time.Second * time.Duration(expTimeSec)).Unix(),
	}
	return GenerateJwtToken(claims, jwtKey)
}

func GenerateLoginToken(userID int, role string, jwtKey []byte, expTimeSec int) (string, error) {
	claims := jwt.MapClaims{
		GetStringFromConfig("jwt.user_id_name"): userID,
		GetStringFromConfig("jwt.role_name"):    role,
		GetStringFromConfig("jwt.exp_name"):     time.Now().Add(time.Second * time.Duration(expTimeSec)).Unix(),
	}
	return GenerateJwtToken(claims, jwtKey)
}

func GenerateKeepAliveToken(jwtKey []byte, meetingID uuid.UUID, expTimeSec int) (string, error) {
	claims := jwt.MapClaims{
		GetStringFromConfig("jwt.meeting_id_name"): meetingID,
		GetStringFromConfig("jwt.exp_name"):        time.Now().Add(time.Second * time.Duration(expTimeSec)).Unix(),
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
