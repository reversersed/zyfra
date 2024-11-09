package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	t.Run("existing config", func(t *testing.T) {
		dir := t.TempDir() + "/cfg.json"
		file, err := os.OpenFile(dir, os.O_CREATE, os.FileMode(0777))
		assert.NoError(t, err)

		_, err = file.WriteString(fmt.Sprintf("{\"user\": \"%s\"}", []byte("password")))
		assert.NoError(t, err)

		_, err = ReadFromFile(dir)
		assert.NoError(t, err)

		err = file.Close()
		assert.NoError(t, err)
	})
	t.Run("non-existing config", func(t *testing.T) {
		cfg, err := ReadFromFile("non-file.dot")
		assert.Error(t, err)
		assert.Empty(t, cfg)
	})
	t.Run("invalid file content", func(t *testing.T) {
		dir := t.TempDir() + "/cfg.json"
		file, err := os.OpenFile(dir, os.O_CREATE, os.FileMode(0777))
		assert.NoError(t, err)

		file.WriteString("oh my god is it a config file![][][][]")

		file.Close()

		cfg, err := ReadFromFile(dir)
		assert.Error(t, err)
		assert.Empty(t, cfg)
	})
}
