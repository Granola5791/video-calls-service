package mywebsocket

import (
	"net/http"

	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/gorilla/websocket"
)

var wsUpgrader websocket.Upgrader

func InitWsUpgrader() {
	wsUpgrader = websocket.Upgrader{
		ReadBufferSize:  config.GetIntFromConfig("stream.read_buffer_size"),
		WriteBufferSize: config.GetIntFromConfig("stream.write_buffer_size"),
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin") == config.GetStringFromConfig("server.frontend_addr")
		},
	}
}

func UpgradeToWebsocket(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	return wsUpgrader.Upgrade(w, r, responseHeader)
}
