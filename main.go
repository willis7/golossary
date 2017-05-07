package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/willis7/golossary/models"
	"github.com/willis7/slack"
)

type App struct {
	DB *bolt.DB
}

func init() {
	viper.SetConfigName("app_config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func main() {

	db, err := models.InitDB(viper.GetString("database.path"))
	if err != nil {

	}
	defer db.Close()
	app := &App{DB: db}

	mux := slack.NewEventMux()
	mux.Handle("message", RTMMessage(app))

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

// RTMMessage is app HandlerFunc implementation which handles the "message" event
func RTMMessage(app *App) slack.Handler {
	return slack.HandlerFunc(func(msg *slack.Message, c*slack.Client) {
		parts := strings.Fields(msg.Text)

		switch parts[1] {
		case "define":
			result := models.Get(app.DB, parts[2])
			c.PostMessage(&slack.Message{Type: msg.Type, Channel: msg.Channel, Text: fmt.Sprintf("%s means - %s", parts[2], result)})
		case "insert":
			word := models.Word{Name: parts[2], Description: strings.Join(parts[3:], " ")}
			models.Update(app.DB, word)
			c.PostMessage(&slack.Message{Type: msg.Type, Channel: msg.Channel, Text: fmt.Sprintf("%s added ", parts[2])})
		default:
			msg.Text = fmt.Sprintf("sorry, that does not compute\n")
			c.PostMessage(msg)
		}
	})
}
