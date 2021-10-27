package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/termoose/irccloud/config"
	"github.com/termoose/irccloud/events"
	"github.com/termoose/irccloud/requests"
	"github.com/termoose/irccloud/ui"
	"log"
)

func main() {
	// Set this so we don't overwrite the default terminal
	// background color
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault

	conf := config.Parse()
	sessionData, err := requests.GetSessionToken(conf.Username, conf.Password)

	if err != nil {
		log.Println(err)
		return
	}

	wsConn := requests.NewConnection(sessionData)
	view := ui.NewView(wsConn, &conf)

	defer func() {
		current := view.GetCurrentChannel()
		config.WriteLatestChannel(conf, current)
		view.Stop()
	}()

	eventHandler := events.NewHandler(sessionData.APIHost,
		sessionData.Session, view)

	go func() {
		for {
			msg, err := wsConn.ReadMessage()

			if err != nil {
				view.Stop()
				log.Print(err)

				return
			}

			eventHandler.Enqueue(msg)
		}
	}()

	view.Start()
}
