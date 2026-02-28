package usecase

import (
	"context"
	"errors"

	"portfolio-backend/internal/domain/entities"
	"portfolio-backend/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SkillService struct {
	repo repository.SkillRepository
}

func NewSkillService(repo repository.SkillRepository) *SkillService {
	return &SkillService{repo: repo}
}

type CreateSkillInput struct {
	Category  string `json:"category" binding:"required,max=100"`
	Name      string `json:"name" binding:"required,max=120"`
	SortOrder int    `json:"sort_order"`
}

type UpdateSkillInput struct {
	Category  string `json:"category" binding:"required,max=100"`
	Name      string `json:"name" binding:"required,max=120"`
	SortOrder int    `json:"sort_order"`
}

func (s *SkillService) Create(ctx context.Context, in CreateSkillInput) (*entities.Skill, error) {
	skill := &entities.Skill{Category: in.Category, Name: in.Name, SortOrder: in.SortOrder}
	if err := s.repo.Create(ctx, skill); err != nil {
		return nil, err
	}
	return skill, nil
}

func (s *SkillService) List(ctx context.Context) ([]entities.Skill, error) {
	return s.repo.List(ctx)
}

func (s *SkillService) Update(ctx context.Context, id uuid.UUID, in UpdateSkillInput) (*entities.Skill, error) {
	skill, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	skill.Category = in.Category
	skill.Name = in.Name
	skill.SortOrder = in.SortOrder
	if err := s.repo.Update(ctx, skill); err != nil {
		return nil, err
	}
	return skill, nil
}

func (s *SkillService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}
