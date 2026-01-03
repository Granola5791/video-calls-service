package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleCreateMeeting(c *gin.Context) {
	meetingID, _ := c.Get(GetStringFromConfig("meeting.meeting_id_name"))
	
	err := os.MkdirAll(fmt.Sprintf("%s/%d", GetStringFromConfig("meeting.dir_path"), meetingID.(int)), os.FileMode(GetIntFromConfig("meeting.dir_perms")))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusCreated)
}