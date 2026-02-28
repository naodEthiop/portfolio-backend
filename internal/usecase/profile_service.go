package usecase

import (
	"context"
	"errors"

	"portfolio-backend/internal/domain/entities"
	"portfolio-backend/internal/domain/repository"

	"gorm.io/gorm"
)

type ProfileService struct {
	repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) *ProfileService {
	return &ProfileService{repo: repo}
}

type UpsertProfileInput struct {
	FullName     string `json:"full_name" binding:"required,max=150"`
	Handle       string `json:"handle" binding:"max=80"`
	Headline     string `json:"headline" binding:"max=200"`
	Location     string `json:"location" binding:"max=160"`
	Summary      string `json:"summary" binding:"max=500"`
	Bio          string `json:"bio"`
	AvatarURL    string `json:"avatar_url"`
	BannerURL    string `json:"banner_url"`
	ResumeURL    string `json:"resume_url"`
	CTAPrimary   string `json:"cta_primary" binding:"max=120"`
	CTASecondary string `json:"cta_secondary" binding:"max=120"`
	CTATertiary  string `json:"cta_tertiary" binding:"max=120"`
}

func (s *ProfileService) Get(ctx context.Context) (*entities.Profile, error) {
	profile, err := s.repo.Get(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return profile, nil
}

func (s *ProfileService) Upsert(ctx context.Context, in UpsertProfileInput) (*entities.Profile, error) {
	profile := &entities.Profile{
		FullName:     in.FullName,
		Handle:       in.Handle,
		Headline:     in.Headline,
		Location:     in.Location,
		Summary:      in.Summary,
		Bio:          in.Bio,
		AvatarURL:    in.AvatarURL,
		BannerURL:    in.BannerURL,
		ResumeURL:    in.ResumeURL,
		CTAPrimary:   in.CTAPrimary,
		CTASecondary: in.CTASecondary,
		CTATertiary:  in.CTATertiary,
	}
	if err := s.repo.Upsert(ctx, profile); err != nil {
		return nil, err
	}
	return s.repo.Get(ctx)
}
