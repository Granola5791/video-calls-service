package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HandleStream(c *gin.Context) {
	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.String(http.StatusForbidden, GetStringFromConfig("error.forbidden"))
		return
	}
	defer ws.Close()

	cmd, stdin, err := InitMpegDash()
	if err != nil {
		c.String(http.StatusInternalServerError, GetStringFromConfig("error.internal"))
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
	cmd.Wait()
}
