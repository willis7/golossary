package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/viper"
	"github.com/willis7/golossary/slack"
)

func init() {
	viper.SetConfigName("app_config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func main() {

	token := viper.GetString("slack.token")

	client := slack.NewClient(token)
	client.Connect()
	defer client.Close()
	go client.Dispatch()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
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
