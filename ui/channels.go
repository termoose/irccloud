package ui

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
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

func headerString(name, topic string) string {
	return fmt.Sprintf("[gold:-:b]%s[-:-:-]: [lime:-:-]%s[-:-:-]", name, topic)
}

func (v *View) getChannel(name string) (int, *channel) {
	for i, c := range v.channels {
		if c.name == name {
			return i, &v.channels[i]
		}
	}

	return 0, nil
}

func (v *View) HasChannel(channel string) bool {
	_, c := v.getChannel(channel)

	return c != nil
}

func (v *View) ChangeTopic(channel, author, newtopic string, time int64) {
	_, c := v.getChannel(channel)

	if c != nil {
		ts := getTimestamp(time)
		c.info.SetText(headerString(channel, newtopic))
		line := fmt.Sprintf("[-:-:d]%s[-:-:-]  [-:-:b]%s[-:-:-] changed topic: [lime:-:-]%s[-:-:-]", ts, author, newtopic)
		v.writeToBuffer(line, c)
		v.app.Draw()
	}
}

func (v *View) AddChannel(name, topic string, cid int, user_list []string) {
	new_chan := channel{
		layout: tview.NewGrid().
			SetRows(1, 0, 1).
			SetColumns(20, 0, 20).
			SetBorders(true),
		name:  name,
		chat:  newTextView(""),
		users: newListView(),
		input: newTextInput(),
		info:  newTextView(headerString(name, topic)),
		cid:   cid,
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
	new_chan.layout.AddItem(new_chan.chat, 1, 0, 1, 2, 0, 0, false)
	new_chan.layout.AddItem(new_chan.input, 2, 0, 1, 2, 0, 0, false)
	new_chan.layout.AddItem(new_chan.info, 0, 0, 1, 2, 0, 0, false)

	v.pages.AddAndSwitchToPage(name, new_chan.layout, true)

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

	index, _ := v.getChannel(name)
	v.channels = remove(v.channels, index)

	v.app.Draw()
}
