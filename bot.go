package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

var api *slack.Client

var Token string
var TLDR string
var Username string
var Emoji string
var Version = "dev"
var showVersion = flag.Bool("version", false, "show version and exit")

func init() {
	flag.StringVar(&Token, "token", "", "Slack API bot token (can also be set with SLACK_TLDR_TOKEN)")
	flag.StringVar(&TLDR, "command", "+tldr", "Command to type in Slack to trigger TLDR functions")
	flag.StringVar(&Username, "username", "TL;DR", "Username for the bot")
	flag.StringVar(&Emoji, "emoji", ":stopwatch:", "Emoji icon for the bot")
}

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Println("slack-tldr version " + Version + "\nCopyright 2016 Michael D. Ivey <ivey@gweezlebur.com>")
		return
	}

	if t := os.Getenv("SLACK_TLDR_TOKEN"); t != "" {
		Token = t
	}
	if Token == "" {
		fmt.Println("SLACK_TLDR_TOKEN was not set - use -token or set env var")
		return
	}

	logger := log.New(os.Stdout, "slack-tldr: ", log.Lshortfile|log.LstdFlags)
	api = slack.New(Token)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				if strings.HasPrefix(ev.Text, TLDR) {
					handleTLDR(ev)
					continue
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				// Ignore other events...
			}
		}
	}
}

func handleTLDR(ev *slack.MessageEvent) {
	params := slack.PostMessageParameters{
		AsUser:    false,
		Username:  Username,
		IconEmoji: Emoji,
	}

	message := strings.TrimSpace(ev.Text)

	// TODO: skip this if adding
	listPins, _, err := api.ListPins(ev.Channel)
	if err != nil {
		fmt.Printf("Error listing pins: %s\n", err)
		return
	}
	for i, j := 0, len(listPins)-1; i < j; i, j = i+1, j-1 {
		listPins[i], listPins[j] = listPins[j], listPins[i]
	}

	if message == TLDR {
		msg := ""
		for i, item := range listPins {
			msg += fmt.Sprintf("*%d*. %s", i+1, item.Message.Text)
			user, err := api.GetUserInfo(item.Message.User)
			if err != nil {
				fmt.Printf("%s\n", err)
			} else {
				msg += " _-- @" + user.Name + "_"
			}
			msg += "\n"
		}
		msg += "\n\n_You can also use pinned messages to add/remove TL;DR posts_"
		_, _, err = api.PostMessage(ev.Channel, msg, params)
		return
	}

	if p := strings.TrimPrefix(message, TLDR+" remove "); p != message {
		pos, err := strconv.Atoi(p)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		if pos > len(listPins) {
			return
		}
		delete := listPins[pos-1]
		msgRef := slack.NewRefToMessage(ev.Channel, delete.Message.Timestamp)
		err = api.RemovePin(ev.Channel, msgRef)
		if err != nil {
			fmt.Printf("Error remove pin: %s\n", err)
			return
		}
		_, _, err = api.PostMessage(ev.Channel, "Unpinned "+delete.Message.Text, params)
		return
	}

	if n := strings.TrimPrefix(message, TLDR+" "); n != message {
		user, err := api.GetUserInfo(ev.User)
		if err != nil {
			fmt.Printf("%s\n", err)
		} else {
			params.Username = user.Name
			n += " _-- @" + user.Name + "_"
		}

		channelID, timestamp, err := api.PostMessage(ev.Channel, n, params)
		if err != nil {
			fmt.Printf("Error posting message: %s\n", err)
			return
		}
		msgRef := slack.NewRefToMessage(channelID, timestamp)

		if err := api.AddPin(channelID, msgRef); err != nil {
			fmt.Printf("Error adding pin: %s\n", err)
			return
		}
		return
	}

}
