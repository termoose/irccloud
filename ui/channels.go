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

func (v *View) ChangeTopic(channel, author, newTopic string, time int64) {
	_, c := v.getChannel(channel)

	if c != nil {
		ts := getTimestamp(time)
		c.info.SetText(headerString(channel, newTopic))
		line := fmt.Sprintf("[-:-:d]%s[-:-:-]  [-:-:b]%s[-:-:-] changed topic: [lime:-:-]%s[-:-:-]", ts, author, newTopic)
		v.writeToBuffer(line, c)
		//v.app.Draw()
	}
}

func (v *View) AddChannel(name, topic string, cid int, userList []string) {
	newChan := channel{
		layout: tview.NewGrid().
			SetRows(1, 0, 1).
			SetColumns(20, 0, 20).
			SetBorders(false),
		name:  name,
		chat:  newTextView(""),
		users: newListView(),
		input: newTextInput(),
		info:  newTextView(headerString(name, topic)),
		cid:   cid,
	}

	// Set callback for handling message sending
	newChan.input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			v.sendToBuffer(cid, name, newChan.input.GetText())
			newChan.input.SetText("")
		}
	})

	v.app.QueueUpdateDraw(func() {
		for _, user := range userList {
			newChan.users.AddItem(user, user, 0, nil)
		}

		// Layout
		newChan.layout.AddItem(newChan.users, 0, 2, 3, 1, 0, 0, false)
		newChan.layout.AddItem(newChan.chat, 1, 0, 1, 2, 0, 0, false)
		newChan.layout.AddItem(newChan.input, 2, 0, 1, 2, 0, 0, false)
		newChan.layout.AddItem(newChan.info, 0, 0, 1, 2, 0, 0, false)

		v.pages.AddAndSwitchToPage(name, newChan.layout, true)

		v.channels = append(v.channels, newChan)

		v.app.SetFocus(newChan.input)
	})
}

func remove(s []channel, i int) []channel {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (v *View) RemoveChannel(name string) {
	v.app.QueueUpdateDraw(func() {
		v.pages.RemovePage(name)
	})

	index, _ := v.getChannel(name)
	v.channels = remove(v.channels, index)
}
