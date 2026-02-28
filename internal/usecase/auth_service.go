package usecase

import (
	"context"
	"errors"

	"portfolio-backend/internal/domain/entities"
	"portfolio-backend/internal/domain/repository"
	"portfolio-backend/pkg/auth"
	"portfolio-backend/pkg/password"

	"gorm.io/gorm"
)

type AuthService struct {
	users repository.UserRepository
	jwt   *auth.JWTManager
}

func NewAuthService(users repository.UserRepository, jwt *auth.JWTManager) *AuthService {
	return &AuthService{users: users, jwt: jwt}
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (*AuthResponse, error) {
	user, err := s.users.GetByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	if err := password.Verify(user.PasswordHash, in.Password); err != nil {
		return nil, ErrUnauthorized
	}

	token, err := s.jwt.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, Role: user.Role}, nil
}

func (s *AuthService) BootstrapAdmin(ctx context.Context, email, rawPassword string) error {
	count, err := s.users.CountAdmins(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, err := password.Hash(rawPassword)
	if err != nil {
		return err
	}

	admin := &entities.User{
		Email:        email,
		PasswordHash: hash,
		Role:         entities.RoleAdmin,
	}

	return s.users.Create(ctx, admin)
}
