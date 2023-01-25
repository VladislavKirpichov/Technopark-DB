package models

type Vote struct {
	Username string `json:"nickname"`
	Voice    int    `json:"voice"`
}
