package ui

import (
	"github.com/rivo/tview"
)

type window struct {
	Main *tview.Grid

	chat *tview.TextView
	users *tview.TextView
	input *tview.TextView
	info *tview.TextView
}

type channel struct {
	layout *tview.Grid

	chat *tview.TextView
	users *tview.TextView
	input *tview.TextView
	info *tview.TextView
}

func NewWindow() (*window) {
	window := &window{
		Main: tview.NewGrid().
			SetRows(1, 0, 1).SetColumns(20, 0, 20).SetBorders(true),
		chat: newTextView("lol chat in struct"),
		users: newTextView("users"),
		input: newTextView("input text here"),
		info: newTextView("info at the top!"),
	}

	// Layout
	window.Main.AddItem(window.users, 0, 2, 3, 1, 0, 0, false)
	window.Main.AddItem(window.chat,  1, 0, 1, 2, 0, 0, false)
	window.Main.AddItem(window.input, 2, 0, 1, 2, 0, 0, false)
	window.Main.AddItem(window.info,  0, 0, 1, 2, 0, 0, false)

	return window
}

func (w *window) Run() {
	if err := tview.NewApplication().SetRoot(w.Main, true).Run(); err != nil {
	 	panic(err)
	}
}

func newTextView(text string) *tview.TextView {
	return tview.NewTextView().SetText(text)
}
