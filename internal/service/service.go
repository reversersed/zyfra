package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	sessions []session
}

func New() *service {
	service := new(service)

	file, _ := os.OpenFile(dataFilePath, os.O_CREATE, os.FileMode(0777))
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&service.sessions); err != nil {
		if err.Error() != "EOF" {
			log.Printf("\n\n! Error parsing data file: %v\n! File may be corrupted\n\n", err)
		}
		service.sessions = make([]session, 0)
	} else {
		log.Println("Session data loaded successfully")
	}

	return service
}
func (s *service) saveData() {
	file, _ := os.OpenFile(dataFilePath, os.O_CREATE, os.FileMode(0777))
	defer file.Close()

	if err := json.NewEncoder(file).Encode(&s.sessions); err != nil {
		log.Printf("! Error saving data: %v", err)
	}
}
func (s *service) Close() error {
	s.saveData()
	return nil
}

func (s *service) CreateSession() string {
	key := primitive.NewObjectID().Hex()

	hash, _ := bcrypt.GenerateFromPassword([]byte(key), bcrypt.MinCost)
	s.sessions = append(s.sessions, session{Id: hash, Expiration: time.Now().UTC().Add(time.Minute)})

	s.saveData()
	return string(key)
}

func (s *service) CheckSession(key string) error {
	// P.S. You can get direct access to session using s.sessions[key] in case there is no encryption for session keys

	for _, v := range s.sessions {
		if err := bcrypt.CompareHashAndPassword(v.Id, []byte(key)); err == nil {
			if v.Expiration.Before(time.Now().UTC()) {
				return fmt.Errorf("! Session is expired %.0f seconds ago\n! Please log in again", time.Now().UTC().Sub(v.Expiration).Seconds())
			}
			return nil
		}
	}
	return errors.New("! Session not found")
}
func (s *service) Delete(key string) error {
	for i, v := range s.sessions {
		if err := bcrypt.CompareHashAndPassword(v.Id, []byte(key)); err == nil {
			s.sessions[i] = s.sessions[len(s.sessions)-1]
			s.sessions = s.sessions[:len(s.sessions)-1]
			s.saveData()
			return nil
		}
	}
	return errors.New("! Session not found")
}
