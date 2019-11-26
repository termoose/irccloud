package events

import (
	"encoding/json"
	"github.com/termoose/irccloud/requests"
	"github.com/termoose/irccloud/ui"
	"log"
	"net/http"
	"sort"
)

type member struct {
	Nick     string `json:"nick"`
	RealName string `json:"realname"`
	Server   string `json:"ircserver"`
	UserHost string `json:"userhost"`
	UserMask string `json:"usermask"`
	Mode     string `json:"mode"`
}

type topic struct {
	Text string `json:"text"`
}

type oob_include struct {
	Url string
}

type eventData struct {
	Type       string
	Time       int64           `json:"eid"`
	Chan       string          `json:"chan"`
	Members    []member        `json:"members"`
	From       string          `json:"from"`
	Msg        string          `json:"msg"`
	Cid        int             `json:"cid"`
	Bid        int             `json:"bid"`
	Hostmask   string          `json:"hostmask"`
	Nick       string          `json:"nick"`
	NewNick    string          `json:"newnick"`
	OldNick    string          `json:"oldnick"`
	Topic      json.RawMessage `json:"topic"`
	Author     string          `json:"author"`
	BufferType string          `json:"buffer_type"`
	Name       string          `json:"name"`
	Archived   bool            `json:"archived"`
	Created    int64           `json:"created"`
	Data       []byte
}

func getTopicText(e json.RawMessage) string {
	var dst string
	err := json.Unmarshal(e, &dst)

	if err != nil {
		log.Fatal(err)
	}

	return dst
}

func UserModeString(mode string) string {
	switch mode {
	case "o":
		return "@"
	case "h":
		return "%"
	case "v":
		return "+"
	default:
		return ""
	}
}

func getTopicName(e json.RawMessage) string {
	dst := &topic{}
	err := json.Unmarshal(e, dst)

	if err != nil {
		log.Fatal(err)
	}

	return dst.Text
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
			window.AddChannel(event.Chan, topic, event.Cid, event.Bid, user_strings)
		}
	}

	// Then we fill them with the message backlog
	for _, event := range backlogData {
		if event.Type == "buffer_msg" {
			window.AddBufferMsg(event.Chan, event.From, event.Msg, event.Time, event.Bid)
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
