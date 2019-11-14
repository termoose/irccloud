package ui

import (
	"fmt"
	_ "github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"time"
	"sort"
)

type activityBuffer struct {
	displayName  string
	lastActivity time.Time
}

type activityBar struct {
	bar     *tview.TextView
	buffers map[string]activityBuffer
}

func NewActivityBar() *activityBar {
	bar := &activityBar{
		bar: tview.NewTextView().
			SetDynamicColors(true).
			SetWrap(false).
			SetScrollable(false),
		//SetText("[-:green:-]lol[-:-:-][-:blue:-]dude[-:-:-]"),
		buffers: make(map[string]activityBuffer),
	}

	return bar
}

func (b *activityBar) updateActivityBar() {
	list := []activityBuffer{}
	for _, b := range b.buffers {
		list = append(list, b)
	}

	sort.SliceStable(list, func(i, j int) bool {
		return list[i].lastActivity.After(list[j].lastActivity)
	})

	reduced := []activityBuffer{}
	if len(list) >= 5 {
		reduced = list[0:5]
	}

	b.bar.SetText(generateBar(reduced))
}

func bufferToBarElement(buffer activityBuffer) string {
	return fmt.Sprintf("[-:blueviolet:-]%s[-:-:-] ", buffer.displayName)
}

func generateBar(buffers []activityBuffer) string {
	var result string
	for _, b := range buffers {
		result += bufferToBarElement(b)
	}

	return result
}

func (b *activityBar) RegisterActivity(buffer string) {
	b.buffers[buffer] = activityBuffer{
		displayName:  buffer,
		lastActivity: time.Now(),
	}

	b.updateActivityBar()
}
