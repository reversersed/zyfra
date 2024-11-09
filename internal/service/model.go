package service

import "time"

const (
	dataFilePath = "./data.json"
)

type session map[string]time.Time
