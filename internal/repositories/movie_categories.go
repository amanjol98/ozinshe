package repositories

import (
	"context"
	"errors"
	"ozinshe/internal/models"

	"github.com/jackc/pgx/v5"
)

type MovieCategoryRepositry struct {
	db *pgx.Conn
}

func NewMovieCategoryRepositry(db *pgx.Conn) *MovieCategoryRepositry {
	return &MovieCategoryRepositry{db: db}
}

var ErrMovieCategoryNotFound = errors.New("Нет фильма с такой категорией")

func (r *MovieCategoryRepositry) AddCategoryToMovie(ctx context.Context, movieID, categoryID int) error {
	sqlQuery := `
	INSERT INTO movie_categories (movie_id, category_id)
	VALUES ($1,$2);
	`

	_, err := r.db.Exec(ctx, sqlQuery, movieID, categoryID)

	return err
}

func (r *MovieCategoryRepositry) GetMovieCategories(ctx context.Context, movieID int) ([]models.Category, error) {
	sqlQuery := `
	SELECT
	categories.id,
	categories.name
	FROM movie_categories
	JOIN categories ON movie_categories.category_id=categories.id
	WHERE movie_categories.movie_id=$1
	`

	rows, err := r.db.Query(ctx, sqlQuery, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movieCategories []models.Category

	for rows.Next() {
		var movieCategory models.Category

		err := rows.Scan(
			&movieCategory.ID,
			&movieCategory.Name,
		)

		if err != nil {
			return nil, err
		}

		movieCategories = append(movieCategories, movieCategory)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movieCategories, nil

}

func (r *MovieCategoryRepositry) DeleteCategoryFromMovie(ctx context.Context, movieID, categoryID int) error {
	sqlQuery := `
	DELETE FROM movie_categories
	WHERE movie_id=$1
	AND category_id=$2;
	`
	result, err := r.db.Exec(ctx, sqlQuery, movieID, categoryID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrMovieCategoryNotFound
	}

	return nil
}

func (r *MovieCategoryRepositry) DeleteAllCategoriesFromMovie(ctx context.Context, movieID int) error {
	sqlQuery := `
	DELETE FROM movie_categories
	WHERE movie_id=$1;
	`

	result, err := r.db.Exec(ctx, sqlQuery, movieID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrMovieNotFound
	}
	return nil
}
