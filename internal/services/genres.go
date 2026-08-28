package services

import (
	"context"
	"fmt"
	"ozinshe/internal/models"
	"ozinshe/internal/repositories"
	"strings"
)

type GenreService struct {
	repo *repositories.GenreRepository
}

func NewGenreService(repo *repositories.GenreRepository) *GenreService {
	return &GenreService{
		repo: repo,
	}
}

func (s *GenreService) GetAll(ctx context.Context) ([]models.Genre, error) {
	return s.repo.GetAll(ctx)
}

func (s *GenreService) GetGenreByID(ctx context.Context, id int) (models.Genre, error) {
	return s.repo.GetGenreByID(ctx, id)
}

func (s *GenreService) CreateGenre(ctx context.Context, name string) (models.Genre, error) {
	if strings.TrimSpace(name) == "" {
		return models.Genre{}, fmt.Errorf("Вы ввели пустое значение!")
	}

	return s.repo.CreateGenre(ctx, name)
}

func (s *GenreService) DeleteGenre(ctx context.Context, id int) error {
	return s.repo.DeleteGenre(ctx, id)
}

func (s *GenreService) UpdateGenre(ctx context.Context, name string, id int) (models.Genre, error) {
	if strings.TrimSpace(name) == "" {
		return models.Genre{}, fmt.Errorf("Вы ввели пустое значение!")
	}

	return s.repo.UpdateGenre(ctx, name, id)
}
