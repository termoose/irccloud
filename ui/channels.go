package ui

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"strings"
)

type channel struct {
	layout *tview.Grid
	name   string
	chat   *tview.TextView
	users  *tview.List
	input  *tview.InputField
	info   *tview.TextView
	cid    int
	bid    int
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

func (v *View) getChannelByName(name string) (int, *channel) {
	for i, c := range v.channels {
		if c.name == name {
			return i, &v.channels[i]
		}
	}

	return 0, nil
}

func (v *View) getChannel(name string, bid int) (int, *channel) {
	for i, c := range v.channels {
		if c.name == name && c.bid == bid {
			return i, &v.channels[i]
		}
	}

	return 0, nil
}

func (v *View) HasChannel(channel string) bool {
	_, c := v.getChannelByName(channel)

	return c != nil
}

func (v *View) ChangeTopic(channel, author, newTopic string, time int64, bid int) {
	_, c := v.getChannel(channel, bid)

	if c != nil {
		ts := getTimestamp(time)
		c.info.SetText(headerString(channel, newTopic))
		line := fmt.Sprintf("[-:-:d]%s[-:-:-]  [-:-:b]%s[-:-:-] changed topic: [lime:-:-]%s[-:-:-]", ts, author, newTopic)
		v.writeToBuffer(line, c)
		//v.app.Draw()
	}
}

func (v *View) AddChannel(name, topic string, cid, bid int, userList []string) {
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
		bid:   bid,
	}

	// Set callback for handling message sending
	newChan.input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			v.sendToBuffer(cid, name, newChan.input.GetText())
			newChan.input.SetText("")
		}

		if key == tcell.KeyTab {
			currText := newChan.input.GetText()
			words := strings.Split(currText, " ")

			if len(words) > 0 {
				lastWord := words[len(words) - 1]
				foundUsersIdxs := newChan.users.FindItems(lastWord, lastWord, false, true)

				if len(foundUsersIdxs) > 0 {
					foundUser, _ := newChan.users.GetItemText(foundUsersIdxs[0])
					newInput := strings.Join(append(words[:len(words) - 1], foundUser), " ")
					newChan.input.SetText(newInput)
				}
			}
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

	index, _ := v.getChannelByName(name)
	v.channels = remove(v.channels, index)
}
