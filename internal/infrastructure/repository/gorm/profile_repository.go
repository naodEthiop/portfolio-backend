package gormrepo

import (
	"context"

	"portfolio-backend/internal/domain/entities"

	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) Get(ctx context.Context) (*entities.Profile, error) {
	var profile entities.Profile
	err := r.db.WithContext(ctx).Order("created_at ASC").First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfileRepository) Upsert(ctx context.Context, profile *entities.Profile) error {
	var existing entities.Profile
	err := r.db.WithContext(ctx).Order("created_at ASC").First(&existing).Error
	if err == nil {
		profile.ID = existing.ID
		return r.db.WithContext(ctx).Model(&existing).Updates(map[string]any{
			"full_name":     profile.FullName,
			"handle":        profile.Handle,
			"headline":      profile.Headline,
			"location":      profile.Location,
			"summary":       profile.Summary,
			"bio":           profile.Bio,
			"avatar_url":    profile.AvatarURL,
			"banner_url":    profile.BannerURL,
			"resume_url":    profile.ResumeURL,
			"cta_primary":   profile.CTAPrimary,
			"cta_secondary": profile.CTASecondary,
			"cta_tertiary":  profile.CTATertiary,
		}).Error
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.WithContext(ctx).Create(profile).Error
}
