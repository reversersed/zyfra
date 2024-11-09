package reader

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var t_reader *reader

func TestMain(m *testing.M) {
	t_reader = New()
	os.Exit(m.Run())
}

func TestRead(t *testing.T) {
	t.Run("read string", func(t *testing.T) {
		read := strings.NewReader("example string\n")
		str := t_reader.WaitForInput(read)
		assert.Equal(t, "example string\n", str)
	})
	t.Run("empty reader", func(t *testing.T) {
		read := strings.NewReader("")
		str := t_reader.WaitForInput(read)
		assert.Empty(t, str)
	})
}
