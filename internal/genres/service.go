package genres

import (
	"context"
	"fmt"
	"strings"
)

type GenreService struct {
	repo *GenreRepository
}

func NewGenreService(repo *GenreRepository) *GenreService {
	return &GenreService{
		repo: repo,
	}
}

func (s *GenreService) GetAll(ctx context.Context) ([]Genre, error) {
	return s.repo.GetAll(ctx)
}

func (s *GenreService) GetGenreByID(ctx context.Context, id int) (Genre, error) {
	return s.repo.GetGenreByID(ctx, id)
}

func (s *GenreService) CreateGenre(ctx context.Context, name string) (Genre, error) {
	if strings.TrimSpace(name) == "" {
		return Genre{}, fmt.Errorf("Вы ввели пустое значение!")
	}

	return s.repo.CreateGenre(ctx, name)
}

func (s *GenreService) DeleteGenre(ctx context.Context, id int) error {
	return s.repo.DeleteGenre(ctx, id)
}

func (s *GenreService) UpdateGenre(ctx context.Context, name string, id int) (Genre, error) {
	if strings.TrimSpace(name) == "" {
		return Genre{}, fmt.Errorf("Вы ввели пустое значение!")
	}

	return s.repo.UpdateGenre(ctx, name, id)
}
