package events

import (
	"encoding/json"
	"log"
)

// {"nick":"sytse","ident_prefix":"","user":"sytse","userhost":"swielinga.nl","usermask":"sytse@swielinga.nl","realname":"Sytse Wielinga","account":null,"ircserver":"leguin.freenode.net","mode":"","away":false,"avatar":null,"avatar_url":null}

type member struct {
	Nick     string `json:"nick"`
	RealName string `json:"realname"`
	Server   string `json:"ircserver"`
	UserHost string `json:"userhost"`
	UserMask string `json:"usermask"`
}

type channelInit struct {
	Name    string   `json:"chan"`
	Members []member `json:"members"`
}

type backlogData struct {
	Type    string
	Chan    string   `json:"chan"`
	Members []member `json:"members"`
	//Events []event
}

func parseBacklog(backlog []byte) []backlogData {
	backlog_data := []backlogData{}
	err := json.Unmarshal(backlog, &backlog_data)

	if err != nil {
		return []backlogData{}
	}

	return backlog_data
}

func parseChannelInit(event []byte) channelInit {
	chan_init := channelInit{}
	json.Unmarshal(event, &chan_init)

	log.Printf("Channel: %s\n", event)
	//fmt.Printf("Users: %v\n", chan_init.Members)

	return chan_init
}
