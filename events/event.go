package events

import (
	"encoding/json"
	"sort"
	"github.com/termoose/irccloud/requests"
	"github.com/termoose/irccloud/ui"
	"log"
	"net/http"
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

// "topic":{"text":"http://www.pvv.org/~birkedal/term_0z_-_kose_til_krampa_tar_meg.mpg | HACK THE PLANET","time":1442082313,"nick":"ehamberg","ident_prefix":"","user":"sid18208","userhost":"gateway/web/irccloud.com/x-opdwzifkkmmqwndd","usermask":"sid18208@gateway/web/irccloud.com/x-opdwzifkkmmqwndd"}
type topic struct {
	Text string `json:"text"`
}

type oob_include struct {
	Url string
}

type eventData struct {
	Type     string
	Time     int64            `json:"eid"`
	Chan     string           `json:"chan"`
	Members  []member         `json:"members"`
	From     string           `json:"from"`
	Msg      string           `json:"msg"`
	Cid      int              `json:"cid"`
	Hostmask string           `json:"hostmask"`
	Nick     string           `json:"nick"`
	NewNick  string           `json:"newnick"`
	OldNick  string           `json:"oldnick"`
	Topic    json.RawMessage  `json:"topic"`
	Data     []byte
}

func getTopicName(e json.RawMessage) string {
	dst := &topic{}
	err := json.Unmarshal(e, dst)

	if err != nil {
		return dst.Text;
	}

	return "";
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

			topic := getTopicName(event.Topic)
			window.AddChannel(event.Chan, topic, event.Cid, user_strings)
		}
	}

	// Then we fill them with the message backlog
	for _, event := range backlogData {
		if event.Type == "buffer_msg" {
			window.AddBufferMsg(event.Chan, event.From, event.Msg, event.Time)
		}
	}
}

func parseBacklog(backlog *http.Response) []eventData {
	backlogData := []eventData{}
	decoder := json.NewDecoder(backlog.Body)
	err := decoder.Decode(&backlogData)

	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(backlogData, func(i, j int) bool {
		return backlogData[i].Time < backlogData[j].Time
	})

	return backlogData
}
