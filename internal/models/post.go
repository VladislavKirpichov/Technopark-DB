package models

import (
	"time"
)

//easyjson:json
type Post struct {
	ID       int       `json:"id"`
	Parent   int       `json:"parent"`
	Author   string    `json:"author"`
	Forum    string    `json:"forum"`
	Thread   int       `json:"thread"`
	Created  time.Time `json:"created"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Message  string    `json:"message"`
}

//easyjson:json
type Posts []*Post

//easyjson:json
type PostsQueryParams struct {
	Limit    int    `form:"limit"`
	Since    int    `form:"since"`
	SortType string `form:"sort"`
	Desc     bool   `form:"desc"`
}

//easyjson:json
type ParamsPost struct {
	Post   *Post   `json:"post,omitempty"`
	Author *User   `json:"author,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
}
