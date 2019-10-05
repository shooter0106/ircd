package main

import (
	"fmt"
	"strings"
)

type Channel struct {
	name  string
	users []*User
	topic string
}

type Channels map[string]*Channel

func newChannel(name string) *Channel {
	return &Channel{
		name: name,
	}
}

func (c Channels) get(name string) *Channel {
	channel, ok := c[name]
	if !ok {
		return newChannel(name)
	}

	return channel
}

func (c *Channels) sendList(user *User) {
	message := fmt.Sprintf(":localhost.localdomain 321 %s Channel :Users  Name\n\r", user.nick)
	user.send(message)

	for _, ch := range *c {
		message = fmt.Sprintf(":localhost.localdomain 322 %s %s %d :%s\n\r", user.nick, ch.name, ch.getUsersCount(), ch.topic)
		user.send(message)
	}

	message = fmt.Sprintf(":localhost.localdomain 323 %s :End of /LIST\n\r", user.nick)
	user.send(message)
}

func (c *Channel) addUser(user *User) {
	c.users = append(c.users, user)
	message := ":" + user.nick + "!~" + user.nick + "@127.0.0.1 JOIN " + c.name + "\n\r"
	user.send(message)

	c.sendTopic(user)
	for _, u := range c.users {
		c.sendUsersList(u)
	}
}

// Remove user from channel
func (c *Channel) removeUser(user *User) {
	s := make([]*User, 0, len(c.users)-1)
	for _, v := range c.users {
		if v.connection != user.connection {
			s = append(s, v)
		}
	}

	c.users = s
	for _, u := range c.users {
		c.sendUsersList(u)
	}
}

func (c *Channel) sendMessage(from *User, message string) {
	for _, user := range c.users {
		if from.nick == user.nick {
			continue
		}

		out := fmt.Sprintf(":%s!~%s PRIVMSG %s %s\n\r", from.nick, from.host, c.name, message)

		user.send(out)
	}
}

// Send channel's topic to user
func (c Channel) sendTopic(user *User) {
	var message string

	if c.topic == "" {
		message = fmt.Sprintf("331 %s :No topic is set\n\r", c.name)
	} else {
		message = fmt.Sprintf(":localhost.localdomain 332 %s %s :%s\n\r", user.nick, c.name, c.topic)
	}

	user.send(message)
}

// Set channel's topic
func (c *Channel) setTopic(topic string) {
	c.topic = topic
	for _, user := range c.users {
		c.sendTopic(user)
	}
}

func (c *Channel) sendUsersList(user *User) {
	message := ":localhost.localdomain 353 " + user.nick + " = " + c.name + " :" + c.getUsersNicknamesString() + "\n\r"
	message += ":localhost.localdomain 366 " + user.nick + " " + c.name + " :End of /NAMES list.\n\r"

	user.send(message)
}

func (c *Channel) getUsersNicknames() []string {
	var buf []string

	for _, user := range c.users {
		buf = append(buf, user.nick)
	}

	return buf
}

func (c *Channel) getUsersNicknamesString() string {
	return strings.Join(c.getUsersNicknames(), " ")
}

func (c *Channel) getUsersCount() int {
	return len(c.users)
}
