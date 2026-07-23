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

func (s *GenreService) CreateGenre(ctx context.Context, name string) (Genre, error) {
	if strings.TrimSpace(name) == "" {
		return Genre{}, fmt.Errorf("Вы ввели пустое значение!")
	}

	return s.repo.CreateGenre(ctx, name)
}
