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
	token    string
	api      *slack.Client
	team     string
	teamID   string
	channels map[string]string
}

func NewSlack(token string) *SlackWrapper {
	return &SlackWrapper{
		token:    token,
		api:      slack.New(token),
		channels: map[string]string{},
	}
}

func (self *SlackWrapper) Start(messages chan Message) {
	if err := self.combineTeam(); err != nil {
		fmt.Printf("Authentication failed. token: %s\n", self.token)
		os.Exit(1)
	}
	if err := self.combineChannel(); err != nil {
		fmt.Printf("Authentication failed. token: %s\n", self.token)
		os.Exit(1)
	}

	rtm := self.api.NewRTM()
	go rtm.ManageConnection()
	for {
		msg := <-rtm.IncomingEvents
		switch event := msg.Data.(type) {
		case *slack.HelloEvent:
			fmt.Printf("Logging for %s was started.\n", self.team)
		case *slack.MessageEvent:
			messages <- self.createMessage(event)
		case *slack.RTMError:
			fmt.Printf("Error: %s\n", event.Error())
		default:
		}
	}
}

func (self *SlackWrapper) createMessage(event *slack.MessageEvent) Message {
	channelName, ok := self.channels[event.Channel]
	if !ok {
		channelName = event.Channel
	}

	return Message{
		Timestamp: event.Timestamp,
		Team:      self.team,
		TeamID:    self.teamID,
		Channel:   channelName,
		ChannelID: event.Channel,
		UserID:    event.User,
		Text:      event.Text,
	}
}

func (self *SlackWrapper) combineTeam() error {
	auth, err := self.api.AuthTest()
	if err != nil {
		return err
	}
	self.team = auth.Team
	self.teamID = auth.TeamID
	return nil
}

func (self *SlackWrapper) combineChannel() error {
	channels, err := self.api.GetChannels(false)
	if err != nil {
		return err
	}
	for i := 0; i < len(channels); i++ {
		channel := channels[i]
		self.channels[channel.ID] = channel.Name
	}

	groups, err := self.api.GetGroups(false)
	if err != nil {
		return err
	}
	for i := 0; i < len(groups); i++ {
		group := groups[i]
		self.channels[group.ID] = group.Name
	}

	ims, err := self.api.GetIMChannels()
	if err != nil {
		return err
	}
	for i := 0; i < len(ims); i++ {
		im := ims[i]
		self.channels[im.ID] = im.User
	}

	return nil
}
