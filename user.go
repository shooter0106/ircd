package main

import (
	"log"
	"net"
)

type Users map[net.Conn]*User

func (users Users) newUser(conn net.Conn, nick string) Users {
	user := &User{
		connection: conn,
		nick:       nick,
		host:       "localhost",
	}
	users[conn] = user

	return users
}

func (users Users) getCount() int {
	return len(users)
}

type User struct {
	connection net.Conn
	nick       string
	host       string
}

func (u User) send(message string) {
	if u.connection == nil {
		return
	}

	log.Println("-> " + message)

	u.connection.Write([]byte(message))
}
