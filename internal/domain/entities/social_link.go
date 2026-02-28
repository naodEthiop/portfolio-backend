package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SocialLink struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Platform  string    `gorm:"size:80;not null" json:"platform"`
	URL       string    `gorm:"size:255;not null" json:"url"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`
	Visible   bool      `gorm:"not null;default:true" json:"visible"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *SocialLink) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
