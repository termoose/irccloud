package ui

import (
	_ "github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"time"
)

type activityBuffer struct {
	displayName  string
	lastActivity time.Time
}

type activityBar struct {
	bar     *tview.TextView
	buffers map[string]activityBuffer
}

func (b *activityBar) RegisterActivity(buffer string) {
	b.buffers[buffer] = activityBuffer{
		displayName:  buffer,
		lastActivity: time.Now(),
	}
}