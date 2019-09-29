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

		for _, line := range strings.Split(string(buf), "\n\r") {
			strings.Trim(line, " ")

			if line == "" {
				continue
			}

			log.Println("<- " + line)

			//message := parseMessage(line)
			//execCommand(message, c)
		}
	}
}

/*func execCommand(message message, c net.Conn) {
	user := usersList[c]
	fmt.Println(message)

	switch message.command {
	case "NICK":
		usersList.newUser(c, message.params[1])

	case "PING":
		response := []byte("PONG")

		c.Write(response)

	case "QUIT":
		c.Close()

	case "JOIN":
		channel, ok := channelList[message.params[1]]
		if !ok {
			channelList[message.params[1]] = newChannel(message.params[1])

			channel = channelList[message.params[1]]
		}

		channel.addUser(&user)

	case "PRIVMSG":
		switch {
		case message.params[1][0] == '#':
			channel := channelList[message.params[1]]
			channel.sendMessage(&user, message.params[1])
		}

	case "LIST":
		fmt.Println(channelList)
		channelList.sendList(&user)

	case "TOPIC":
		channel := channelList.get(message.params[1])
		channel.setTopic(message.params[2])
	}
}
*/
