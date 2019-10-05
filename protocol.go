package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type baseCommand struct {
	connection *net.Conn
	from       *User
	cmd        string
}

type parseCommandFunc func(msg []string) interface{}

var commandHandlers = map[string]parseCommandFunc{
	"NICK":    parseNickMessage,
	"QUIT":    parseQuitMessage,
	"JOIN":    parseJoinCommand,
	"PRIVMSG": parsePrivmsgCommand,
	"LIST":    parseListCommand,
	"TOPIC":   parseTopicCommand,
}

// Validate command name
func validateCommand(command string) error {
	if command == "" {
		return errors.New("Empty command")
	}

	_, ok := commandHandlers[command]
	if !ok {
		return fmt.Errorf("Unknown command: %s", command)
	}

	return nil
}

// Parse message: extract command name and call handler
func parseLine(line string) (interface{}, error) {
	msg := splitLine(line)
	if len(msg) == 0 {
		return nil, errors.New("Can not parse message line")
	}

	cmd := msg[0]
	err := validateCommand(cmd)
	if err != nil {
		return nil, err
	}

	handler := commandHandlers[cmd]

	return handler(msg), nil
}

// Splits message line to array
func splitLine(line string) (args []string) {
	if line == "" {
		return
	}

	fields := strings.Fields(line)
	for i, f := range fields {
		if strings.HasPrefix(f, ":") {
			// ignore message prefix
			if i == 0 {
				continue
			}

			trailing := strings.Join(fields[i:], " ")
			trailing = strings.TrimLeft(trailing, ":")
			args = append(args, trailing)

			return
		}

		args = append(args, f)
	}

	return
}

// NICK

type nickCommand struct {
	nickname string
	baseCommand
}

func parseNickMessage(msg []string) interface{} {
	cmd := nickCommand{}
	cmd.cmd = msg[0]
	cmd.nickname = msg[1]

	return cmd
}

// QUIT

type quitCommand struct {
	message string
	baseCommand
}

func parseQuitMessage(msg []string) interface{} {
	cmd := quitCommand{}
	cmd.cmd = msg[0]
	cmd.message = msg[1]

	return cmd
}

// JOIN

type joinCommand struct {
	channel string
	key     string
	baseCommand
}

func parseJoinCommand(msg []string) interface{} {
	cmd := joinCommand{}
	cmd.cmd = msg[0]
	cmd.channel = msg[1]
	//cmd.key = msg[2]

	return cmd
}

// PRIVMSG

type privmsgCommand struct {
	receiver string
	message  string
	baseCommand
}

func parsePrivmsgCommand(msg []string) interface{} {
	cmd := privmsgCommand{}
	cmd.cmd = msg[0]
	cmd.receiver = msg[1]
	cmd.message = msg[2]

	return cmd
}

// LIST

type listCommand struct {
	// channel and server args ignored
	baseCommand
}

func parseListCommand(msg []string) interface{} {
	cmd := listCommand{}
	cmd.cmd = msg[0]

	return cmd
}

// TOPIC

type getTopicCommand struct {
	channel string
	baseCommand
}

type setTopicCommand struct {
	channel string
	topic   string
	baseCommand
}

func parseTopicCommand(msg []string) interface{} {
	switch len(msg) {
	case 2:
		cmd := getTopicCommand{}
		cmd.cmd = msg[0]
		cmd.channel = msg[1]

		return cmd
	case 3:
		cmd := setTopicCommand{}
		cmd.cmd = msg[0]
		cmd.channel = msg[1]
		cmd.topic = msg[2]

		return cmd
	default:
		panic("Wrong TOPIC message")
	}
}
