package stream

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/Granola5791/video-calls-service/internal/dash"
	"github.com/Granola5791/video-calls-service/internal/mywebsocket"
	"github.com/Granola5791/video-calls-service/internal/db"
)

func HandleStream(c *gin.Context) {
	meetingID := uuid.MustParse(c.Param(config.GetStringFromConfig("meeting.meeting_id_name")))
	userID := c.GetInt(config.GetStringFromConfig("auth_jwt.user_id_name"))

	ws, err := mywebsocket.UpgradeToWebsocket(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		c.String(http.StatusForbidden, config.GetStringFromConfig("error.forbidden"))
		return
	}
	defer ws.Close()

	cmd, stdin, err := dash.InitMpegDash(meetingID.String(), uint(userID))
	if err != nil {
		c.String(http.StatusInternalServerError, config.GetStringFromConfig("error.internal"))
		return
	}
	err = ws.WriteMessage(websocket.TextMessage, []byte(config.GetStringFromConfig("stream.ready_msg")))
	if err != nil {
		log.Println(err)
		c.String(http.StatusInternalServerError, config.GetStringFromConfig("error.internal"))
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
			go db.SaveVideoChunkToDB(data, meetingID, uint(userID), uint(chunkNumber))
			chunkNumber++
			dash.PrepareForMpegDash(stdin, data)
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

	err = os.RemoveAll(fmt.Sprintf("%s/%s/%d", config.GetStringFromConfig("meeting.dir_path"), meetingID, userID))
	if err != nil {
		log.Println(err)
	}
}