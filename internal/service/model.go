package service

import "time"

const (
	dataFilePath = "./data.json"
)

type session struct {
	Id         []byte
	Expiration time.Time
}
