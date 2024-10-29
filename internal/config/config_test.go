package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	t.Run("existing config", func(t *testing.T) {
		dir := t.TempDir() + "/cfg.json"
		file, err := os.OpenFile(dir, os.O_CREATE, os.FileMode(0777))
		assert.NoError(t, err)

		file.WriteString("{\"user\": \"password\"}")

		file.Close()

		cfg, err := Read(dir)
		assert.NoError(t, err)

		assert.Equal(t, map[string]string{"user": "password"}, cfg)
	})
	t.Run("non-existing config", func(t *testing.T) {
		cfg, err := Read("non-file.dot")
		assert.Error(t, err)
		assert.Empty(t, cfg)
	})
	t.Run("invalid file content", func(t *testing.T) {
		dir := t.TempDir() + "/cfg.json"
		file, err := os.OpenFile(dir, os.O_CREATE, os.FileMode(0777))
		assert.NoError(t, err)

		file.WriteString("oh my god is it a config file![][][][]")

		file.Close()

		cfg, err := Read(dir)
		assert.Error(t, err)
		assert.Empty(t, cfg)
	})
}
