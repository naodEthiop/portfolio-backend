package gormrepo

import (
	"context"

	"portfolio-backend/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AwardRepository struct {
	db *gorm.DB
}

func NewAwardRepository(db *gorm.DB) *AwardRepository {
	return &AwardRepository{db: db}
}

func (r *AwardRepository) Create(ctx context.Context, item *entities.Award) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *AwardRepository) ListVisible(ctx context.Context) ([]entities.Award, error) {
	var items []entities.Award
	err := r.db.WithContext(ctx).
		Where("visible = ?", true).
		Order("sort_order ASC, award_date DESC NULLS LAST, created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *AwardRepository) ListAdmin(ctx context.Context) ([]entities.Award, error) {
	var items []entities.Award
	err := r.db.WithContext(ctx).
		Order("visible DESC, sort_order ASC, award_date DESC NULLS LAST, created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *AwardRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Award, error) {
	var item entities.Award
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *AwardRepository) Update(ctx context.Context, item *entities.Award) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *AwardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Award{}, "id = ?", id).Error
}
