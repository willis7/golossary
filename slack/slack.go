package slack

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const RTMMessage = "message"

// respRtmStart is the structure of the introductory response
// from Slack
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

// Message is the conversation data structure
type Message struct {
	Id   uint64 `json:"id"`
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

type Client struct {
	conn    *websocket.Conn
	apiUrl  string
	token   string
	counter uint64
}

// NewClient
func NewClient(token string) *Client {
	return &Client{
		apiUrl: "https://slack.com/api",
		token:  token,
	}
}

// start does a rtm.start, and returns a websocket URL and user ID.
func (c *Client) start() (wsurl string, id string, err error) {
	url := fmt.Sprintf("%s/rtm.start?token=%s", c.apiUrl, c.token)
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

// Connect opens a websocket connection with Slack.
func (c *Client) Connect() error {
	wsurl, _, err := c.start()
	if err != nil {
		log.Fatal(err)
		return errors.Errorf("Client Connect error: %s", err.Error)
	}

	c.conn, _, err = websocket.DefaultDialer.Dial(wsurl, nil)
	if err != nil {
		log.Fatal(err)
		return errors.Errorf("Client Websocket Connect error: %s", err.Error)
	}

	return nil
}

// Dispatch reads the events from Slack and sends them to the correct Performer
func (c *Client) Dispatch() {
	defer c.conn.Close()

	for {
		msg, err := c.getMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		switch msg.Type {
		case RTMMessage:
			log.Printf("text: %s", msg.Text)
			if msg.Text == "hello" {
				c.postMessage(&Message{Type: msg.Type, Channel: msg.Channel, Text: fmt.Sprintf("Hello to you, %s", msg.User)})
			}
			break
		default:
			log.Printf("recv: %+v", msg)
		}
	}
}

// Shutdown cleanly closes a connection. A client should send a close
// frame and wait for the server to close the connection.
func (c *Client) Shutdown() error {
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return errors.Errorf("Client Shutdown error: %s", err.Error)
	}
	return nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) getMessage() (*Message, error) {
	msg := &Message{}
	err := c.conn.ReadJSON(msg)
	if err != nil {
		return msg, errors.Errorf("Client getMessage error: %s", err.Error)
	}
	return msg, nil
}

func (c *Client) postMessage(m *Message) error {
	m.Id = atomic.AddUint64(&c.counter, 1)
	err := c.conn.WriteJSON(m)
	if err != nil {
		return errors.Errorf("Client postMessage error: %s", err.Error)
	}
	return nil
}
