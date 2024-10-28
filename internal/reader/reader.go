package reader

import (
	"bufio"
	"os"
)

type reader struct{}

func New() *reader {
	return new(reader)
}

func (*reader) WaitKey() string {
	var std = bufio.NewReader(os.Stdin)
	key, err := std.ReadString('\n')
	if err != nil {
		return ""
	}
	return key
}
