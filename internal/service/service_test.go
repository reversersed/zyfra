package service

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var s *service
var buf io.Writer

func TestMain(m *testing.M) {
	buf = bytes.NewBufferString("")
	log.SetOutput(buf)
	s = New()
	defer os.Remove("./data.json")

	m.Run()
}
func TestFullWay(t *testing.T) {
	id := s.CreateSession()

	err := s.CheckSession(id)
	assert.NoError(t, err)

	err = s.CheckSession("non-existing session")
	assert.Error(t, err)
}
func TestSessionExpiration(t *testing.T) {
	s.sessions["session_id"] = time.Date(1, time.April, 1, 0, 0, 0, 0, time.UTC)

	err := s.CheckSession("session_id")
	assert.Error(t, err)
}
func TestDeleteSession(t *testing.T) {
	id := s.CreateSession()

	err := s.CheckSession(id)
	assert.NoError(t, err)

	s.Delete(id)

	err = s.CheckSession(id)
	assert.Error(t, err)
}
