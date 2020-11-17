package events

import (
	"encoding/json"
	"fmt"
)

func (e *eventHandler) backlog(event eventData) {
	oobData := &oobInclude{}
	err := json.Unmarshal(event.Data, &oobData)

	if err == nil {
		e.handleBacklog(oobData.Url)
	}
}

func (e *eventHandler) channels(event eventData, backlogEvent bool) {
	if !backlogEvent {
		userStrings := []string{}
		for _, userString := range event.Members {
			userStrings = append(userStrings, userString.Nick)
		}
		topic := getTopicName(event.Topic)
		e.Window.AddChannel(event.Chan, topic, event.Cid, event.Bid, userStrings)
	}
}

func (e *eventHandler) parted(event eventData, backlogEvent bool) {
	if !backlogEvent {
		e.Window.RemoveChannel(event.Chan)
	}
}

func (e *eventHandler) msg(event eventData) {
	if e.Window.HasChannel(event.Chan) {
		e.Window.Activity.RegisterActivity(event.Chan, event.Msg, e.Window)
		e.Window.AddBufferMsg(event.Chan, event.From, event.Msg, event.Time, event.Bid)
	}
}

func (e *eventHandler) join(event eventData, backlogEvent bool) {
	if !backlogEvent {
		e.Window.AddUser(event.Chan, event.Nick, event.Bid)
	}
	e.Window.AddJoinEvent(event.Chan, event.Nick, event.Hostmask, event.Time, event.Bid)
}

func (e *eventHandler) part(event eventData, backlogEvent bool) {
	if !backlogEvent {
		e.Window.RemoveUser(event.Chan, event.Nick, event.Bid)
	}
	e.Window.AddPartEvent(event.Chan, event.Nick, event.Hostmask, event.Time, event.Bid)
}

func (e *eventHandler) nick(event eventData) {
	e.Window.ChangeUserNick(event.Chan, event.OldNick, event.NewNick, event.Time, event.Bid)
}

func (e *eventHandler) topic(event eventData) {
	e.Window.ChangeTopic(event.Chan, event.Author, getTopicText(event.Topic), event.Time, event.Bid)
}

func (e *eventHandler) conversation(event eventData) {
	if event.BufferType == "conversation" {
		header := fmt.Sprintf("Chatting since: %s", unixtimeToDate(event.Created))
		e.Window.AddChannel(event.Name, header, event.Cid, event.Bid, []string{})
	}
}

func (e *eventHandler) meMsg(event eventData) {
	e.Window.AddBufferMsg(event.Chan, event.From, event.Msg, event.Time, event.Bid)
}

func (e *eventHandler) userQuit(event eventData, backlogEvent bool) {
	if !backlogEvent {
		e.Window.RemoveUser(event.Chan, event.Nick, event.Bid)
	}
	e.Window.AddQuitEvent(event.Chan, event.Nick, event.Hostmask, event.Msg, event.Time, event.Bid)
}