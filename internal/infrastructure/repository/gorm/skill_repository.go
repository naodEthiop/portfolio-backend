package gormrepo

import (
	"context"

	"portfolio-backend/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SkillRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) *SkillRepository {
	return &SkillRepository{db: db}
}

func (r *SkillRepository) Create(ctx context.Context, skill *entities.Skill) error {
	return r.db.WithContext(ctx).Create(skill).Error
}

func (r *SkillRepository) List(ctx context.Context) ([]entities.Skill, error) {
	var skills []entities.Skill
	err := r.db.WithContext(ctx).Order("category ASC, sort_order ASC, name ASC").Find(&skills).Error
	return skills, err
}

func (r *SkillRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Skill, error) {
	var skill entities.Skill
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&skill).Error
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *SkillRepository) Update(ctx context.Context, skill *entities.Skill) error {
	return r.db.WithContext(ctx).Save(skill).Error
}

func (r *SkillRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Skill{}, "id = ?", id).Error
}
