package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"os"
	"strings"
)

func isKudosMessage(message string) bool {
	return strings.Contains(message, ":kudos:") && strings.Contains(message, "<@U")
}

func main() {
	api := slack.New("YOUR_TOKEN")
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()

	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		info := rtm.GetInfo()

		switch ev := msg.Data.(type) {

		case *slack.MessageEvent:
			if (info.User.ID != ev.Msg.User) && isKudosMessage(ev.Msg.Text) {
				fmt.Println("Kudos to you!")

				msgRef := slack.NewRefToMessage(ev.Msg.Channel, ev.Msg.Timestamp)
				if err := api.AddReaction("+1", msgRef); err != nil {
					fmt.Printf("Error adding reaction: %s\n", err)
					continue
				}

				// rtm.SendMessage(rtm.NewOutgoingMessage("Kudos", ev.Msg.Channel))
			}

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Println("Invalid credentials")
			return

		default:
			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}
