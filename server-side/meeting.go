package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleCreateMeeting(c *gin.Context) {
	meetingID, err := CreateMeetingInDB(c.GetInt("userID"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	jwtKey := []byte(os.Getenv("MEETING_JWT_SECRET"))
	token, err := GenerateMeetingToken(meetingID, jwtKey, GetIntFromConfig("meeting.token_exp"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": GetStringFromConfig("errors.internal_server_error")})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     GetStringFromConfig("meeting.token_name"),
		Value:    token,
		Path:     "/", // visible to all paths
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   GetIntFromConfig("meeting.token_exp"),
	})

	c.String(http.StatusOK, meetingID.String())
}
