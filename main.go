package main

import (
	_ "github.com/termoose/irccloud/events"
	_ "github.com/termoose/irccloud/requests"
	_ "github.com/termoose/irccloud/ui"
	"fmt"
	"github.com/termoose/irccloud/config"
)

func main() {
	conf := config.Parse()
	fmt.Printf("%+v\n", conf)
	
	// session := requests.GetSessionToken("birkedal85@gmail.com", "SECRET")
	// view := ui.NewView()
	// ws_conn := requests.NewConnection(session)
	// event_handler := events.NewHandler(session, view)

	// go func() {
	// 	for {
	// 		msg, err := ws_conn.ReadMessage()

	// 		if err != nil {
	// 			panic("Connection lost!")
	// 		}

	// 		event_handler.Enqueue(msg)
	// 	}
	// }()

	// view.Start()
}
