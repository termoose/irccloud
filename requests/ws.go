package requests

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

type Connection struct {
	WSConn *websocket.Conn
}

type sayMessage struct {
	Method string `json:"_method"`
	Cid    int    `json:"cid"`
	To     string `json:"to"`
	Msg    string `json:"msg"`
}

type heartbeatMessage struct {
	Method         string `json:"_method"`
	SelectedBuffer int    `json:"selectedBuffer"`
	SeenEids       string `json:"seenEids"`
}

func NewConnection(data sessionReply) *Connection {
	address := url.URL{Scheme: "wss", Host: data.WSHost, Path: data.WSPath}

	headers := http.Header{}
	headers.Add("User-Agent", "irccloud-cli")
	headers.Add("Origin", "https://api.irccloud.com")
	headers.Add("Cookie", fmt.Sprintf("session=%s", data.Session))

	conn, _, err := websocket.DefaultDialer.Dial(address.String(), headers)

	if err != nil {
		log.Fatal(err)
	}

	return &Connection{
		WSConn: conn,
	}
}

func (c *Connection) SendHeartbeat(selected, cid, bid, eid int) {
	msg := &heartbeatMessage{
		Method: "heartbeat",
		SelectedBuffer: selected,
		SeenEids: makeSeenEids(cid, bid, eid),
	}

	data, _ := json.Marshal(msg)
	_ = c.writeMessage(data)
}

func (c *Connection) SendMessage(cid int, channel, message string) {
	msg := &sayMessage{
		Method: "say",
		Cid:    cid,
		To:     channel,
		Msg:    message,
	}

	data, _ := json.Marshal(msg)
	_ = c.writeMessage(data)
}

func (c *Connection) writeMessage(message []byte) error {
	return c.WSConn.WriteMessage(websocket.TextMessage, message)
}

func (c *Connection) ReadMessage() ([]byte, error) {
	_, msg, err := c.WSConn.ReadMessage()

	return msg, err
}

// "{\\"3\":{\\"4\\":1343825583263721}}"
func makeSeenEids(cid, bid, eid int) string {
	return fmt.Sprintf(`{"%d":{"%d":%d}}`, cid, bid, eid)
}
