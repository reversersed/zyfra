package parser

import (
	"errors"
	"strings"
)

type parser struct{}

func New() *parser {
	return new(parser)
}

func (*parser) ParseCommand(input string) (cmd string, args []string, err error) {
	if len(input) == 0 {
		err = errors.New("received empty command")
		return
	}

	inp := strings.Fields(strings.ReplaceAll(strings.ReplaceAll(input, "\r", ""), "\n", ""))

	if len(inp) == 0 || len(inp[0]) == 0 {
		err = errors.New("received empty command")
		return
	}
	cmd = inp[0]
	if len(inp) > 1 {
		args = inp[1:]
	}
	return
}
