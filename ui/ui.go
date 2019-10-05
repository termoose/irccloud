package ui

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Window struct {
	Main *tview.Grid

	App *tview.Application
	Chat *tview.TextView
	users *tview.TextView
	input *tview.TextView
	info *tview.TextView
}

type channel struct {
	layout *tview.Grid

	name string
	chat *tview.TextView
	users *tview.TextView
	input *tview.TextView
	info *tview.TextView
}

type View struct {
	Pages *tview.Pages
	App *tview.Application
	ActiveChannel int
	Channels []channel
}

func NewView() (*View) {
	view := &View{
		Pages: tview.NewPages(),
	}

	return view
}

func (v *View) Start() {
	v.App = tview.NewApplication()

	v.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		total_active := v.Pages.GetPageCount()

		if event.Key() == tcell.KeyLeft {
			v.ActiveChannel = (v.ActiveChannel - 1 + total_active) % total_active
			page := v.Channels[v.ActiveChannel]
			v.Pages.SwitchToPage(page.name)
		}

		if event.Key() == tcell.KeyRight {
			v.ActiveChannel = (v.ActiveChannel + 1) % total_active
			page := v.Channels[v.ActiveChannel]
			v.Pages.SwitchToPage(page.name)
		}

		return event
	})

	if err := v.App.SetRoot(v.Pages, true).SetFocus(v.Pages).Run(); err != nil {
		panic(err)
	}
}

func (v *View) AddChannel(name string) {
	new_chan := channel{
		layout: tview.NewGrid().SetRows(1, 0, 1).SetColumns(20, 0, 20).SetBorders(true),
		name: name,
		chat: newTextView("text here"),
		users: newTextView("users here"),
		input: newTextView("input text here"),
		info: newTextView(name),
	}

	// Layout
	new_chan.layout.AddItem(new_chan.users, 0, 2, 3, 1, 0, 0, false)
	new_chan.layout.AddItem(new_chan.chat,  1, 0, 1, 2, 0, 0, false)
	new_chan.layout.AddItem(new_chan.input, 2, 0, 1, 2, 0, 0, false)
	new_chan.layout.AddItem(new_chan.info,  0, 0, 1, 2, 0, 0, false)

	v.Pages.AddPage(name, new_chan.layout, true, true)

	// FIXME: Might not be the best idea to keep this counter
	v.ActiveChannel = v.Pages.GetPageCount()

	lol := append(v.Channels, new_chan)
	v.Channels = lol
}

func NewWindow() (*Window) {
	window := &Window{
		Main: tview.NewGrid().
			SetRows(1, 0, 1).SetColumns(20, 0, 20).SetBorders(true),
		Chat: newTextView(""),
		users: newTextView("users"),
		input: newTextView("input text here"),
		info: newTextView("info at the top!"),
	}

	// Layout
	window.Main.AddItem(window.users, 0, 2, 3, 1, 0, 0, false)
	window.Main.AddItem(window.Chat,  1, 0, 1, 2, 0, 0, false)
	window.Main.AddItem(window.input, 2, 0, 1, 2, 0, 0, false)
	window.Main.AddItem(window.info,  0, 0, 1, 2, 0, 0, false)

	return window
}

func (w *Window) AddLine(nick, msg string) {
	line := fmt.Sprintf("<%s> %s\n", nick, msg)
	w.Chat.Write([]byte(line))
	w.App.Draw();
}

func (w *Window) Run() {
	w.App = tview.NewApplication()

	if err := w.App.SetRoot(w.Main, true).Run(); err != nil {
	 	panic(err)
	}
}

func newTextView(text string) *tview.TextView {
	return tview.NewTextView().SetText(text)
}
