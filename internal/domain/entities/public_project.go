package entities

import "time"

type ProjectSource string

const (
	ProjectSourceGithub ProjectSource = "github"
	ProjectSourceManual ProjectSource = "manual"
)

// PublicProject is the unified project card shape consumed by the public frontend.
// It can represent either a GitHub repository or a manually-managed (admin) project.
type PublicProject struct {
	ID          string        `json:"id"`
	Source      ProjectSource `json:"source"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`

	TechStack []string `json:"tech_stack"`
	Topics    []string `json:"topics,omitempty"`
	Language  string   `json:"language,omitempty"`

	Stars int `json:"stars"`
	Forks int `json:"forks"`

	DemoURL string `json:"demo_url,omitempty"`
	RepoURL string `json:"repo_url,omitempty"`

	Featured  bool      `json:"featured"`
	Pinned    bool      `json:"pinned"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PublicProjectsResponse struct {
	Pinned   []PublicProject `json:"pinned"`
	Projects []PublicProject `json:"projects"`
}
