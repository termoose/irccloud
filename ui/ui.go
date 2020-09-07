package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/termoose/irccloud/requests"
	"sync"
)

type View struct {
	basePages   *tview.Pages
	pages       *tview.Pages
	layout      *tview.Grid
	app         *tview.Application
	channels    channelList
	websocket   *requests.Connection
	Activity    *activityBar
	lastChan    string
	channelLock sync.Mutex
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

func NewView(socket *requests.Connection, triggerWords []string, lastChannel string) *View {
	view := &View{
		pages:     tview.NewPages(),
		layout:    newGrid(),
		basePages: tview.NewPages(),
		websocket: socket,
		Activity:  NewActivityBar(triggerWords),
		lastChan:  lastChannel,
	}

	return view
}

func (v *View) GetCurrentChannel() string {
	name, _ := v.pages.GetFrontPage()
	return name
}

func (v *View) Start() {
	v.app = tview.NewApplication()

	v.app.SetBeforeDrawFunc(func(screen tcell.Screen) bool {
		screen.Clear()
		return false
	})

	v.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlSpace {
			if v.basePages.HasPage("select_channel") {
				v.hideChannelSelector()
			} else {
				v.showChannelSelector()
			}
		}

		if event.Key() == tcell.KeyCtrlB {
			lastActive, err := v.Activity.GetLatestActiveChannel()

			if err == nil {
				_, channel := v.getChannelByName(lastActive)

				if channel != nil {
					v.gotoPage(channel)
				}
			}
		}

		if event.Key() == tcell.KeyPgUp {
			channelName := v.GetCurrentChannel()
			_, channel := v.getChannelByName(channelName)

			channel.Scroll(-10)
		}

		if event.Key() == tcell.KeyPgDn {
			channelName := v.GetCurrentChannel()
			_, channel := v.getChannelByName(channelName)

			channel.Scroll(10)
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

func (v *View) SetLatestChannel() {
	_, selected := v.getChannelByName(v.lastChan)

	if selected != nil {
		v.app.QueueUpdate(func() {
			v.Activity.MarkAsVisited(selected.name, v)
			v.pages.SwitchToPage(selected.name)
			v.app.SetFocus(selected.input)
		})
	}
}

func (v *View) Redraw() {
	v.app.Draw()
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

func newGrid() *tview.Grid {
	return tview.NewGrid().
		SetRows(1, 0).
		SetColumns(0)
}

func newTextView(text string) *tview.TextView {
	return tview.NewTextView().
		SetText(text).
		SetDynamicColors(true).
		SetWordWrap(true)
}
