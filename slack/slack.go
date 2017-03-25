package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"

	"golang.org/x/net/websocket"
)

// These two structures represent the response of the Slack API rtm.start.
// Only some fields are included. The rest are ignored by json.Unmarshal.
//type responseRtmStart struct {
//	Ok    bool         `json:"ok"`
//	Error string       `json:"error"`
//	Url   string       `json:"url"`
//	Self  responseSelf `json:"self"`
//}

type respRtmStart struct {
	Ok  bool `json:"ok"`
	URL string `json:"url"`
	Self struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Prefs struct {
		} `json:"prefs"`
		Created        int `json:"created"`
		ManualPresence string `json:"manual_presence"`
	} `json:"self"`
	Team struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		EmailDomain string `json:"email_domain"`
		Domain      string `json:"domain"`
		Icon struct {
		} `json:"icon"`
		MsgEditWindowMins int `json:"msg_edit_window_mins"`
		OverStorageLimit  bool `json:"over_storage_limit"`
		Prefs struct {
		} `json:"prefs"`
		Plan string `json:"plan"`
	} `json:"team"`
	Users    []interface{} `json:"users"`
	Channels []interface{} `json:"channels"`
	Groups   []interface{} `json:"groups"`
	Mpims    []interface{} `json:"mpims"`
	Ims      []interface{} `json:"ims"`
	Bots     []interface{} `json:"bots"`
	Error    string       `json:"error"`
}

// slackStart does a rtm.start, and returns a websocket URL and user ID. The
// websocket URL can be used to initiate an RTM session.
func slackStart(token string) (wsurl, id string, err error) {
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("API request failed with code %d", resp.StatusCode)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	var respObj respRtmStart
	err = json.Unmarshal(body, &respObj)
	if err != nil {
		return
	}

	if !respObj.Ok {
		err = fmt.Errorf("Slack error: %s", respObj.Error)
		return
	}

	wsurl = respObj.URL
	id = respObj.Self.ID
	return
}

// These are the messages read off and written into the websocket. Since this
// struct serves as both read and write, we include the "Id" field which is
// required only for writing.

type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func GetMessage(ws *websocket.Conn) (m Message, err error) {
	err = websocket.JSON.Receive(ws, &m)
	return
}

var counter uint64

func postMessage(ws *websocket.Conn, m Message) error {
	m.Id = atomic.AddUint64(&counter, 1)
	return websocket.JSON.Send(ws, m)
}

// Starts a websocket-based Real Time API session and return the websocket
// and the ID of the (bot-)user whom the token belongs to.
func SlackConnect(token string) (*websocket.Conn, string) {
	wsurl, id, err := slackStart(token)
	if err != nil {
		log.Fatal(err)
	}

	ws, err := websocket.Dial(wsurl, "", "https://api.slack.com/")
	if err != nil {
		log.Fatal(err)
	}

	return ws, id
}
