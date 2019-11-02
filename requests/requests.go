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

const formtoken_url = "https://www.irccloud.com/chat/auth-formtoken"
const session_url = "https://www.irccloud.com/chat/login"

type session_reply struct {
	Success bool
	Session string
	Uid     uint32
}

type formtoken_reply struct {
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

	// response, _ := ioutil.ReadAll(resp.Body)

	// return response
}

func GetSessionToken(user, pass string) string {
	http_client := &http.Client{
		Timeout: 10 * time.Second,
	}

	formToken := getFormtoken(http_client)

	form := url.Values{}
	form.Add("token", formToken)
	form.Add("email", user)
	form.Add("password", pass)

	http_request, _ := http.NewRequest("POST", session_url, strings.NewReader(form.Encode()))
	http_request.Header.Add("X-Auth-FormToken", formToken)
	http_request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http_client.Do(http_request)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	return parseSession(resp)
}

func parseSession(response *http.Response) string {
	decoder := json.NewDecoder(response.Body)
	rep := &session_reply{}
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
	http_client := &http.Client{
		Timeout: 10 * time.Second,
	}

	http_request, _ := http.NewRequest("POST", formtoken_url, nil)
	http_request.Header.Add("Content-Length", "0")
	resp, _ := http_client.Do(http_request)
	defer resp.Body.Close()

	return parseToken(resp)
}

func parseToken(response *http.Response) string {
	decoder := json.NewDecoder(response.Body)
	rep := &formtoken_reply{}
	err := decoder.Decode(&rep)

	if err != nil {
		panic(err)
	}

	return rep.Token
}
