package usecase

import (
	"context"
	"errors"
	"time"

	"portfolio-backend/internal/domain/entities"
	"portfolio-backend/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeachingService struct {
	repo repository.TeachingRepository
}

func NewTeachingService(repo repository.TeachingRepository) *TeachingService {
	return &TeachingService{repo: repo}
}

type CreateTeachingInput struct {
	Title        string `json:"title" binding:"required,max=180"`
	Organization string `json:"organization" binding:"required,max=180"`
	Location     string `json:"location" binding:"max=160"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Description  string `json:"description"`
	LinkURL      string `json:"link_url" binding:"max=255"`
	SortOrder    int    `json:"sort_order"`
	Visible      bool   `json:"visible"`
}

type UpdateTeachingInput = CreateTeachingInput

func (s *TeachingService) Create(ctx context.Context, in CreateTeachingInput) (*entities.Teaching, error) {
	start, err := parseOptionalDate(in.StartDate)
	if err != nil {
		return nil, err
	}
	end, err := parseOptionalDate(in.EndDate)
	if err != nil {
		return nil, err
	}

	item := &entities.Teaching{
		Title:        in.Title,
		Organization: in.Organization,
		Location:     in.Location,
		StartDate:    start,
		EndDate:      end,
		Description:  in.Description,
		LinkURL:      in.LinkURL,
		SortOrder:    in.SortOrder,
		Visible:      in.Visible,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *TeachingService) ListVisible(ctx context.Context) ([]entities.Teaching, error) {
	return s.repo.ListVisible(ctx)
}

func (s *TeachingService) ListAdmin(ctx context.Context) ([]entities.Teaching, error) {
	return s.repo.ListAdmin(ctx)
}

func (s *TeachingService) Update(ctx context.Context, id uuid.UUID, in UpdateTeachingInput) (*entities.Teaching, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	start, err := parseOptionalDate(in.StartDate)
	if err != nil {
		return nil, err
	}
	end, err := parseOptionalDate(in.EndDate)
	if err != nil {
		return nil, err
	}

	item.Title = in.Title
	item.Organization = in.Organization
	item.Location = in.Location
	item.StartDate = start
	item.EndDate = end
	item.Description = in.Description
	item.LinkURL = in.LinkURL
	item.SortOrder = in.SortOrder
	item.Visible = in.Visible

	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *TeachingService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}

func parseOptionalDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
