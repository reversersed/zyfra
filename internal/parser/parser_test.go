package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var t_parser *parser

func TestMain(m *testing.M) {
	t_parser = New()
	os.Exit(m.Run())
}

func TestParseCommand(t *testing.T) {
	t.Run("empty input", func(t *testing.T) {
		cmd, args, err := t_parser.ParseCommand("")

		assert.Error(t, err)
		assert.Empty(t, cmd)
		assert.Empty(t, args)
	})
	t.Run("wrong argument", func(t *testing.T) {
		cmd, args, err := t_parser.ParseCommand("\r")

		assert.Error(t, err)
		assert.Empty(t, cmd)
		assert.Empty(t, args)
	})

	t.Run("non-args command parsed", func(t *testing.T) {
		cmd, args, err := t_parser.ParseCommand("command")

		assert.NoError(t, err)
		assert.Equal(t, "command", cmd)
		assert.Empty(t, args)
	})
	t.Run("single-arg command parsed", func(t *testing.T) {
		cmd, args, err := t_parser.ParseCommand("command arg")

		assert.NoError(t, err)
		assert.Equal(t, "command", cmd)
		assert.Equal(t, []string{"arg"}, args)
	})
	t.Run("multi-arg command parsed", func(t *testing.T) {
		cmd, args, err := t_parser.ParseCommand("command arg1 args2")

		assert.NoError(t, err)
		assert.Equal(t, "command", cmd)
		assert.Equal(t, []string{"arg1", "args2"}, args)
	})
}
