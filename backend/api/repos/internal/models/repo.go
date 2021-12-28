package models

import (
	"time"
)

// Repo represents a single git repository
type Repo struct {
	ID          string     `json:"id" firestore:"id"`
	Name        string     `json:"name" firestore:"name"`
	Stars       int        `json:"stars" firestore:"stars"`
	Forks       int        `json:"forks" firestore:"forks"`
	Url         string     `json:"url" firestore:"url"`
	Description *string    `json:"description" firestore:"description"`
	CreatedOn   time.Time  `json:"created_on" firestore:"firestore"`
	UpdatedOn   time.Time  `json:"updated_on" firestore:"updated_on"`
	License     string     `json:"license" firestore:"license"`
	Languages   []language `json:"languages" firestore:"languages"`
	Score       int        `json:"score" firestore:"score"`
	UpdatedAt   time.Time  `json:"updated_at" firestore:"updated_at,serverTimestamp"`
	Source      string     `json:"source" firestore:"source"`
}

// Id returns the if this Repo
func (r Repo) Id() string {
	return r.ID
}

// language represent a repo's language and color code
type language struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}
