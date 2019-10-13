package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

var usersList = make(Users)
var channelList = make(Channels)

func main() {
	l, err := net.Listen("tcp", ":6667")
	if err != nil {
		panic(err)
	}

	defer l.Close()

	for {
		// wait for connection
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go connectionListener(conn)
	}
}

func connectionListener(c net.Conn) {
	for {
		buf := make([]byte, 512)

		_, err := c.Read(buf)
		if err != nil {
			// Ignore error
			// panic(err)
		}

		if len(buf) == 0 {
			continue
		}

		for _, line := range strings.Split(string(buf), "\r\n") {
			line = strings.Trim(line, "\x00")

			if line == "" {
				continue
			}

			log.Println("<- " + line)

			cmd, err := parseLine(line)
			if err != nil {
				// just ignore unknown commands
				//panic(err)
				continue
			}

			execCommand(cmd, c)
		}
	}
}

func execCommand(cmd interface{}, c net.Conn) {
	user := usersList[c]

	switch cmd := cmd.(type) {
	case nickCommand:
		usersList = usersList.newUser(c, cmd.nickname)

	case quitCommand:
		c.Close()

	case joinCommand:
		channel, ok := channelList[cmd.channel]
		if !ok {
			channelList[cmd.channel] = newChannel(cmd.channel)
			channel = channelList[cmd.channel]
		}

		channel.addUser(user)

	case partCommand:
		channel, ok := channelList[cmd.channel]
		if !ok {
			channelList[cmd.channel] = newChannel(cmd.channel)
			channel = channelList[cmd.channel]
		}

		channel.removeUser(user, cmd.message)

	case privmsgCommand:
		switch {
		// send message to channel
		case cmd.receiver[0] == '#':
			channel := channelList[cmd.receiver]
			channel.sendMessage(user, cmd.message)
		}

	case listCommand:
		channelList.sendList(user)

	case setTopicCommand:
		channel := channelList.get(cmd.channel)
		channel.setTopic(cmd.topic)

	}
}
