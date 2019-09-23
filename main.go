package main

import (
	// "github.com/termoose/irccloud/events"
	// "github.com/termoose/irccloud/requests"
	"github.com/termoose/irccloud/ui"
)

func main() {
	// session := requests.GetSessionToken("birkedal85@gmail.com", "SECRET")

	// ws_conn := requests.NewConnection(session)
	// event_handler := events.NewHandler(session)

	window := ui.NewWindow()
	window.Run()

	// for {
	// 	msg, err := ws_conn.ReadMessage()

	// 	if err != nil {
	// 		panic("Connection lost")
	// 	}

	// 	event_handler.Enqueue(msg)
	// }
}
