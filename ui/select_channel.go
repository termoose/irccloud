package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/sahilm/fuzzy"
)

func (v *View) showChannelSelector() {
	input := tview.NewInputField().
		SetPlaceholder("Select channel").
		SetFieldBackgroundColor(tcell.ColorGold).
		SetFieldTextColor(tcell.ColorBlack)

	input.SetAutocompleteFunc(
		func(currentText string) []string {
			results := fuzzy.FindFrom(currentText, v.channels)
			resultStrs := make([]string, len(results))
			for i, r := range results {
				resultStrs[i] = v.channels[r.Index].name
			}

			// If there's a unique match we switch immediately?
			if len(results) == 1 {
				//_, top_channel := v.getChannel(resultStrs[0])
				//v.gotoPage(top_channel)
			}

			input.SetDoneFunc(func(key tcell.Key) {
				if len(results) > 0 && key == tcell.KeyEnter {
					_, selected := v.getChannel(input.GetText())
					if selected != nil {
						v.gotoPage(selected)
					} else {
						_, first_pick := v.getChannel(resultStrs[0])
						v.gotoPage(first_pick)
					}
				}
			})

			return resultStrs
		})

	v.basePages.AddPage("select_channel", floatingModal(input, 40, 10), true, true)
	v.app.SetFocus(input)
}

func (v *View) hideChannelSelector() {
	v.basePages.RemovePage("select_channel")
}

func (v *View) gotoPage(c *channel) {
	v.hideChannelSelector()
	v.pages.SwitchToPage(c.name)
	v.app.SetFocus(c.input)
}
