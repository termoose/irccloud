package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"time"
)

func (v *View) ChangeUserNick(channel, oldnick, newnick string, time int64) {
	index, c, err := v.getUserIndex(channel, oldnick)

	if err == nil {
		ts := getTimestamp(time)
		c.users.SetItemText(index, newnick, newnick)
		line := fmt.Sprintf("[-:-:d]%s[-:-:-]  [coral]%s[-:-:-] is now known as [gold]%s[-:-:-]", ts, oldnick, newnick)

		v.writeToBuffer(line, c)
		v.app.Draw()
	}
}

func (v *View) AddQuitEvent(channel, nick, hostmask, reason string, time int64) {
	_, c := v.getChannel(channel)

	if c != nil {
		ts := getTimestamp(time)
		line := fmt.Sprintf("[-:-:d]%s[-:-:-][blueviolet]  <- [blueviolet:-:b]%s[-:-:-] quit (%s): [blueviolet]%s[-:-:-]", ts, nick, hostmask, reason)
		v.writeToBuffer(line, c)
	}
}

func (v *View) AddPartEvent(channel, nick, hostmask string, time int64) {
	_, c := v.getChannel(channel)

	if c != nil {
		ts := getTimestamp(time)
		line := fmt.Sprintf("[-:-:d]%s[-:-:-][blueviolet]  <- [blueviolet:-:b]%s[-:-:-] left (%s)", ts, nick, hostmask)
		v.writeToBuffer(line, c)
	}
}

func (v *View) AddJoinEvent(channel, nick, hostmask string, time int64) {
	_, c := v.getChannel(channel)

	if c != nil {
		ts := getTimestamp(time)
		line := fmt.Sprintf("[-:-:d]%s[-:-:-][aquamarine]  -> [aquamarine:-:b]%s[-:-:-] joined (%s)", ts, nick, hostmask)
		v.writeToBuffer(line, c)
	}
}

func getTimestamp(t int64) string {
	tm := time.Unix(t/1000000, 0)
	hour, min, _ := tm.Clock()
	return tview.Escape(fmt.Sprintf("[%02d:%02d]", hour, min))
}

func (v *View) AddBufferMsg(channel, from, msg string, time int64) {
	_, c := v.getChannel(channel)

	if c != nil {
		ts := getTimestamp(time)
		line := fmt.Sprintf("[-:-:d]%s[-:-:-] <[-:-:b]%s[-:-:-]> %s", ts, from, msg)
		v.writeToBuffer(line, c)
	}
}

func (v *View) writeToBuffer(line string, c *channel) {
	_, _ = c.chat.Write([]byte("\n" + line))
	c.chat.ScrollToEnd()
	v.app.Draw()
}
