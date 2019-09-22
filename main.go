package main

import (
	// "github.com/termoose/irccloud/events"
	// "github.com/termoose/irccloud/requests"
	"github.com/rivo/tview"
)

func main() {
	// session := requests.GetSessionToken("birkedal85@gmail.com", "SECRET")

	// ws_conn := requests.NewConnection(session)
	// event_handler := events.NewHandler(session)

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	grid := tview.NewGrid().
		SetRows(1, 0, 1).SetColumns(20, 0, 20).SetBorders(true)
	users := newPrimitive("Users")
	chat := newPrimitive("Chat")
	input := newPrimitive("Input")
	info := newPrimitive("Info")

	grid.AddItem(users, 0, 2, 3, 1, 0, 0, false)
	grid.AddItem(chat,  1, 0, 1, 2, 0, 0, false)
	grid.AddItem(input, 2, 0, 1, 2, 0, 0, false)
	grid.AddItem(info,  0, 0, 1, 2, 0, 0, false)

	if err := tview.NewApplication().SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}

	// for {
	// 	msg, err := ws_conn.ReadMessage()

	// 	if err != nil {
	// 		panic("Connection lost")
	// 	}

	// 	event_handler.Enqueue(msg)
	// }
}
