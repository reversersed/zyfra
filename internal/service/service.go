package service

import (
	"fmt"
	"time"

	model "github.com/reversersed/zyfra/internal/storage"
)

type storage interface {
	LoginUser(string, string) (string, error)
	GetSession(string) (*model.UserSession, error)
	DeleteSession(string) error
}
type service struct {
	storage storage
}

func New(storage storage) *service {
	return &service{storage: storage}
}
func (s *service) CreateSession(login, password string) (string, error) {
	key, err := s.storage.LoginUser(login, password)
	if err != nil {
		return "", err
	}
	return key, nil
}

func (s *service) CheckSession(key string) error {
	session, err := s.storage.GetSession(key)
	if err != nil {
		return err
	}

	if session.Expiration.Before(time.Now().UTC()) {
		return fmt.Errorf("! Session is expired %.0f seconds ago\n! Please log in again", time.Now().UTC().Sub(session.Expiration).Seconds())
	}
	return nil
}
func (s *service) Delete(key string) error {
	if err := s.storage.DeleteSession(key); err != nil {
		return err
	}
	return nil
}
