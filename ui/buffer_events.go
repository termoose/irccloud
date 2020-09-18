package ui

import (
	"fmt"
	"github.com/rivo/tview"
	"regexp"
	"strings"
	"time"
)

func (v *View) ChangeUserNick(channel, oldnick, newnick string, time int64, bid int) {
	v.app.QueueUpdateDraw(func() {
		index, c, err := v.getUserIndex(channel, oldnick, bid)

		if err == nil {
			ts := getTimestamp(time)
			c.users.SetItemText(index, newnick, newnick)
			line := fmt.Sprintf("[-:-:d]%s[-:-:-]  [coral]%s[-:-:-] is now known as [gold]%s[-:-:-]", ts, tview.Escape(oldnick), tview.Escape(newnick))

			v.writeToBuffer(line, c)
		}
	})
}

func (v *View) AddQuitEvent(channel, nick, hostmask, reason string, time int64, bid int) {
	v.app.QueueUpdateDraw(func() {
		_, c := v.getChannel(channel, bid)

		if c != nil {
			ts := getTimestamp(time)
			line := fmt.Sprintf("[-:-:d]%s[-:-:-][blueviolet]  <- [blueviolet:-:b]%s[-:-:-] quit (%s): [blueviolet]%s[-:-:-]", ts, tview.Escape(nick), hostmask, tview.Escape(reason))
			v.writeToBuffer(line, c)
		}
	})
}

func (v *View) AddPartEvent(channel, nick, hostmask string, time int64, bid int) {
	v.app.QueueUpdateDraw(func() {
		_, c := v.getChannel(channel, bid)

		if c != nil {
			ts := getTimestamp(time)
			line := fmt.Sprintf("[-:-:d]%s[-:-:-][blueviolet]  <- [blueviolet:-:b]%s[-:-:-] left (%s)", ts, tview.Escape(nick), hostmask)
			v.writeToBuffer(line, c)
		}
	})
}

func (v *View) AddJoinEvent(channel, nick, hostmask string, time int64, bid int) {
	v.app.QueueUpdateDraw(func() {
		_, c := v.getChannel(channel, bid)

		if c != nil {
			ts := getTimestamp(time)
			line := fmt.Sprintf("[-:-:d]%s[-:-:-][aquamarine]  -> [aquamarine:-:b]%s[-:-:-] joined (%s)", ts, tview.Escape(nick), hostmask)
			v.writeToBuffer(line, c)
		}
	})
}

func getTimestamp(t int64) string {
	tm := time.Unix(t/1000000, 0)
	hour, min, _ := tm.Clock()
	return tview.Escape(fmt.Sprintf("[%02d:%02d]", hour, min))
}

func (v *View) AddBufferMsg(channel, from, msg string, time int64, bid int) {
	v.app.QueueUpdateDraw(func() {
		_, c := v.getChannel(channel, bid)

		if c != nil {
			ts := getTimestamp(time)
			line := fmt.Sprintf("[-:-:d]%s[-:-:-] <[-:-:b]%s[-:-:-]> %s", ts, tview.Escape(from), tview.Escape(msg))

			line = tagifyCodeQuotes(line, `[grey:black:b]`, `[-:-:-]`)
			v.writeToBuffer(line, c)
		}
	})
}

func (v *View) writeToBuffer(line string, c *channel) {
	c.SendToChannel(line)
	//_, _ = c.chat.Write([]byte("\n" + line))
	//c.chat.ScrollToEnd()
}

func trimString(s string) string {
	if len(s) < 2 {
		return s
	}

	return s[1 : len(s)-1]
}

func tagifyCodeQuotes(input, startTag, endTag string) string {
	regex := regexp.MustCompile("`(.*?)`")
	found := regex.FindAllString(input, -1)

	for _, f := range found {
		tag := fmt.Sprintf("%s%s%s", startTag, trimString(f), endTag)
		input = strings.Replace(input, f, tag, -1)
	}

	return input
}