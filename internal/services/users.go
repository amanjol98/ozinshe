package services

import (
	"context"
	"errors"
	"ozinshe/internal/middleware"
	"ozinshe/internal/models"
	"ozinshe/internal/repositories"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Register(
		ctx context.Context,
		name, email, password, phone_number string,
		born_at time.Time,
	) error

	GetByEmail(
		ctx context.Context,
		email string,
	) (
		models.User,
		error,
	)

	GetByID(ctx context.Context, id int) (models.User, error)
}

type UserService struct {
	repo         UserRepository
	tokenService *middleware.TokenService
}

func NewUserService(repo UserRepository, tokenService *middleware.TokenService) *UserService {
	return &UserService{
		repo:         repo,
		tokenService: tokenService,
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

	token, err := s.tokenService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) GetByID(ctx context.Context, id int) (models.User, error) {
	return s.repo.GetByID(ctx, id)
}
