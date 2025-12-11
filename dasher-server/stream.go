package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  GetIntFromConfig("stream.read_buffer_size"),
	WriteBufferSize: GetIntFromConfig("stream.write_buffer_size"),
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == GetStringFromConfig("server.frontend_addr")
	},
}

func HandleStream(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.String(http.StatusForbidden, GetStringFromConfig("error.forbidden"))
		return
	}
	defer ws.Close()

	for {
		messageType, data, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}

		if messageType == websocket.BinaryMessage {
			// Here we do the fun stuff with the received video data
		}
	}
}
