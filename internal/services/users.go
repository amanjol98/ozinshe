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
		email, passwordstring string,
	) error

	GetByEmail(
		ctx context.Context,
		email string,
	) (
		models.User,
		error,
	)

	GetByID(ctx context.Context, id int) (models.User, error)

	Update(
		ctx context.Context,
		userID int,
		name string,
		phoneNumber string,
		bornAt *time.Time,
	) error

	UpdatePassword(ctx context.Context, userID int, password string) error

	GetByIDWithPassword(ctx context.Context, userID int) (models.User, error)
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
	email, password string,
) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.Register(ctx, email, string(hash))
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

func (s *UserService) Update(
	ctx context.Context,
	userID int,
	name, phoneNumber string,
	bornAt *time.Time,
) error {
	return s.repo.Update(ctx, userID, name, phoneNumber, bornAt)
}

func (s *UserService) UpdatePassword(ctx context.Context, userID int, oldPassword string, newPassword string) error {
	user, err := s.repo.GetByIDWithPassword(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("старый пароль указан неверно")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(ctx, userID, string(hashedPassword))
}
