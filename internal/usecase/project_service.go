package usecase

import (
	"context"
	"errors"
	"mime/multipart"

	"portfolio-backend/internal/domain/entities"
	"portfolio-backend/internal/domain/repository"
	"portfolio-backend/internal/infrastructure/storage"
	"portfolio-backend/pkg/validator"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ProjectService struct {
	repo    repository.ProjectRepository
	storage storage.ImageStorage
}

func NewProjectService(repo repository.ProjectRepository, storage storage.ImageStorage) *ProjectService {
	return &ProjectService{repo: repo, storage: storage}
}

type CreateProjectInput struct {
	Title            string   `json:"title" binding:"required,min=3,max=140"`
	ShortDescription string   `json:"short_description" binding:"required,max=255"`
	Description      string   `json:"description" binding:"required"`
	Status           string   `json:"status" binding:"required,oneof=complete in_progress archived"`
	TechStack        []string `json:"tech_stack"`
	Achievements     []string `json:"achievements"`
	DemoURL          string   `json:"demo_url"`
	RepoURL          string   `json:"repo_url"`
	Featured         bool     `json:"featured"`
}

type UpdateProjectInput struct {
	Title            string   `json:"title" binding:"required,min=3,max=140"`
	ShortDescription string   `json:"short_description" binding:"required,max=255"`
	Description      string   `json:"description" binding:"required"`
	Status           string   `json:"status" binding:"required,oneof=complete in_progress archived"`
	TechStack        []string `json:"tech_stack"`
	Achievements     []string `json:"achievements"`
	DemoURL          string   `json:"demo_url"`
	RepoURL          string   `json:"repo_url"`
	Featured         bool     `json:"featured"`
}

func (s *ProjectService) Create(ctx context.Context, in CreateProjectInput) (*entities.Project, error) {
	project := &entities.Project{
		Title:            in.Title,
		Slug:             validator.Slug(in.Title),
		ShortDescription: in.ShortDescription,
		Description:      in.Description,
		Status:           in.Status,
		TechStack:        pq.StringArray(in.TechStack),
		Achievements:     pq.StringArray(in.Achievements),
		DemoURL:          in.DemoURL,
		RepoURL:          in.RepoURL,
		Featured:         in.Featured,
	}
	if project.Slug == "" {
		project.Slug = uuid.NewString()
	}

	if err := s.repo.Create(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) ListPublic(ctx context.Context) ([]entities.Project, error) {
	return s.repo.ListPublic(ctx)
}

func (s *ProjectService) ListAdmin(ctx context.Context) ([]entities.Project, error) {
	return s.repo.ListAdmin(ctx)
}

func (s *ProjectService) GetByID(ctx context.Context, id uuid.UUID) (*entities.Project, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) Update(ctx context.Context, id uuid.UUID, in UpdateProjectInput) (*entities.Project, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	project.Title = in.Title
	project.Slug = validator.Slug(in.Title)
	if project.Slug == "" {
		project.Slug = id.String()
	}
	project.ShortDescription = in.ShortDescription
	project.Description = in.Description
	project.Status = in.Status
	project.TechStack = pq.StringArray(in.TechStack)
	project.Achievements = pq.StringArray(in.Achievements)
	project.DemoURL = in.DemoURL
	project.RepoURL = in.RepoURL
	project.Featured = in.Featured

	if err := s.repo.Update(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *ProjectService) SetFeatured(ctx context.Context, id uuid.UUID, featured bool) (*entities.Project, error) {
	project, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	project.Featured = featured
	if err := s.repo.Update(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) UploadImage(ctx context.Context, id uuid.UUID, fileHeader *multipart.FileHeader) (*entities.Project, error) {
	project, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	path, err := s.storage.SaveImage("projects", fileHeader)
	if err != nil {
		return nil, err
	}
	project.ImageURL = path
	if err := s.repo.Update(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}
