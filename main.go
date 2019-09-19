package main

import (
	"fmt"
	"irccloud/requests"

	"net/http"
	"net/url"
	"github.com/gorilla/websocket"
)

func main() {
	session := requests.GetSessionToken("birkedal85@gmail.com", "SECRET")
	fmt.Println("session:", session)

	address := url.URL{Scheme: "wss", Host: "api.irccloud.com", Path: "/"}

	fmt.Printf("Connecting to: %s\n", address.String())

	headers := http.Header{}
	headers.Add("User-Agent", "irccloud-cli")
	headers.Add("Origin", "https://api.irccloud.com")
	headers.Add("Cookie", fmt.Sprintf("session=%s", session))

	conn, _, err := websocket.DefaultDialer.Dial(address.String(), headers)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			fmt.Printf("Nothing...\n")
		}

		fmt.Printf("Received: %s\n", msg)
	}
}
