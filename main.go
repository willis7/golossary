package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/willis7/slack"
)

func init() {
	viper.SetConfigName("app_config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func main() {

	mux := slack.NewEventMux()
	mux.Handle("message", slack.HandlerFunc(RTMMessage))

	token := viper.GetString("slack.token")
	client := slack.NewClient(token, mux)
	client.Connect()
	defer client.Close()
	go client.Dispatch()

	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, os.Interrupt)
	for {
		select {
		case <-sigterm:
			log.Println("terminate signal recvd")
			err := client.Shutdown()
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-time.After(time.Second):
			}
			client.Close()
			return
		}
	}
}

// RTMMessage is a HandlerFunc implementation which handles the "message" event
func RTMMessage(msg *slack.Message, c *slack.Client) {
	parts := strings.Fields(msg.Text)
	if len(parts) == 3 && parts[1] == "define" {
		// TODO: concurrently call to the DB and postMessage
		c.PostMessage(&slack.Message{Type: msg.Type, Channel: msg.Channel, Text: fmt.Sprintf("%s means - ", parts[2])})
		// NOTE: the Message object is copied, this is intentional
	} else {
		// huh?
		msg.Text = fmt.Sprintf("sorry, that does not compute\n")
		c.PostMessage(msg)
	}
}
