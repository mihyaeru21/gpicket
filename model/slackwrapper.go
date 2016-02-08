package model

import (
	"fmt"
	"github.com/nlopes/slack"
	"os"
)

type Message struct {
	Timestamp string
	Team      string
	TeamID    string
	Channel   string
	ChannelID string
	User      string
	UserID    string
	Text      string
}

type SlackWrapper struct {
	token  string
	api    *slack.Client
	team   string
	teamID string
}

func NewSlack(token string) *SlackWrapper {
	return &SlackWrapper{
		token: token,
		api:   slack.New(token),
	}
}

func (self *SlackWrapper) Start(messages chan Message) {
	auth, err := self.api.AuthTest()
	if err != nil {
		fmt.Printf("Authentication failed. token: %s\n", self.token)
		os.Exit(1)
	}
	self.team = auth.Team
	self.teamID = auth.TeamID

	rtm := self.api.NewRTM()
	go rtm.ManageConnection()
	for {
		msg := <-rtm.IncomingEvents
		switch event := msg.Data.(type) {
		case *slack.HelloEvent:
			fmt.Printf("Logging for %s was started.\n", auth.Team)
		case *slack.MessageEvent:
			messages <- self.createMessage(event)
		case *slack.RTMError:
			fmt.Printf("Error: %s\n", event.Error())
		default:
		}
	}
}

func (self *SlackWrapper) createMessage(event *slack.MessageEvent) Message {
	return Message{
		Timestamp: event.Timestamp,
		Team:      self.team,
		TeamID:    self.teamID,
		ChannelID: event.Channel,
		UserID:    event.User,
		Text:      event.Text,
	}
}
