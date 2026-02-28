package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Teaching struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Title        string     `gorm:"size:180;not null" json:"title"`
	Organization string     `gorm:"size:180;not null" json:"organization"`
	Location     string     `gorm:"size:160" json:"location"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	Description  string     `gorm:"type:text" json:"description,omitempty"`
	LinkURL      string     `gorm:"size:255" json:"link_url,omitempty"`
	SortOrder    int        `gorm:"not null;default:0" json:"sort_order"`
	Visible      bool       `gorm:"not null;default:true" json:"visible"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Teaching) TableName() string {
	return "teaching"
}

func (t *Teaching) BeforeCreate(_ *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
