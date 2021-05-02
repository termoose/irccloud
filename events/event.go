package events

import (
	"encoding/json"
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

type oobInclude struct {
	Url string
}

// {"84415":{"4440297":1605131885388611}}}
type BidToEid map[string]int
// cid -> bid -> eid
type Seen map[string]BidToEid

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
	LastEid    int             `json:"last_seen_eid"`
	SeenEids   Seen            `json:"seenEids"`
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
