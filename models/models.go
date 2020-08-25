package models

import "time"

// GitResponse is the response given for gitlab
type GitResponse struct {
	Commits []Commit `json:"commits"`
}

// Commit is commit info
type Commit struct {
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	ShortID   string    `json:"short_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Parameter is the data required to get commit info
type Parameter struct {
	Channel string `form:"channel" json:"channel" xml:"channel" binding:"required"`
	Repo    string `form:"repo" json:"repo" xml:"repo" binding:"required"`
	From    string `form:"from" json:"from" xml:"from" binding:"required"`
	To      string `form:"to" json:"to" xml:"to"`
}
