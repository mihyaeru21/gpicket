package model

import (
	"fmt"
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

func (self *Message) MakeStringForLog() string {
	return fmt.Sprintf("[%s][#%s][%s]%s", self.Team, self.Channel, self.User, self.Text)
}
