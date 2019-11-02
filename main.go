package main

import (
	"github.com/termoose/irccloud/config"
	"github.com/termoose/irccloud/events"
	"github.com/termoose/irccloud/requests"
	"github.com/termoose/irccloud/ui"
	"log"
)

func main() {
	conf := config.Parse()

	session := requests.GetSessionToken(conf.Username, conf.Password)

	ws_conn := requests.NewConnection(session)
	view := ui.NewView(ws_conn)
	defer view.Stop()

	event_handler := events.NewHandler(session, view)

	go func() {
		for {
			msg, err := ws_conn.ReadMessage()

			if err != nil {
				view.Stop()
				log.Fatal(err)
			}

			event_handler.Enqueue(msg)
		}
	}()

	view.Start()
}
