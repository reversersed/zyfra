package models

type LoginCommand struct {
	Login    string `json:"login" form:"login"`
	Password string `json:"password" form:"password"`
}
