package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/termoose/irccloud/requests"
)

type View struct {
	basePages *tview.Pages
	pages     *tview.Pages
	layout    *tview.Grid
	app       *tview.Application
	channels  channelList
	websocket *requests.Connection
	Activity  *activityBar
}

func floatingModal(p tview.Primitive, width, height int) tview.Primitive {
	return tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(p, height, 1, false).
			AddItem(nil, 0, 1, false), width, 1, false).
		AddItem(nil, 0, 1, false)
}

func NewView(socket *requests.Connection) *View {
	view := &View{
		pages:     tview.NewPages(),
		layout:    tview.NewGrid().
			SetRows(1, 0).
			SetColumns(0),
		basePages: tview.NewPages(),
		websocket: socket,
		Activity:  NewActivityBar(),
	}

	return view
}

func (v *View) Start() {
	v.app = tview.NewApplication()

	v.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlSpace {
			if v.basePages.HasPage("select_channel") {
				v.hideChannelSelector()
			} else {
				v.showChannelSelector()
			}
		}

		return event
	})

	v.basePages.AddPage("channel", v.pages, true, true)
	// v.basePages.AddPage("splash", floatingModal(newANSIView(), 100, 35),
	// 	true, true)

	v.layout.AddItem(v.basePages, 1, 0, 1, 1, 0, 0, true)
	v.layout.AddItem(v.Activity.bar, 0, 0, 1, 1, 0, 0, false)

	if err := v.app.
		SetRoot(v.layout, true).
		SetFocus(v.layout).
		Run(); err != nil {
		panic(err)
	}
}

func (v *View) HideSplash() {
	v.basePages.RemovePage("splash")
}

func (v *View) Stop() {
	v.app.Stop()
}

func (v *View) sendToBuffer(cid int, channel, message string) {
	v.websocket.SendMessage(cid, channel, message)
}

func newTextInput() *tview.InputField {
	return tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorDimGray).
		SetFieldTextColor(tcell.ColorWhite).
		SetPlaceholderTextColor(tcell.ColorWhiteSmoke).
		SetPlaceholder("type here...")
}

func newListView() *tview.List {
	return tview.NewList().
		ShowSecondaryText(false).
		SetSelectedFocusOnly(true).
		SetMainTextColor(tcell.ColorLightSkyBlue)
}

func newANSIView() *tview.TextView {
	art := readFile("test.ans")
	return tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		//	SetText(art)
		SetText(tview.TranslateANSI(art))
}

func newTextView(text string) *tview.TextView {
	return tview.NewTextView().
		SetText(text).
		SetDynamicColors(true).
		SetWordWrap(true)
}
