package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HandleStream(c *gin.Context) {
	meetingID := c.Param(GetStringFromConfig("meeting.meeting_id_name"))
	userID, _ := c.Get(GetStringFromConfig("auth_jwt.user_id_name"))

	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.String(http.StatusForbidden, GetStringFromConfig("error.forbidden"))
		return
	}
	defer ws.Close()

	cmd, stdin, err := InitMpegDash(meetingID, userID.(int))
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

	for {
		messageType, data, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		if messageType == websocket.BinaryMessage {
			// Here we do the fun stuff with the received video data
			log.Printf(GetStringFromConfig("stream.got_chunk_msg"), len(data))
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
}
