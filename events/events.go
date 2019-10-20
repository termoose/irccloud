package events

import (
	"encoding/json"
	_ "fmt"
	"github.com/termoose/irccloud/requests"
	"github.com/termoose/irccloud/ui"
	_ "log"
)

type event struct {
	Bid int32   `json:"bid"`
	Eid int32   `json:"eid"`
	Type string `json:"type"`
	Data []byte
}

type oob_include struct {
	Url string
}

type buffer_msg struct {
	From string
	Chan string
	Msg string
}

type eventHandler struct {
	Queue chan event
	SessionToken string
	Window *ui.View
}

func NewHandler(token string, w *ui.View) (*eventHandler) {
	handler := &eventHandler{
		Queue: make(chan event, 8),
		SessionToken: token,
		Window: w,
	}

	// Start consumer thread
	go func() {
		for curr_event := range handler.Queue {
			handler.handle(curr_event)
		}
	}()

	return handler
}

func (e *eventHandler) Enqueue(msg []byte) {
	current := event{}
	json.Unmarshal(msg, &current)

	// Attach raw message data
	current.Data = msg

	e.Queue <- current
}

func (e *eventHandler) handle(curr_event event) {
	switch curr_event.Type {
	case "oob_include":
		oob_data := &oob_include{}
		json.Unmarshal(curr_event.Data, &oob_data)
		backlog := requests.GetBacklog(e.SessionToken, oob_data.Url)

		//log.Printf("BACKLOG: %s\n", backlog)

		backlogData := parseBacklog(backlog)

		// First we initialize all channels
		for _, event := range backlogData {
			if event.Type == "channel_init" {
				user_strings := []string{}
				for _, user_string := range event.Members {
					user_strings = append(user_strings, user_string.Nick)
				}
				
				e.Window.AddChannel(event.Chan, event.Cid, user_strings)
				//log.Printf("event: %v\n", event.Chan)
			}
		}

		for _, event := range backlogData {
			if event.Type == "buffer_msg" {
				e.Window.AddBufferMsg(event.Chan, event.From, event.Msg)
			}
		}

	case "buffer_msg":
		msg_data := &eventData{}
		json.Unmarshal(curr_event.Data, &msg_data)
		e.Window.AddBufferMsg(msg_data.Chan, msg_data.From, msg_data.Msg)
	}
}
