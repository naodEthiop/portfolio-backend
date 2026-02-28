package usecase

import (
	"context"
	"errors"

	"portfolio-backend/internal/domain/entities"
	"portfolio-backend/internal/domain/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SocialLinkService struct {
	repo repository.SocialLinkRepository
}

func NewSocialLinkService(repo repository.SocialLinkRepository) *SocialLinkService {
	return &SocialLinkService{repo: repo}
}

type CreateSocialLinkInput struct {
	Platform  string `json:"platform" binding:"required,max=80"`
	URL       string `json:"url" binding:"required,url,max=255"`
	SortOrder int    `json:"sort_order"`
	Visible   bool   `json:"visible"`
}

type UpdateSocialLinkInput struct {
	Platform  string `json:"platform" binding:"required,max=80"`
	URL       string `json:"url" binding:"required,url,max=255"`
	SortOrder int    `json:"sort_order"`
	Visible   bool   `json:"visible"`
}

func (s *SocialLinkService) Create(ctx context.Context, in CreateSocialLinkInput) (*entities.SocialLink, error) {
	link := &entities.SocialLink{
		Platform:  in.Platform,
		URL:       in.URL,
		SortOrder: in.SortOrder,
		Visible:   in.Visible,
	}
	if err := s.repo.Create(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *SocialLinkService) ListVisible(ctx context.Context) ([]entities.SocialLink, error) {
	return s.repo.ListVisible(ctx)
}

func (s *SocialLinkService) ListAdmin(ctx context.Context) ([]entities.SocialLink, error) {
	return s.repo.ListAdmin(ctx)
}

func (s *SocialLinkService) Update(ctx context.Context, id uuid.UUID, in UpdateSocialLinkInput) (*entities.SocialLink, error) {
	link, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	link.Platform = in.Platform
	link.URL = in.URL
	link.SortOrder = in.SortOrder
	link.Visible = in.Visible
	if err := s.repo.Update(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *SocialLinkService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.repo.Delete(ctx, id)
}
