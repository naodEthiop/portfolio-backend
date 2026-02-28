package entities

import "time"

// GithubProject is a lightweight, frontend-oriented view of a GitHub repository
// intended for the public portfolio projects endpoint.
type GithubProject struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	FullName    string   `json:"full_name"`
	HTMLURL     string   `json:"html_url"`
	Homepage    string   `json:"homepage,omitempty"`
	Description string   `json:"description,omitempty"`
	Topics      []string `json:"topics"`
	Language    string   `json:"language,omitempty"`

	Stars     int       `json:"stars"`
	Forks     int       `json:"forks"`
	Archived  bool      `json:"archived"`
	UpdatedAt time.Time `json:"updated_at"`
	PushedAt  time.Time `json:"pushed_at"`
}
