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

	input = strings.ReplaceAll(input, "\n", "")
	inp := strings.Split(input[:len(input)-1], " ")

	if len(inp) == 0 || len(input) == 1 {
		err = errors.New("received empty command")
		return
	}
	cmd = inp[0]
	if len(inp) > 1 {
		args = inp[1:]
	}
	return
}
