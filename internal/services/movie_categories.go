package services

import (
	"context"
	"ozinshe/internal/models"
	"ozinshe/internal/repositories"
)

type MovieCategoryService struct {
	repo *repositories.MovieCategoryRepositry
}

func NewMovieCategoryService(repo *repositories.MovieCategoryRepositry) *MovieCategoryService {
	return &MovieCategoryService{repo: repo}
}

func (s *MovieCategoryService) AddCategoryToMovie(ctx context.Context, movieID, categoryID int) error {
	return s.repo.AddCategoryToMovie(ctx, movieID, categoryID)
}

func (s *MovieCategoryService) GetMovieCategories(ctx context.Context, movieID int) ([]models.Category, error) {
	return s.repo.GetMovieCategories(ctx, movieID)
}

func (s *MovieCategoryService) DeleteCategoryFromMovie(ctx context.Context, movieID, categoryID int) error {
	return s.repo.DeleteCategoryFromMovie(ctx, movieID, categoryID)
}
