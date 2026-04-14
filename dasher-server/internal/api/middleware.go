package api

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/Granola5791/video-calls-service/internal/auth"
)

func RequireAuthentication(c *gin.Context) {

	// get cookie
	tokenName := config.GetStringFromConfig("auth_jwt.token_cookie_name")
	tokenString, err := c.Cookie(tokenName)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := auth.ParseToken(tokenString, jwtKey)
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
	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set(config.GetStringFromConfig("auth_jwt.user_id_name"), int(claims[config.GetStringFromConfig("auth_jwt.user_id_name")].(float64)))
	c.Set("role", claims["role"].(string))

	c.Next()
}

func RequireAuthorizedMeeting(c *gin.Context) {
	// get cookie
	tokenName := config.GetStringFromConfig("meeting.token_name")
	tokenString, err := c.Cookie(tokenName)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token
	jwtKey := []byte(os.Getenv("MEETING_JWT_SECRET"))
	token, err := auth.ParseToken(tokenString, jwtKey)
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

	c.Set(config.GetStringFromConfig("meeting.meeting_id_name"), claims[config.GetStringFromConfig("meeting.meeting_id_name")])

	c.Next()
}

func RequireKeepAliveToken(c *gin.Context) {
	// check route
	if strings.HasPrefix(c.Request.URL.Path, "/"+config.GetStringFromConfig("server.api.create_meeting_path")) {
		return
	}

	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("server.api.params.meeting_id_name")))

	// get cookie
	tokenName := config.GetStringFromConfig("keep_alive.token_cookie_name")
	tokenString, err := c.Cookie(tokenName)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token
	jwtKey := []byte(os.Getenv("KEEP_ALIVE_JWT_SECRET"))
	token, err := auth.ParseToken(tokenString, jwtKey)
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

	// check if token is for the correct meeting
	if claims[config.GetStringFromConfig("jwt.meeting_id_name")].(string) != meetingID.String() {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// check if token is expired
	if claims[config.GetStringFromConfig("jwt.exp_name")].(float64) < float64(time.Now().Unix()) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}

func RequireSameOrigin(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	if origin != config.GetStringFromConfig("server.frontend_addr") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
