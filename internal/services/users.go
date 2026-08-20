package services

import (
	"context"
	"errors"
	"ozinshe/internal/middleware"
	"ozinshe/internal/repositories"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo         *repositories.UserRepository
	tokenService *middleware.TokenService
}

func NewUserService(repo *repositories.UserRepository, tokenservice *middleware.TokenService) *UserService {
	return &UserService{
		repo:         repo,
		tokenService: tokenservice,
	}
}

var ErrInvalidCredentials = errors.New("неверный email или пароль")

func (s *UserService) Register(
	ctx context.Context,
	name, email, password, phone_number string,
	born_at time.Time,
) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.Register(ctx, name, email, string(hash), phone_number, born_at)
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.tokenService.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
