package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

const (
	ProjectStatusComplete   = "complete"
	ProjectStatusInProgress = "in_progress"
	ProjectStatusArchived   = "archived"
)

type Project struct {
	ID               uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title            string         `gorm:"size:140;not null" json:"title"`
	Slug             string         `gorm:"size:160;uniqueIndex;not null" json:"slug"`
	ShortDescription string         `gorm:"size:255;not null" json:"short_description"`
	Description      string         `gorm:"type:text;not null" json:"description"`
	Status           string         `gorm:"size:32;not null;default:in_progress" json:"status"`
	TechStack        pq.StringArray `gorm:"type:text[];not null;default:'{}'" json:"tech_stack"`
	Achievements     pq.StringArray `gorm:"type:text[];not null;default:'{}'" json:"achievements"`
	DemoURL          string         `gorm:"size:255" json:"demo_url"`
	RepoURL          string         `gorm:"size:255" json:"repo_url"`
	ImageURL         string         `gorm:"size:255" json:"image_url"`
	Featured         bool           `gorm:"not null;default:false" json:"featured"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (p *Project) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
