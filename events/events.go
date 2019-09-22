package events

import (
	"encoding/json"
	"github.com/termoose/irccloud/requests"
	"log"
)

type event struct {
	Bid int32
	Eid int32
	Type string
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
}

func NewHandler(token string) (*eventHandler) {
	handler := &eventHandler{
		Queue: make(chan event, 8),
		SessionToken: token,
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

		//log.Println("BACKLOG: %s\n", backlog)
		// FIXME: parse the `channel_init` messages from backlog
		// and use them to initialize the joined channels list
	case "buffer_msg":
		msg_data := &buffer_msg{}
		json.Unmarshal(curr_event.Data, &msg_data)
		log.Printf("<%s> %s", msg_data.From, msg_data.Msg)
	}
}
