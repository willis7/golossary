package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

// respRtmStart is the structure of the introductory response
// from slack
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

// Message
type Message struct {
	Id      uint64 `json:"id"`
	Type string    `json:"type"`
	Error struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	} `json:"error"`
	ReplyTo int    `json:"reply_to"`
	Channel string `json:"channel"`
	Ts      string `json:"ts"`
	User    string `json:"user"`
	Text    string `json:"text"`
}

// start does a rtm.start, and returns a websocket URL and user ID. The
// websocket URL can be used to initiate an RTM session.
func start(token string) (wsurl string, id string, err error) {
	api_url := viper.GetString("slack.api_url")
	url := fmt.Sprintf("%s/rtm.start?token=%s", api_url, token)
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


func Connect(token string) *websocket.Conn {
	wsurl, _, err := start(token)
	if err != nil {
		log.Fatal(err)
	}

	c, _, err := websocket.DefaultDialer.Dial(wsurl, nil)
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func Dispatcher(conn *websocket.Conn) {
	defer conn.Close()

	for {
		msg := &Message{}
		err := conn.ReadJSON(msg)
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %+v", msg)
	}
}
