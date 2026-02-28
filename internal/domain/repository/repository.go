package repository

import (
	"context"

	"portfolio-backend/internal/domain/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	Create(ctx context.Context, user *entities.User) error
	CountAdmins(ctx context.Context) (int64, error)
}

type ProjectRepository interface {
	Create(ctx context.Context, project *entities.Project) error
	ListPublic(ctx context.Context) ([]entities.Project, error)
	ListAdmin(ctx context.Context) ([]entities.Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Project, error)
	Update(ctx context.Context, project *entities.Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CertificateRepository interface {
	Create(ctx context.Context, certificate *entities.Certificate) error
	List(ctx context.Context) ([]entities.Certificate, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Certificate, error)
	Update(ctx context.Context, certificate *entities.Certificate) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProfileRepository interface {
	Get(ctx context.Context) (*entities.Profile, error)
	Upsert(ctx context.Context, profile *entities.Profile) error
}

type SkillRepository interface {
	Create(ctx context.Context, skill *entities.Skill) error
	List(ctx context.Context) ([]entities.Skill, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Skill, error)
	Update(ctx context.Context, skill *entities.Skill) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SocialLinkRepository interface {
	Create(ctx context.Context, link *entities.SocialLink) error
	ListVisible(ctx context.Context) ([]entities.SocialLink, error)
	ListAdmin(ctx context.Context) ([]entities.SocialLink, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.SocialLink, error)
	Update(ctx context.Context, link *entities.SocialLink) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type TeachingRepository interface {
	Create(ctx context.Context, item *entities.Teaching) error
	ListVisible(ctx context.Context) ([]entities.Teaching, error)
	ListAdmin(ctx context.Context) ([]entities.Teaching, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Teaching, error)
	Update(ctx context.Context, item *entities.Teaching) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type AwardRepository interface {
	Create(ctx context.Context, item *entities.Award) error
	ListVisible(ctx context.Context) ([]entities.Award, error)
	ListAdmin(ctx context.Context) ([]entities.Award, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Award, error)
	Update(ctx context.Context, item *entities.Award) error
	Delete(ctx context.Context, id uuid.UUID) error
}
