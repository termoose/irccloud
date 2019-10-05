package main

import (
	_ "github.com/termoose/irccloud/events"
	_ "github.com/termoose/irccloud/requests"
	"github.com/termoose/irccloud/ui"
)

func main() {
	// session := requests.GetSessionToken("birkedal85@gmail.com", "SECRET")

	view := ui.NewView()
	view.AddChannel("lol channel")
	view.AddChannel("another channel")
	view.Start()
	// ws_conn := requests.NewConnection(session)
	//window := ui.NewWindow()
	// event_handler := events.NewHandler(session, window)

	// go func() {
	// 	for {
	// 		msg, err := ws_conn.ReadMessage()

	// 		if err != nil {
	// 			panic("Connection lost!")
	// 		}

	// 		event_handler.Enqueue(msg)
	// 	}
	// }()

	//window.Run()

	// for {
	// 	msg, err := ws_conn.ReadMessage()

	// 	if err != nil {
	// 		panic("Connection lost")
	// 	}

	// 	event_handler.Enqueue(msg)
	// }
}
