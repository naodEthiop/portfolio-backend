package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	FullName     string    `gorm:"size:150;not null" json:"full_name"`
	Handle       string    `gorm:"size:80" json:"handle"`
	Headline     string    `gorm:"size:200" json:"headline"`
	Location     string    `gorm:"size:160" json:"location"`
	Summary      string    `gorm:"size:500" json:"summary"`
	Bio          string    `gorm:"type:text" json:"bio"`
	AvatarURL    string    `gorm:"size:255" json:"avatar_url"`
	BannerURL    string    `gorm:"size:255" json:"banner_url"`
	ResumeURL    string    `gorm:"size:255" json:"resume_url"`
	CTAPrimary   string    `gorm:"size:120" json:"cta_primary"`
	CTASecondary string    `gorm:"size:120" json:"cta_secondary"`
	CTATertiary  string    `gorm:"size:120" json:"cta_tertiary"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Profile) TableName() string {
	return "profile"
}

func (p *Profile) BeforeCreate(_ *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
