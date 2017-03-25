package main

import (
	"fmt"
	"log"
	"strings"
	"github.com/willis7/golossary/slack"
)

func main() {
	ws, id := slack.SlackConnect("xoxb-158974596516-mPovn5Bbqd6wAiejarChjn0q")

	fmt.Println("Golossary ready, ^C exits")

	for {
		// read each incoming message
		m, err := slack.GetMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
			fmt.Sprintf("From Golossary: %s", m.Text)
		}
	}
}
