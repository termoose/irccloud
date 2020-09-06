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

func GetSessionToken(user, pass string) (string, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	formToken, err := getFormtoken(httpClient)

	if err != nil {
		return "", fmt.Errorf("error getting session token: %v", err)
	}

	form := url.Values{}
	form.Add("token", formToken)
	form.Add("email", user)
	form.Add("password", pass)

	httpRequest, _ := http.NewRequest("POST", sessionUrl, strings.NewReader(form.Encode()))
	httpRequest.Header.Add("X-Auth-FormToken", formToken)
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(httpRequest)

	if err != nil {
		log.Print(err)
		return "", err
	}

	defer resp.Body.Close()
	return parseSession(resp)
}

func parseSession(response *http.Response) (string, error) {
	decoder := json.NewDecoder(response.Body)
	rep := &sessionReply{}
	err := decoder.Decode(&rep)

	if err != nil {
		return "", fmt.Errorf("error parsing auth reply: %w", err)
	}

	if !rep.Success {
		return "", fmt.Errorf("invalid login: %w", err)
	}

	return rep.Session, nil
}

func getFormtoken(client *http.Client) (string, error) {
	httpRequest, _ := http.NewRequest("POST", formtokenUrl, nil)
	httpRequest.Header.Add("Content-Length", "0")
	resp, err := client.Do(httpRequest)

	if err != nil {
		return "", fmt.Errorf("error getting form token: %w", err)
	}

	defer resp.Body.Close()

	return parseToken(resp)
}

func parseToken(response *http.Response) (string, error) {
	decoder := json.NewDecoder(response.Body)
	rep := &formtokenReply{}
	err := decoder.Decode(&rep)

	if err != nil {
		return "", fmt.Errorf("can't parse token response: %w", err)
	}

	return rep.Token, nil
}
