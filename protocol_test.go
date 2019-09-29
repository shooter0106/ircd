package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateCommand(t *testing.T) {
	t.Parallel()

	for cmd := range commandHandlers {
		err := validateCommand(cmd)
		require.NoError(t, err)
	}

	err := validateCommand("PINGG")
	require.Error(t, err)

	err = validateCommand("")
	require.Error(t, err)
}

func TestSplitLine(t *testing.T) {
	var tests = map[string][]string{
		"":                 []string(nil),
		"NICK Wiz":         []string{"NICK", "Wiz"},
		":WiZ NICK Kilroy": []string{"NICK", "Kilroy"},
		"USER guest tolmoon tolsun :Ronnie Reagan":                []string{"USER", "guest", "tolmoon", "tolsun", "Ronnie Reagan"},
		":testnick USER guest tolmoon tolsun :Ronnie Reagan":      []string{"USER", "guest", "tolmoon", "tolsun", "Ronnie Reagan"},
		":Trillian SQUIT cm22.eng.umd.edu :Server out of control": []string{"SQUIT", "cm22.eng.umd.edu", "Server out of control"},
	}

	t.Parallel()

	for in, out := range tests {
		args := splitLine(in)
		require.IsType(t, make([]string, 0), args)
		require.Equal(t, out, args)
	}
}

func TestParseLine(t *testing.T) {
	t.Parallel()

	cmd, err := parseLine("NICK Wiz")
	require.IsType(t, nickCommand{}, cmd)
	require.NoError(t, err)

	nickCmd := cmd.(nickCommand)
	require.Equal(t, "NICK", nickCmd.cmd)
	require.Equal(t, "Wiz", nickCmd.nickname)
}
