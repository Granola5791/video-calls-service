package main

import (
	"context"

	"github.com/gorilla/websocket"
	"log"
)

type ParticipantNotifierStruct struct {
	ID      uint
	channel chan ParticipantNotification
	ctx     context.Context
	cancel  context.CancelFunc
}

func (m *ParticipantNotifierStruct) Init(participantID uint, parentCtx context.Context) {
	m.ID = participantID
	m.ctx, m.cancel = context.WithCancel(parentCtx)
	m.channel = make(chan ParticipantNotification, GetIntFromConfig("notifications.channel_buffer_size"))
}

func (m *ParticipantNotifierStruct) Close() {
	m.cancel()
}

func (m *ParticipantNotifierStruct) Run(ws *websocket.Conn) {
	for {
		select {
		case <-m.ctx.Done():
			ws.Close()
			return
		case notification := <-m.channel:
			if notification.ParticipantID != m.ID {
				err := ws.WriteJSON(notification)
				if err != nil {
					log.Println(err)
					ws.Close()
					return
				}
			}
		}
	}
}
