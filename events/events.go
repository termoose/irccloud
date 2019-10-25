package events

import (
	"encoding/json"
	_ "fmt"
	_ "github.com/termoose/irccloud/requests"
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
		InitBacklog(e.SessionToken, oob_data.Url, e.Window)

	case "buffer_msg":
		msg_data := &eventData{}
		json.Unmarshal(curr_event.Data, &msg_data)
		e.Window.AddBufferMsg(msg_data.Chan, msg_data.From, msg_data.Msg)

	case "joined_channel":
		join_data := &eventData{}
		json.Unmarshal(curr_event.Data, &join_data)
		e.Window.AddUser(join_data.Chan, join_data.Nick)
		e.Window.AddJoinEvent(join_data.Chan, join_data.Nick, join_data.Hostmask)

	case "parted_channel":
		part_data := &eventData{}
		json.Unmarshal(curr_event.Data, &part_data)
		e.Window.RemoveUser(part_data.Chan, part_data.Nick)
		e.Window.AddPartEvent(part_data.Chan, part_data.Nick, part_data.Hostmask)
	}
}
