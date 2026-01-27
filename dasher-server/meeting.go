package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func HandleCreateMeeting(c *gin.Context) {
	meetingID, _ := c.Get(GetStringFromConfig("meeting.meeting_id_name"))

	err := os.MkdirAll(fmt.Sprintf("%s/%s", GetStringFromConfig("meeting.dir_path"), meetingID), os.FileMode(GetIntFromConfig("meeting.dir_perms")))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	c.Status(http.StatusCreated)
}

func JoinMeeting(meetingID string, userID int) error {
	path := fmt.Sprintf("%s/%s/%d", GetStringFromConfig("meeting.dir_path"), meetingID, userID)
	err := os.Mkdir(path, os.FileMode(GetIntFromConfig("meeting.dir_perms")))
	if err != nil {
		return err
	}
	return nil
}

func HandleJoinMeeting(c *gin.Context) {
	userID, _ := c.Get(GetStringFromConfig("auth_jwt.user_id_name"))
	meetingID := c.Param(GetStringFromConfig("meeting.meeting_id_name"))

	err := JoinMeeting(meetingID, userID.(int))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	c.Status(http.StatusCreated)
}
