package usecase

import (
	"context"
	"errors"

	"portfolio-backend/internal/domain/entities"
	"portfolio-backend/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AwardService struct {
	repo repository.AwardRepository
}

func NewAwardService(repo repository.AwardRepository) *AwardService {
	return &AwardService{repo: repo}
}

type CreateAwardInput struct {
	Title       string `json:"title" binding:"required,max=180"`
	Issuer      string `json:"issuer" binding:"required,max=180"`
	AwardDate   string `json:"award_date"`
	Description string `json:"description"`
	LinkURL     string `json:"link_url" binding:"max=255"`
	SortOrder   int    `json:"sort_order"`
	Visible     bool   `json:"visible"`
}

type UpdateAwardInput = CreateAwardInput

func (s *AwardService) Create(ctx context.Context, in CreateAwardInput) (*entities.Award, error) {
	date, err := parseOptionalDate(in.AwardDate)
	if err != nil {
		return nil, err
	}
	item := &entities.Award{
		Title:       in.Title,
		Issuer:      in.Issuer,
		AwardDate:   date,
		Description: in.Description,
		LinkURL:     in.LinkURL,
		SortOrder:   in.SortOrder,
		Visible:     in.Visible,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *AwardService) ListVisible(ctx context.Context) ([]entities.Award, error) {
	return s.repo.ListVisible(ctx)
}

func (s *AwardService) ListAdmin(ctx context.Context) ([]entities.Award, error) {
	return s.repo.ListAdmin(ctx)
}

func (s *AwardService) Update(ctx context.Context, id uuid.UUID, in UpdateAwardInput) (*entities.Award, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	date, err := parseOptionalDate(in.AwardDate)
	if err != nil {
		return nil, err
	}

	item.Title = in.Title
	item.Issuer = in.Issuer
	item.AwardDate = date
	item.Description = in.Description
	item.LinkURL = in.LinkURL
	item.SortOrder = in.SortOrder
	item.Visible = in.Visible

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *AwardService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}
