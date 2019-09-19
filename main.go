package main

import (
	"fmt"
	"irccloud/requests"
)

func main() {
	session := requests.GetSessionToken("birkedal85@gmail.com", "SECRET")
	fmt.Println("session:", session)
}
