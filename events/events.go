package events

import (
	"encoding/json"
	_ "fmt"
	"github.com/termoose/irccloud/requests"
	"github.com/termoose/irccloud/ui"
	_ "log"
)

type eventHandler struct {
	Queue chan eventData
	SessionToken string
	Window *ui.View
}

func NewHandler(token string, w *ui.View) (*eventHandler) {
	handler := &eventHandler{
		Queue: make(chan eventData, 8),
		SessionToken: token,
		Window: w,
	}

	// Start consumer thread
	go func() {
		for curr_event := range handler.Queue {
			handler.handle(curr_event, false)
		}
	}()

	return handler
}

func (e *eventHandler) Enqueue(msg []byte) {
	current := eventData{}
	json.Unmarshal(msg, &current)

	// Attach raw message data
	current.Data = msg

	e.Queue <- current
}

func (e *eventHandler) queueEvent(event eventData) {
	e.Queue <- event
}

func (e *eventHandler) handleBacklog(url string) {
	backlog := requests.GetBacklog(e.SessionToken, url)
	backlogData := parseBacklog(backlog)

	// First we initialize all channels
	for _, event := range backlogData {
		if event.Type == "channel_init" {
			user_strings := []string{}
			for _, user_string := range event.Members {
				user_strings = append(user_strings, user_string.Nick)
			}

			e.Window.AddChannel(event.Chan, event.Topic.Text, event.Cid, user_strings)
		}
	}

	// Then we fill them with the message backlog, should we send these events
	// to the event queue to have them arrive before live events
	for _, event := range backlogData {
		e.handle(event, true)
	}
}

func (e *eventHandler) handle(curr eventData, backlogEvent bool) {
	switch curr.Type {
	case "oob_include":
		oob_data := &oob_include{}
		json.Unmarshal(curr.Data, &oob_data)
		e.handleBacklog(oob_data.Url)

	case "channel_init":
		if !backlogEvent {
			user_strings := []string{}
			for _, user_string := range curr.Members {
				user_strings = append(user_strings, user_string.Nick)
			}

			e.Window.AddChannel(curr.Chan, curr.Topic.Text, curr.Cid, user_strings)
		}

	case "you_parted_channel":
		if !backlogEvent {
			e.Window.RemoveChannel(curr.Chan)
		}

	case "buffer_msg":
		e.Window.AddBufferMsg(curr.Chan, curr.From, curr.Msg)

	case "joined_channel":
		if !backlogEvent {
			e.Window.AddUser(curr.Chan, curr.Nick)
		}
		e.Window.AddJoinEvent(curr.Chan, curr.Nick, curr.Hostmask)

	case "parted_channel":
		if !backlogEvent {
			e.Window.RemoveUser(curr.Chan, curr.Nick)
		}
		e.Window.AddPartEvent(curr.Chan, curr.Nick, curr.Hostmask)

	case "nickchange":
		e.Window.ChangeUserNick(curr.Chan, curr.OldNick, curr.NewNick)

	case "quit":
		if !backlogEvent {
			e.Window.RemoveUser(curr.Chan, curr.Nick)
		}
		e.Window.AddQuitEvent(curr.Chan, curr.Nick, curr.Hostmask, curr.Msg)
	default:
		//fmt.Printf("Event: %s\n", curr.Type)
		return
	}
}
