package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type CategoryRepository struct {
	db *pgx.Conn
}

func NewCategoryRepository(db *pgx.Conn) *CategoryRepository {
	return &CategoryRepository{db: db}
}

var ErrCategoryNotFound = errors.New("Нет категорий с таким ID")

func (r *CategoryRepository) GetCategories(ctx context.Context) ([]Category, error) {
	sqlQuery := `
	SELECT id, name
	FROM categories;
	`
	rows, err := r.db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var category Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
		)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepository) GetCategoryByID(ctx context.Context, id int) (Category, error) {
	sqlQuery := `
	SELECT id, name
	FROM categories
	WHERE id=$1;
	`

	var category Category

	err := r.db.QueryRow(ctx, sqlQuery, id).Scan(
		&category.ID,
		&category.Name,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Category{}, ErrCategoryNotFound
	}

	if err != nil {
		return Category{}, err
	}

	return category, nil
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, name string) (Category, error) {
	sqlQuery := `
	INSERT INTO categories (name)
	VALUES ($1)
	RETURNING
	id,
	name;
	`

	var category Category

	err := r.db.QueryRow(ctx, sqlQuery, name).Scan(
		&category.ID,
		&category.Name,
	)

	if err != nil {
		return Category{}, err
	}
	return category, nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int) error {
	sqlQuery := `
	DELETE FROM categories
	WHERE id=$1;
	`

	result, err := r.db.Exec(ctx, sqlQuery, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, name string, id int) (Category, error) {
	sqlQuery := `
	UPDATE categories
	SET name=$1
	WHERE id=$2
	RETURNING
	id,
	name;
	`

	var category Category

	err := r.db.QueryRow(ctx, sqlQuery, name, id).Scan(
		&category.ID,
		&category.Name,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Category{}, ErrCategoryNotFound
	}
	if err != nil {
		return Category{}, err
	}

	return category, nil
}
