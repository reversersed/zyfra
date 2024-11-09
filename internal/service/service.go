package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type service struct {
	sessions session
}

func New() *service {
	service := new(service)

	file, _ := os.OpenFile(dataFilePath, os.O_CREATE, os.FileMode(0777))
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&service.sessions); err != nil {
		if err.Error() != "EOF" {
			log.Printf("\n\n! Error parsing data file: %v\n! File may be corrupted\n\n", err)
		}
		service.sessions = make(session, 0)
	} else {
		log.Println("Session data loaded successfully")
	}

	return service
}
func (s *service) saveData() {
	file, _ := os.OpenFile(dataFilePath, os.O_CREATE|os.O_TRUNC, os.FileMode(0777))
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("! Error closing data file: %v", err)
		}
	}()

	if err := json.NewEncoder(file).Encode(&s.sessions); err != nil {
		log.Printf("! Error saving data: %v", err)
	}
}
func (s *service) CreateSession() string {
	key := primitive.NewObjectID().Hex()

	s.sessions[key] = time.Now().UTC().Add(time.Minute)

	s.saveData()
	return string(key)
}

func (s *service) CheckSession(key string) error {
	session, ok := s.sessions[key]
	if !ok {
		return errors.New("! Session not found")
	}

	if session.Before(time.Now().UTC()) {
		return fmt.Errorf("! Session is expired %.0f seconds ago\n! Please log in again", time.Now().UTC().Sub(session).Seconds())
	}
	return nil

}
func (s *service) Delete(key string) error {
	if _, exist := s.sessions[key]; !exist {
		return errors.New("session does not exist")
	}

	delete(s.sessions, key)
	s.saveData()
	return nil
}
