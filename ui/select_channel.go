package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/sahilm/fuzzy"
)

func (v *View) inputDone(key tcell.Key, resultStrs []string, input *tview.InputField) {
	if len(resultStrs) > 0 && key == tcell.KeyEnter {
		_, selected := v.getChannelByName(input.GetText())
		if selected != nil {
			v.gotoPage(selected)
		} else {
			_, first_pick := v.getChannelByName(resultStrs[0])

			if first_pick != nil {
				v.gotoPage(first_pick)
			}
		}
	}
}

func (v *View) showChannelSelector() {
	input := tview.NewInputField().
		SetPlaceholder("Select channel").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldTextColor(tcell.ColorWhite)
		//SetFieldTextColor(tcell.ColorBlack).SetFieldTextColor(tcell.ColorBlack)

	input.SetAutocompleteFunc(
		func(currentText string) []string {
			results := fuzzy.FindFrom(currentText, v.channels)
			resultStrs := make([]string, len(results))
			for i, r := range results {
				resultStrs[i] = v.channels[r.Index].name
			}

			// If we don't have any results we go to the channel
			// with most recent activity
			if len(resultStrs) == 0 {
				last, err := v.Activity.GetLatestActiveChannel()
				if err == nil {
					resultStrs = append(resultStrs, last)
				}
			}

			input.SetDoneFunc(func(key tcell.Key) {
				v.inputDone(key, resultStrs, input)
			})

			return resultStrs
		})

	v.basePages.AddPage("select_channel", floatingModal(input, 40, 10), true, true)
	v.app.SetFocus(input)
}

func (v *View) hideChannelSelector() {
	v.basePages.RemovePage("select_channel")

	page, primitive := v.pages.GetFrontPage()

	if primitive != nil {
		_, channel := v.getChannelByName(page)

		if channel != nil {
			v.app.SetFocus(channel.input)
		}
	}
}

func (v *View) gotoPage(c *channel) {
	v.Activity.MarkAsVisited(c.name, v)
	v.hideChannelSelector()
	v.pages.SwitchToPage(c.name)
	v.app.SetFocus(c.input)
}
