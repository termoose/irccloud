package ui

import (
	"errors"
)

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

func (v *View) AddUser(channel, nick string) {
	v.app.QueueUpdate(func() {
		_, c := v.getChannel(channel)

		if c != nil {
			c.users.AddItem(nick, nick, 0, nil)
		}
	})
}

func (v *View) RemoveUser(channel, nick string) {
	v.app.QueueUpdate(func() {
		index, c, err := v.getUserIndex(channel, nick)

		if err == nil {
			c.users.RemoveItem(index)
		}
	})
}
