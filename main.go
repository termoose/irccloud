package main

import (
	"github.com/termoose/irccloud/events"
	"github.com/termoose/irccloud/requests"
	"github.com/rivo/tview"
)

func main() {
	session := requests.GetSessionToken("birkedal85@gmail.com", "SECRET")

	ws_conn := requests.NewConnection(session)
	event_handler := events.NewHandler(session)

	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		panic(err)
	}

	for {
		msg, err := ws_conn.ReadMessage()

		if err != nil {
			panic("Connection lost")
		}

		event_handler.Enqueue(msg)
	}
}
