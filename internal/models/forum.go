package models

import "time"

//easyjson:json
type ForumQueryParams struct {
	Limit int       `form:"limit"`
	Since time.Time `form:"since"`
	Desc  bool      `form:"desc"`
}

//easyjson:json
type Forum struct {
	ID      int    `json:"id"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	User    string `json:"user"`
	Posts   int    `json:"posts"`
	Threads int    `json:"threads"`
}

//easyjson:json
type ForumUserQueryParams struct {
	Limit int    `form:"limit"`
	Since string `form:"since"`
	Desc  bool   `form:"desc"`
}
