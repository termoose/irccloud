package ui

import (
	"errors"
	"fmt"
	_ "github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"sort"
	"strings"
	"sync"
	"time"
)

type activityBuffer struct {
	displayName  string
	lastActivity time.Time
	visited      bool
	triggerWord  bool
}

type activityBar struct {
	bar          *tview.TextView
	buffers      map[string]activityBuffer
	sorted       []activityBuffer
	triggerWords []string
	buffersLock  sync.Mutex
}

func NewActivityBar(triggers []string) *activityBar {
	return &activityBar{
		bar: tview.NewTextView().
			SetDynamicColors(true).
			SetWrap(false).
			SetScrollable(false),
		buffers:      make(map[string]activityBuffer),
		triggerWords: triggers,
	}
}

func (b *activityBar) updateActivityBar(view *View) {
	b.buffersLock.Lock()
	defer b.buffersLock.Unlock()

	var list []activityBuffer
	for _, b := range b.buffers {
		list = append(list, b)
	}

	sort.SliceStable(list, func(i, j int) bool {
		return list[i].lastActivity.After(list[j].lastActivity)
	})

	b.sorted = list
	b.bar.SetText(generateBar(list))
}

func bufferToBarElement(buffer activityBuffer) string {
	if buffer.triggerWord {
		return fmt.Sprintf("[yellow:blueviolet:b]%s[-:-:-] ", buffer.displayName)
	}

	return fmt.Sprintf("[-:blueviolet:-]%s[-:-:-] ", buffer.displayName)
}

func bufferToVisitedBarElement(buffer activityBuffer) string {
	return fmt.Sprintf("[-:grey:-]%s[-:-:-] ", buffer.displayName)
}

func generateBar(buffers []activityBuffer) string {
	var result string
	for _, b := range buffers {
		if !b.visited {
			result += bufferToBarElement(b)
		} else {
			result += bufferToVisitedBarElement(b)
		}
	}

	return result
}

func (b *activityBar) hasTriggerWord(line string) bool {
	for _, word := range b.triggerWords {
		if strings.Contains(line, word) {
			return true
		}
	}

	return false
}

func (b *activityBar) MarkAsVisited(buffer string, view *View) {
	b.buffersLock.Lock()
	elem, ok := b.buffers[buffer]

	if ok {
		elem.visited = true

		// If you have buffers open with activity from 34 years ago
		// then this is going to look weird! Not sure of a better way
		// to keep the order than a translation down memory lane
		elem.lastActivity = elem.lastActivity.AddDate(-34, 0, 0)

		// Assign the value back to the map
		b.buffers[buffer] = elem
	}
	b.buffersLock.Unlock()

	b.updateActivityBar(view)
}

func (b *activityBar) GetLatestActiveChannel() (string, error) {
	b.buffersLock.Lock()
	defer b.buffersLock.Unlock()

	if len(b.sorted) == 0 {
		return "", errors.New("No channels!")
	}

	return b.sorted[0].displayName, nil
}

func (b *activityBar) RegisterActivity(buffer, msg string, eid int, view *View) {
	channel := view.GetChannel(buffer)
	if channel != nil {
		channel.SetEid(eid)
	}

	// Do not register activity if it's in our current active channel
	if view.GetCurrentChannel() == buffer {
		return
	}

	trigger := b.hasTriggerWord(msg)
	b.buffersLock.Lock()
	b.buffers[buffer] = activityBuffer{
		displayName:  buffer,
		lastActivity: time.Now(),
		visited:      false,
		triggerWord:  trigger,
	}
	b.buffersLock.Unlock()

	b.updateActivityBar(view)
}
