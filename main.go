package main

import (
	"fmt"
	"github.com/termoose/irccloud/events"
	"github.com/termoose/irccloud/requests"
)

func main() {
	session := requests.GetSessionToken("birkedal85@gmail.com", "SECRET")

	ws_conn := requests.NewConnection(session)
	event_handler := events.NewHandler(session)

	for {
		msg, err := ws_conn.ReadMessage()

		if err != nil {
			panic("Connection lost")
		}

		event_handler.Enqueue(msg)
	}
}
