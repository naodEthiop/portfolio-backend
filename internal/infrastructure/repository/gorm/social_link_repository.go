package gormrepo

import (
	"context"

	"portfolio-backend/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SocialLinkRepository struct {
	db *gorm.DB
}

func NewSocialLinkRepository(db *gorm.DB) *SocialLinkRepository {
	return &SocialLinkRepository{db: db}
}

func (r *SocialLinkRepository) Create(ctx context.Context, link *entities.SocialLink) error {
	return r.db.WithContext(ctx).Create(link).Error
}

func (r *SocialLinkRepository) ListVisible(ctx context.Context) ([]entities.SocialLink, error) {
	var links []entities.SocialLink
	err := r.db.WithContext(ctx).
		Where("visible = ?", true).
		Order("sort_order ASC, created_at ASC").
		Find(&links).Error
	return links, err
}

func (r *SocialLinkRepository) ListAdmin(ctx context.Context) ([]entities.SocialLink, error) {
	var links []entities.SocialLink
	err := r.db.WithContext(ctx).Order("sort_order ASC, created_at ASC").Find(&links).Error
	return links, err
}

func (r *SocialLinkRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.SocialLink, error) {
	var link entities.SocialLink
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *SocialLinkRepository) Update(ctx context.Context, link *entities.SocialLink) error {
	return r.db.WithContext(ctx).Save(link).Error
}

func (r *SocialLinkRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.SocialLink{}, "id = ?", id).Error
}
