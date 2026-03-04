package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func HandleStream(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(GetStringFromConfig("meeting.meeting_id_name")))
	userID := c.GetInt(GetStringFromConfig("auth_jwt.user_id_name"))

	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.String(http.StatusForbidden, GetStringFromConfig("error.forbidden"))
		return
	}
	defer ws.Close()

	cmd, stdin, err := InitMpegDash(meetingID.String(), uint(userID))
	if err != nil {
		c.String(http.StatusInternalServerError, GetStringFromConfig("error.internal"))
		return
	}
	err = ws.WriteMessage(websocket.TextMessage, []byte(GetStringFromConfig("stream.ready_msg")))
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, GetStringFromConfig("error.internal"))
		return
	}

	chunkNumber := uint(0)
	for {
		messageType, data, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		if messageType == websocket.BinaryMessage {
			// Here we do the fun stuff with the received video data
			go SaveVideoChunkToDB(data, meetingID, uint(userID), uint(chunkNumber))
			chunkNumber++
			PrepareForMpegDash(stdin, data)
		}
	}

	err = stdin.Close()
	if err != nil {
		log.Println(err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Println(err)
	}

	err = os.RemoveAll(fmt.Sprintf("%s/%s/%d", GetStringFromConfig("meeting.dir_path"), meetingID, userID))
	if err != nil {
		log.Println(err)
	}
}