package requests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const formtokenUrl = "https://www.irccloud.com/chat/auth-formtoken"
const sessionUrl = "https://www.irccloud.com/chat/login"

type sessionReply struct {
	Success bool
	Session string
	Uid     uint32
}

type formtokenReply struct {
	Id      uint32
	Success bool
	Token   string
}

func GetBacklog(token, endpoint string) *http.Response {
	path := fmt.Sprintf("https://api.irccloud.com%s", endpoint)
	client := http.Client{}

	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Add("User-Agent", "irccloud-cli")
	req.Header.Add("Origin", "https://api.irccloud.com")
	req.Header.Add("Cookie", fmt.Sprintf("session=%s", token))

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error fetching %s\n", path)
	}

	return resp
}

func GetSessionToken(user, pass string) string {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	formToken := getFormtoken(httpClient)

	form := url.Values{}
	form.Add("token", formToken)
	form.Add("email", user)
	form.Add("password", pass)

	httpRequest, _ := http.NewRequest("POST", sessionUrl, strings.NewReader(form.Encode()))
	httpRequest.Header.Add("X-Auth-FormToken", formToken)
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(httpRequest)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	return parseSession(resp)
}

func parseSession(response *http.Response) string {
	decoder := json.NewDecoder(response.Body)
	rep := &sessionReply{}
	err := decoder.Decode(&rep)

	if err != nil {
		panic(err)
	}

	if !rep.Success {
		panic("Incorrect login!")
	}

	return rep.Session
}

func getFormtoken(client *http.Client) string {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	httpRequest, _ := http.NewRequest("POST", formtokenUrl, nil)
	httpRequest.Header.Add("Content-Length", "0")
	resp, _ := httpClient.Do(httpRequest)
	defer resp.Body.Close()

	return parseToken(resp)
}

func parseToken(response *http.Response) string {
	decoder := json.NewDecoder(response.Body)
	rep := &formtokenReply{}
	err := decoder.Decode(&rep)

	if err != nil {
		panic(err)
	}

	return rep.Token
}
