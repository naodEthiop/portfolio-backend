package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Award struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string     `gorm:"size:180;not null" json:"title"`
	Issuer      string     `gorm:"size:180;not null" json:"issuer"`
	AwardDate   *time.Time `json:"award_date,omitempty"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	LinkURL     string     `gorm:"size:255" json:"link_url,omitempty"`
	SortOrder   int        `gorm:"not null;default:0" json:"sort_order"`
	Visible     bool       `gorm:"not null;default:true" json:"visible"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (a *Award) BeforeCreate(_ *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
