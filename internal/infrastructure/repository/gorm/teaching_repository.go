package gormrepo

import (
	"context"

	"portfolio-backend/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeachingRepository struct {
	db *gorm.DB
}

func NewTeachingRepository(db *gorm.DB) *TeachingRepository {
	return &TeachingRepository{db: db}
}

func (r *TeachingRepository) Create(ctx context.Context, item *entities.Teaching) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *TeachingRepository) ListVisible(ctx context.Context) ([]entities.Teaching, error) {
	var items []entities.Teaching
	err := r.db.WithContext(ctx).
		Where("visible = ?", true).
		Order("sort_order ASC, start_date DESC NULLS LAST, created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *TeachingRepository) ListAdmin(ctx context.Context) ([]entities.Teaching, error) {
	var items []entities.Teaching
	err := r.db.WithContext(ctx).
		Order("visible DESC, sort_order ASC, start_date DESC NULLS LAST, created_at DESC").
		Find(&items).Error
	return items, err
}

func (r *TeachingRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Teaching, error) {
	var item entities.Teaching
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *TeachingRepository) Update(ctx context.Context, item *entities.Teaching) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *TeachingRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Teaching{}, "id = ?", id).Error
}
