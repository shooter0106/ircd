package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	users := make(Users)
	require.Len(t, users, 0)

	users.newUser(nil, "test")
	require.Len(t, users, 1)
}

func TestUsersgetCount(t *testing.T) {
	t.Parallel()

	users := make(Users)
	require.IsType(t, 0, users.getCount())
	require.Equal(t, 0, users.getCount())

	users.newUser(nil, "test")
	require.Equal(t, 1, users.getCount())
}
