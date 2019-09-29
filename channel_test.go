package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChannelCreate(t *testing.T) {
	t.Parallel()

	c := newChannel("#test")
	require.IsType(t, &Channel{}, c)
	require.Equal(t, "#test", c.name)
}

func TestGetUsersCount(t *testing.T) {
	t.Parallel()

	c := Channel{}
	require.IsType(t, 0, c.getUsersCount())
	require.Equal(t, 0, c.getUsersCount())

	c.addUser(&User{})
	require.Equal(t, 1, c.getUsersCount())
}

func TestGetUsersNicknames(t *testing.T) {
	t.Parallel()

	c := Channel{}
	require.IsType(t, make([]string, 0), c.getUsersNicknames())
	require.Len(t, c.getUsersNicknames(), 0)

	user := &User{
		nick: "test",
	}
	c.addUser(user)
	require.Len(t, c.getUsersNicknames(), 1)
	require.Contains(t, c.getUsersNicknames(), user.nick)
}

func TestGetUsersNicknamesString(t *testing.T) {
	t.Parallel()

	c := Channel{}
	require.IsType(t, "", c.getUsersNicknamesString())
	require.Empty(t, c.getUsersNicknamesString())

	user := &User{
		nick: "test",
	}
	c.addUser(user)
	require.Len(t, c.getUsersNicknamesString(), len(user.nick))
	require.Equal(t, c.getUsersNicknamesString(), user.nick)

	user2 := &User{
		nick: "shooter",
	}
	c.addUser(user2)
	require.Len(t, c.getUsersNicknamesString(), len(user.nick)+len(user2.nick)+1)
	require.Equal(t, c.getUsersNicknamesString(), fmt.Sprintf("%s %s", user.nick, user2.nick))
}
