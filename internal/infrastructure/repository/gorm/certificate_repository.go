package gormrepo

import (
	"context"

	"portfolio-backend/internal/domain/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CertificateRepository struct {
	db *gorm.DB
}

func NewCertificateRepository(db *gorm.DB) *CertificateRepository {
	return &CertificateRepository{db: db}
}

func (r *CertificateRepository) Create(ctx context.Context, certificate *entities.Certificate) error {
	return r.db.WithContext(ctx).Create(certificate).Error
}

func (r *CertificateRepository) List(ctx context.Context) ([]entities.Certificate, error) {
	var certificates []entities.Certificate
	err := r.db.WithContext(ctx).Order("issue_date DESC NULLS LAST, created_at DESC").Find(&certificates).Error
	return certificates, err
}

func (r *CertificateRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.Certificate, error) {
	var certificate entities.Certificate
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&certificate).Error
	if err != nil {
		return nil, err
	}
	return &certificate, nil
}

func (r *CertificateRepository) Update(ctx context.Context, certificate *entities.Certificate) error {
	return r.db.WithContext(ctx).Save(certificate).Error
}

func (r *CertificateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Certificate{}, "id = ?", id).Error
}
