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
		buf := make([]byte, 1024)

		_, err := c.Read(buf)
		if err != nil {
			// Ignore error
			// panic(err)
		}

		if len(buf) > 0 {
			parseCommand(string(buf), c)
		}
	}
}

func parseCommand(input string, c net.Conn) {
	for _, cmd := range strings.Split(input, "\n") {
		if cmd == "" {
			return
		}

		log.Println("<-- " + cmd)

		user := usersList[c]
		command := strings.Fields(cmd)

		switch command[0] {
		case "NICK":
			usersList.newUser(c, command[1])

		case "PING":
			response := []byte("PONG")

			c.Write(response)

		case "QUIT":
			c.Close()

		case "JOIN":
			channel, ok := channelList[command[1]]
			if !ok {
				channelList[command[1]] = newChannel(command[1])

				channel = channelList[command[1]]
			}

			channel.addUser(&user)

		case "PRIVMSG":
			switch {
			case command[1][0] == '#':
				channel := channelList[command[1]]
				channel.sendMessage(&user, command[2])
			}

		case "LIST":
			fmt.Println(channelList)
			channelList.sendList(&user)

		case "TOPIC":
			channel := channelList.get(command[1])
			channel.setTopic(command[2])
		}
	}
}
