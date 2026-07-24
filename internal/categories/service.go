package categories

import (
	"context"
	"fmt"
	"strings"
)

type CategoryService struct {
	repo *CategoryRepository
}

func NewCategoryService(repo *CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *CategoryService) GetCategoryByID(ctx context.Context, id int) (Category, error) {
	return s.repo.GetCategoryByID(ctx, id)
}

func (s *CategoryService) CreateCategory(ctx context.Context, name string) (Category, error) {
	if strings.TrimSpace(name) == "" {
		return Category{}, fmt.Errorf("Вы ввели пустое название категорий!")
	}
	return s.repo.CreateCategory(ctx, name)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *CategoryService) UpdateCategory(ctx context.Context, name string, id int) (Category, error) {
	if strings.TrimSpace(name) == "" {
		return Category{}, fmt.Errorf("Вы ввели пустое название категорий!")
	}
	return s.repo.UpdateCategory(ctx, name, id)
}
