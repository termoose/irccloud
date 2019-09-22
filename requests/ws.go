package requests

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"github.com/gorilla/websocket"
)

type connection struct {
	WSConn *websocket.Conn
}

func NewConnection(token string) (*connection) {
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

	return &connection {
		WSConn: conn,
	}
}

func (c *connection) ReadMessage() ([]byte, error) {
	_, msg, err := c.WSConn.ReadMessage()

	return msg, err
}
