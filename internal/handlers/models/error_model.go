package models

type RequestError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
