package models

import "time"

type Thread struct {
	ID      int       `json:"id"`
	Slug    string    `json:"slug"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Title   string    `json:"title"`
	Msg     string    `json:"message"`
	Created time.Time `json:"created"`
	Votes   int       `json:"votes"`
}
