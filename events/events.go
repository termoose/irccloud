package events

import (
	"encoding/json"
	"fmt"
	"github.com/termoose/irccloud/requests"
	"github.com/termoose/irccloud/ui"
	_ "log"
	"time"
)

type eventHandler struct {
	Queue        chan eventData
	SessionToken string
	Window       *ui.View
}

func NewHandler(token string, w *ui.View) *eventHandler {
	handler := &eventHandler{
		Queue:        make(chan eventData, 8),
		SessionToken: token,
		Window:       w,
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
	err := json.Unmarshal(msg, &current)

	if err == nil {
		// Attach raw message data
		current.Data = msg

		e.Queue <- current
	}
}

func (e *eventHandler) handleBacklog(url string) {
	backlogResponse := requests.GetBacklog(e.SessionToken, url)
	backlogData := parseBacklog(backlogResponse)

	// First we initialize all channels
	for _, event := range backlogData {
		if event.Type == "channel_init" {
			user_strings := []string{}
			for _, user_string := range event.Members {
				user_strings = append(user_strings, user_string.Nick)
			}

			topic := getTopicName(event.Topic)
			e.Window.AddChannel(event.Chan, topic, event.Cid, user_strings)
		}
	}

	// Then we fill them with the message backlog, should we send these events
	// to the event queue to have them arrive before live events
	for _, event := range backlogData {
		e.handle(event, true)
	}

	// Remove splash art when backlog has been parsed
	//e.Window.HideSplash()
}

func (e *eventHandler) handle(curr eventData, backlogEvent bool) {
	switch curr.Type {
	case "oob_include":
		oob_data := &oob_include{}
		err := json.Unmarshal(curr.Data, &oob_data)

		if err == nil {
			e.handleBacklog(oob_data.Url)
		}

	case "channel_init":
		if !backlogEvent {
			user_strings := []string{}
			for _, user_string := range curr.Members {
				user_strings = append(user_strings, user_string.Nick)
			}
			topic := getTopicName(curr.Topic)
			e.Window.AddChannel(curr.Chan, topic, curr.Cid, user_strings)
		}

	case "you_parted_channel":
		if !backlogEvent {
			e.Window.RemoveChannel(curr.Chan)
		}

	case "buffer_msg":
		e.Window.Activity.RegisterActivity(curr.Chan)
		e.Window.AddBufferMsg(curr.Chan, curr.From, curr.Msg, curr.Time)

	case "joined_channel":
		if !backlogEvent {
			e.Window.AddUser(curr.Chan, curr.Nick)
		}
		e.Window.AddJoinEvent(curr.Chan, curr.Nick, curr.Hostmask, curr.Time)

	case "parted_channel":
		if !backlogEvent {
			e.Window.RemoveUser(curr.Chan, curr.Nick)
		}
		e.Window.AddPartEvent(curr.Chan, curr.Nick, curr.Hostmask, curr.Time)

	case "nickchange":
		e.Window.ChangeUserNick(curr.Chan, curr.OldNick, curr.NewNick, curr.Time)

	case "channel_topic":
		e.Window.ChangeTopic(curr.Chan, curr.Author, getTopicText(curr.Topic), curr.Time)

	case "makebuffer":
		if curr.BufferType == "conversation" && !curr.Archived {
			header := fmt.Sprintf("Chatting since: %s", unixtimeToDate(curr.Created))
			e.Window.AddChannel(curr.Name, header, curr.Cid, []string{})
		}

	case "buffer_me_msg":
		e.Window.AddBufferMsg(curr.Chan, curr.From, curr.Msg, curr.Time)

	case "quit":
		if !backlogEvent {
			e.Window.RemoveUser(curr.Chan, curr.Nick)
		}
		e.Window.AddQuitEvent(curr.Chan, curr.Nick, curr.Hostmask, curr.Msg, curr.Time)
	default:
		//fmt.Printf("Event: %s\n", curr.Type)
		return
	}
}

func unixtimeToDate(t int64) string {
	tm := time.Unix(t/1000000, 0)
	return tm.Format("Mon Jan 2 15:04:05 UTC 2006")
}
