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

func NewConnection(token string) *Connection {
	address := url.URL{Scheme: "wss", Host: "api.irccloud.com", Path: "/"}
	log.Printf("Connecting to: %s\n", address.String())

	headers := http.Header{}
	headers.Add("User-Agent", "irccloud-cli")
	headers.Add("Origin", "https://api.irccloud.com")
	headers.Add("Cookie", fmt.Sprintf("session=%s", token))

	conn, _, err := websocket.DefaultDialer.Dial(address.String(), headers)

	if err != nil {
		panic(err)
	}

	log.Printf("Connected!\n")

	return &Connection{
		WSConn: conn,
	}
}

func (c *Connection) SendMessage(cid int, channel, message string) {
	msg := &sayMessage{
		Method: "say",
		Cid:    cid,
		To:     channel,
		Msg:    message,
	}

	data, _ := json.Marshal(msg)

	c.writeMessage([]byte(data))
}

func (c *Connection) writeMessage(message []byte) error {
	return c.WSConn.WriteMessage(websocket.TextMessage, message)
}

func (c *Connection) ReadMessage() ([]byte, error) {
	_, msg, err := c.WSConn.ReadMessage()

	return msg, err
}
