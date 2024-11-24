package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	t.Run("wrong database config", func(t *testing.T) {
		dir := t.TempDir() + "/w.env"
		file, err := os.OpenFile(dir, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
		assert.NoError(t, err)

		_, err = file.WriteString("SERVER_HOST=localhost\nSERVER_PORT=80")
		assert.NoError(t, err)

		_, err = GetConfig(dir)
		assert.Error(t, err)

		err = file.Close()
		assert.NoError(t, err)
	})
	t.Run("existing config", func(t *testing.T) {
		dir := t.TempDir() + "/.env"
		file, err := os.OpenFile(dir, os.O_CREATE|os.O_RDWR, os.ModePerm)
		assert.NoError(t, err)

		_, err = file.WriteString("SERVER_HOST=localhost\nSERVER_PORT=80\nDB_AUTHDB=base\nDB_BASE=database\nDB_HOST=local\nDB_NAME=root\nDB_PASS=pass\nDB_PORT=10")
		assert.NoError(t, err)

		_, err = GetConfig(dir)
		assert.NoError(t, err)

		err = file.Close()
		assert.NoError(t, err)
	})
	t.Run("non-existing config", func(t *testing.T) {
		cfg, err := GetConfig("non-existing-file.env")
		assert.Error(t, err)
		assert.Empty(t, cfg)
	})
	t.Run("invalid file content", func(t *testing.T) {
		dir := t.TempDir() + "/cfg.json"
		file, err := os.OpenFile(dir, os.O_CREATE|os.O_RDWR, os.ModePerm)
		assert.NoError(t, err)

		file.WriteString("oh my god is it a config file![][][][]")

		file.Close()

		cfg, err := GetConfig(dir)
		assert.Error(t, err)
		assert.Empty(t, cfg)
	})
}
