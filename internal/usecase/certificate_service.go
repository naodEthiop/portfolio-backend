package usecase

import (
	"context"
	"errors"
	"mime/multipart"
	"time"

	"portfolio-backend/internal/domain/entities"
	"portfolio-backend/internal/domain/repository"
	"portfolio-backend/internal/infrastructure/storage"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CertificateService struct {
	repo    repository.CertificateRepository
	storage storage.ImageStorage
}

func NewCertificateService(repo repository.CertificateRepository, storage storage.ImageStorage) *CertificateService {
	return &CertificateService{repo: repo, storage: storage}
}

type CreateCertificateInput struct {
	Name          string `json:"name" binding:"required,max=180"`
	Issuer        string `json:"issuer" binding:"required,max=180"`
	IssueDate     string `json:"issue_date"`
	Description   string `json:"description"`
	CredentialURL string `json:"credential_url"`
}

type UpdateCertificateInput struct {
	Name          string `json:"name" binding:"required,max=180"`
	Issuer        string `json:"issuer" binding:"required,max=180"`
	IssueDate     string `json:"issue_date"`
	Description   string `json:"description"`
	CredentialURL string `json:"credential_url"`
}

func parseDate(val string) (*time.Time, error) {
	if val == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", val)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *CertificateService) Create(ctx context.Context, in CreateCertificateInput) (*entities.Certificate, error) {
	issueDate, err := parseDate(in.IssueDate)
	if err != nil {
		return nil, err
	}
	certificate := &entities.Certificate{
		Name:          in.Name,
		Issuer:        in.Issuer,
		IssueDate:     issueDate,
		Description:   in.Description,
		CredentialURL: in.CredentialURL,
	}
	if err := s.repo.Create(ctx, certificate); err != nil {
		return nil, err
	}
	return certificate, nil
}

func (s *CertificateService) List(ctx context.Context) ([]entities.Certificate, error) {
	return s.repo.List(ctx)
}

func (s *CertificateService) GetByID(ctx context.Context, id uuid.UUID) (*entities.Certificate, error) {
	certificate, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return certificate, nil
}

func (s *CertificateService) Update(ctx context.Context, id uuid.UUID, in UpdateCertificateInput) (*entities.Certificate, error) {
	certificate, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	issueDate, err := parseDate(in.IssueDate)
	if err != nil {
		return nil, err
	}

	certificate.Name = in.Name
	certificate.Issuer = in.Issuer
	certificate.IssueDate = issueDate
	certificate.Description = in.Description
	certificate.CredentialURL = in.CredentialURL

	if err := s.repo.Update(ctx, certificate); err != nil {
		return nil, err
	}
	return certificate, nil
}

func (s *CertificateService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func (s *CertificateService) UploadImage(ctx context.Context, id uuid.UUID, fileHeader *multipart.FileHeader) (*entities.Certificate, error) {
	certificate, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	path, err := s.storage.SaveImage("certificates", fileHeader)
	if err != nil {
		return nil, err
	}
	certificate.ImageURL = path
	if err := s.repo.Update(ctx, certificate); err != nil {
		return nil, err
	}
	return certificate, nil
}
