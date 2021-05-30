package models

import "time"

// Repo represents a single git repository
type Repo struct {
	Name        string     `json:"name"`
	Stars       int        `json:"stars"`
	Forks       int        `json:"forks"`
	Url         string     `json:"url"`
	Description *string    `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	License     string     `json:"license"`
	Languages   []language `json:"languages"`
}

// language represent a repo's language and color code
type language struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Score returns a sum of the stars and forks of this repo
func (r Repo) Score() int {
	return r.Stars + r.Forks
}
