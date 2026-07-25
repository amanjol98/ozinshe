package movie_categories

import (
	"context"
	"ozinshe/internal/categories"
)

type MovieCategoryService struct {
	repo *MovieCategoryRepositry
}

func NewMovieCategoryService(repo *MovieCategoryRepositry) *MovieCategoryService {
	return &MovieCategoryService{repo: repo}
}

func (s *MovieCategoryService) AddCategoryToMovie(ctx context.Context, movieID, categoryID int) error {
	return s.repo.AddCategoryToMovie(ctx, movieID, categoryID)
}

func (s *MovieCategoryService) GetMovieCategories(ctx context.Context, movieID int) ([]categories.Category, error) {
	return s.repo.GetMovieCategories(ctx, movieID)
}

func (s *MovieCategoryService) DeleteCategoryFromMovie(ctx context.Context, movieID, categoryID int) error {
	return s.repo.DeleteCategoryFromMovie(ctx, movieID, categoryID)
}
