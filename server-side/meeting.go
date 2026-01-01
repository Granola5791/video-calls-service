package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleCreateMeeting(c *gin.Context) {
	meetingID, err := CreateMeetingInDB(c.GetInt("userID"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{GetStringFromConfig("meeting.meeting_id_name"): meetingID})
}
