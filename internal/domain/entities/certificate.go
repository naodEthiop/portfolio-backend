package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Certificate struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name          string     `gorm:"size:180;not null" json:"name"`
	Issuer        string     `gorm:"size:180;not null" json:"issuer"`
	IssueDate     *time.Time `json:"issue_date"`
	Description   string     `gorm:"type:text" json:"description"`
	CredentialURL string     `gorm:"size:255" json:"credential_url"`
	ImageURL      string     `gorm:"size:255" json:"image_url"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (c *Certificate) BeforeCreate(_ *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
