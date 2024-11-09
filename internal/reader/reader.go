package reader

import (
	"bufio"
	"io"
)

type reader struct{}

func New() *reader {
	return new(reader)
}

func (*reader) WaitForInput(read io.Reader) string {
	var std = bufio.NewReader(read)
	key, err := std.ReadString('\n')
	if err != nil {
		return ""
	}
	return key
}
