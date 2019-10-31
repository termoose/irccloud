package ui

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/termoose/irccloud/requests"
)

type channel struct {
	layout *tview.Grid

	name   string
	chat   *tview.TextView
	users  *tview.List
	input  *tview.InputField
	info   *tview.TextView
	cid    int
}

type channelList []channel

func (c channelList) String(i int) string {
	return c[i].name
}

func (c channelList) Len() int {
	return len(c)
}

type View struct {
	basePages     *tview.Pages
	pages         *tview.Pages
	app           *tview.Application
	activeChannel int
	channels      channelList
	websocket     *requests.Connection
}

func NewView(socket *requests.Connection) (*View) {
	view := &View{
		pages: tview.NewPages(),
		basePages: tview.NewPages(),
		websocket: socket,
	}

	return view
}

func (v *View) Start() {
	v.app = tview.NewApplication()

	v.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		total_active := v.pages.GetPageCount()

		if event.Key() == tcell.KeyLeft {
			v.activeChannel = (v.activeChannel - 1 + total_active) % total_active
			page := v.channels[v.activeChannel]
			v.pages.SwitchToPage(page.name)
			v.app.SetFocus(page.input)
		}

		if event.Key() == tcell.KeyRight {
			v.activeChannel = (v.activeChannel + 1) % total_active
			page := v.channels[v.activeChannel]
			v.pages.SwitchToPage(page.name)
			v.app.SetFocus(page.input)
		}

		if event.Key() == tcell.KeyCtrlSpace {
			if v.basePages.HasPage("select_channel") {
				v.hideChannelSelector()
			} else {
				v.showChannelSelector()
			}
		}

		return event
	})

	// Create "channels" layer at the bottom of `basePages`
	//v.basePages.AddAndSwitchToPage("channels", v.pages, true)
	v.basePages.AddPage("channel", v.pages, true, true)

	if err := v.app.
		SetRoot(v.basePages, true).
		SetFocus(v.basePages).
		Run(); err != nil {
		panic(err)
	}
}

func (v *View) Stop() {
	v.app.Stop()
}

func (v *View) AddChannel(name, topic string, cid int, user_list []string) {
	headerStr := fmt.Sprintf("[gold:-:b]%s[-:-:-]: %s", name, topic)
	new_chan := channel{
		layout: tview.NewGrid().
			SetRows(1, 0, 1).
			SetColumns(20, 0, 20).
			SetBorders(true),
		name: name,
		chat: newTextView(""),
		users: newListView(),
		input: newTextInput(),
		info: newTextView(headerStr),
		cid: cid,
	}

	// Set callback for handling message sending
	new_chan.input.SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				v.sendToBuffer(cid, name, new_chan.input.GetText())
				new_chan.input.SetText("")
			}
	})

	for _, user := range user_list {
		new_chan.users.AddItem(user, user, 0, nil)
	}

	// Layout
	new_chan.layout.AddItem(new_chan.users, 0, 2, 3, 1, 0, 0, false)
	new_chan.layout.AddItem(new_chan.chat,  1, 0, 1, 2, 0, 0, false)
	new_chan.layout.AddItem(new_chan.input, 2, 0, 1, 2, 0, 0, false)
	new_chan.layout.AddItem(new_chan.info,  0, 0, 1, 2, 0, 0, false)

	v.pages.AddAndSwitchToPage(name, new_chan.layout, true)

	// FIXME: Might not be the best idea to keep this counter
	v.activeChannel = v.pages.GetPageCount()

	v.channels = append(v.channels, new_chan)

	v.app.SetFocus(new_chan.input)
	v.app.Draw()
}

func remove(s []channel, i int) []channel {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (v *View) RemoveChannel(name string) {
	v.pages.RemovePage(name)
	v.activeChannel--;

	index, _ := v.getChannel(name)
	v.channels = remove(v.channels, index)

	v.app.Draw()
}

func (v *View) sendToBuffer(cid int, channel, message string) {
	v.websocket.SendMessage(cid, channel, message)
}

func (v *View) getChannel(name string) (int, *channel) {
	for i, c := range v.channels {
		if c.name == name {
			return i, &v.channels[i]
		}
	}

	return 0, nil
}

func (v *View) getUserIndex(channel, name string) (int, *channel, error) {
	_, c := v.getChannel(channel)

	if c != nil {
		list := c.users.FindItems(name, name, true, false)

		for _, elem := range list {
			found_name, _ := c.users.GetItemText(elem)
			if found_name == name {
				return elem, c, nil
			}
		}
	}

	return 0, nil, errors.New("Could not find user and/or channel")
}

func (v *View) ChangeUserNick(channel, oldnick, newnick string) {
	index, c, err := v.getUserIndex(channel, oldnick)

	if err == nil {
		c.users.SetItemText(index, newnick, newnick)
		line := fmt.Sprintf("  [coral]%s[-:-:-] is now known as [gold]%s[-:-:-]", oldnick, newnick)

		v.writeToBuffer(line, c)
		v.app.Draw()
	}
}

func (v *View) AddUser(channel, nick string) {
	_, c := v.getChannel(channel)

	if c != nil {
		c.users.AddItem(nick, nick, 0, nil)
		v.app.Draw()
	}
}

func (v *View) RemoveUser(channel, nick string) {
	index, c, err := v.getUserIndex(channel, nick)

	if err != nil {
		c.users.RemoveItem(index)
		v.app.Draw()
	}
}

func (v *View) AddQuitEvent(channel, nick, hostmask, reason string) {
	_, c := v.getChannel(channel)

	if c != nil {
		line := fmt.Sprintf("[blueviolet]  <- [blueviolet:-:b]%s[-:-:-] quit (%s): [blueviolet]%s[-:-:-]", nick, hostmask, reason)
		v.writeToBuffer(line, c)
	}
}

func (v *View) AddPartEvent(channel, nick, hostmask string) {
	_, c := v.getChannel(channel)

	if c != nil {
		line := fmt.Sprintf("[blueviolet]  <- [blueviolet:-:b]%s[-:-:-] left (%s)",
			nick, hostmask)
		v.writeToBuffer(line, c)
	}
}

func (v *View) AddJoinEvent(channel, nick, hostmask string) {
	_, c := v.getChannel(channel)

	if c != nil {
		line := fmt.Sprintf("[aquamarine]  -> [aquamarine:-:b]%s[-:-:-] joined (%s)",
			nick, hostmask)
		v.writeToBuffer(line, c)
	}
}

func (v *View) AddBufferMsg(channel, from, msg string) {
	_, c := v.getChannel(channel)

	if c != nil {
		line := fmt.Sprintf("<[-:-:b]%s[-:-:-]> %s", from, msg)
		v.writeToBuffer(line, c)
	}
}

func (v *View) writeToBuffer(line string, c *channel) {
	c.chat.Write([]byte("\n" + line))
	c.chat.ScrollToEnd()
	v.app.Draw()
}

func newTextInput() *tview.InputField {
	return tview.NewInputField().
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetPlaceholder("type here...")
}

func newListView() *tview.List {
	return tview.NewList().ShowSecondaryText(false).
		SetSelectedFocusOnly(true)
}

func newTextView(text string) *tview.TextView {
	return tview.NewTextView().
		SetText(text).
		SetDynamicColors(true).
		SetWordWrap(true)
}
