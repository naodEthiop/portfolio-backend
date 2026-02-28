package gormrepo

import (
	"context"

	"portfolio-backend/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *entities.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *ProjectRepository) ListPublic(ctx context.Context) ([]entities.Project, error) {
	var projects []entities.Project
	err := r.db.WithContext(ctx).
		Order("featured DESC, created_at DESC").
		Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) ListAdmin(ctx context.Context) ([]entities.Project, error) {
	var projects []entities.Project
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Project, error) {
	var project entities.Project
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) Update(ctx context.Context, project *entities.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Project{}, "id = ?", id).Error
}
