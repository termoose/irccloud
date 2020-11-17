package events

import (
	"encoding/json"
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
		for currEvent := range handler.Queue {
			handler.handle(currEvent, false)
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
			userStr := []string{}
			for _, userString := range event.Members {
				userStr = append(userStr, userString.Nick)
			}

			topic := getTopicName(event.Topic)
			e.Window.AddChannel(event.Chan, topic, event.Cid, event.Bid, userStr)
		}
	}

	// Then we fill them with the message backlog, should we send these events
	// to the event queue to have them arrive before live events?
	for _, event := range backlogData {
		e.handle(event, true)
	}

	// Go to the last visited channel if it exists
	e.Window.SetLatestChannel()
	e.Window.Redraw()
}

func (e *eventHandler) handle(curr eventData, backlogEvent bool) {
	switch curr.Type {
	case "oob_include":
		e.backlog(curr)

	case "channel_init":
		e.channels(curr, backlogEvent)

	case "you_parted_channel":
		e.parted(curr, backlogEvent)

	case "buffer_msg":
		e.msg(curr)

	case "joined_channel":
		e.join(curr, backlogEvent)

	case "parted_channel":
		e.part(curr, backlogEvent)

	case "nickchange":
		e.nick(curr)

	case "channel_topic":
		e.topic(curr)

	case "makebuffer":
		e.conversation(curr)

	case "buffer_me_msg":
		e.meMsg(curr)

	case "quit":
		e.userQuit(curr, backlogEvent)
	default:
		//fmt.Printf("Event: %s\n", curr.Type)
		return
	}

	// We only redraw per event if it's not a backlog event to speed
	// up app start time
	if !backlogEvent {
		e.Window.Redraw()
	}
}

func unixtimeToDate(t int64) string {
	tm := time.Unix(t/1000000, 0)
	return tm.Format("Mon Jan 2 15:04:05 UTC 2006")
}
