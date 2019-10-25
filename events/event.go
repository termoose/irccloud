package events

import (
	"encoding/json"
	"github.com/termoose/irccloud/requests"
	"github.com/termoose/irccloud/ui"
	"log"
)

// {"nick":"sytse","ident_prefix":"","user":"sytse","userhost":"swielinga.nl","usermask":"sytse@swielinga.nl","realname":"Sytse Wielinga","account":null,"ircserver":"leguin.freenode.net","mode":"","away":false,"avatar":null,"avatar_url":null}

// {"bid":43026393,"eid":1570625780315817,"type":"buffer_msg","from":"BB-Martino","chan":"#lnd","cid":84415,"statusmsg":false,"msg":"err=non-ascii data < what's that about?","hostmask":"~martino@bitbargain.co.uk","ident_prefix":"~","from_name":"martino","from_host":"bitbargain.co.uk","from_account":"BB-Martino","from_realname":"Martin","avatar":null,"avatar_url":null}
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

type eventData struct {
	Type     string
	Chan     string   `json:"chan"`
	Members  []member `json:"members"`
	From     string   `json:"from"`
	Msg      string   `json:"msg"`
	Cid      int      `json:"cid"`
	Hostmask string   `json:"hostmask"`
	Nick     string   `json:"nick"`
}

func InitBacklog(token, url string, window *ui.View) {
	backlog := requests.GetBacklog(token, url)
	backlogData := parseBacklog(backlog)

	// First we initialize all channels
	for _, event := range backlogData {
		if event.Type == "channel_init" {
			user_strings := []string{}
			for _, user_string := range event.Members {
				user_strings = append(user_strings, user_string.Nick)
			}

			window.AddChannel(event.Chan, event.Cid, user_strings)
			//log.Printf("event: %v\n", event.Chan)
		}
	}

	// Then we fill them with the message backlog
	for _, event := range backlogData {
		if event.Type == "buffer_msg" {
			window.AddBufferMsg(event.Chan, event.From, event.Msg)
		}
	}
}

func parseBacklog(backlog []byte) []eventData {
	backlog_data := []eventData{}
	err := json.Unmarshal(backlog, &backlog_data)

	if err != nil {
		return []eventData{}
	}

	return backlog_data
}

// func initBacklog([]data backlogData) {
// }

func parseChannelInit(event []byte) channelInit {
	chan_init := channelInit{}
	json.Unmarshal(event, &chan_init)

	log.Printf("Channel: %s\n", event)
	//fmt.Printf("Users: %v\n", chan_init.Members)

	return chan_init
}
