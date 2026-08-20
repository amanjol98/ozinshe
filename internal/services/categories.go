package services

import (
	"context"
	"fmt"
	"ozinshe/internal/models"
	"ozinshe/internal/repositories"
	"strings"
)

type CategoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]models.Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (models.Category, error) {
	return s.repo.GetCategoryByID(ctx, id)
}

func (s *CategoryService) CreateCategory(ctx context.Context, name string) (models.Category, error) {
	if strings.TrimSpace(name) == "" {
		return models.Category{}, fmt.Errorf("Вы ввели пустое название категорий!")
	}
	return s.repo.CreateCategory(ctx, name)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, name string, id int) (models.Category, error) {
	if strings.TrimSpace(name) == "" {
		return models.Category{}, fmt.Errorf("Вы ввели пустое название категорий!")
	}
	return s.repo.UpdateCategory(ctx, name, id)
}
