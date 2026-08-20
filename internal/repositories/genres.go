package repositories

import (
	"context"
	"errors"
	"ozinshe/internal/models"

	"github.com/jackc/pgx/v5"
)

type GenreRepository struct {
	db *pgx.Conn
}

func NewGenreRepository(db *pgx.Conn) *GenreRepository {
	return &GenreRepository{
		db: db,
	}
}

var ErrGenreNotFound = errors.New("Нет жанра с таким ID")

func (r *GenreRepository) GetAll(ctx context.Context) ([]models.Genre, error) {
	sqlQuery := `
	SELECT id, name FROM genres;
	`

	rows, err := r.db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var genres []models.Genre

	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(
			&genre.ID,
			&genre.Name,
		)
		if err != nil {
			return nil, err
		}

		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return genres, nil
}

func (r *GenreRepository) GetGenreByID(ctx context.Context, id int) (models.Genre, error) {
	sqlQuery := `
	SELECT *FROM genres
	WHERE id=$1;
	`
	row := r.db.QueryRow(ctx, sqlQuery, id)
	var genre models.Genre
	err := row.Scan(
		&genre.ID,
		&genre.Name,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Genre{}, ErrGenreNotFound
	}

	if err != nil {
		return models.Genre{}, err
	}

	return genre, nil
}

func (r *GenreRepository) CreateGenre(ctx context.Context, name string) (models.Genre, error) {
	sqlQuery := `
	INSERT INTO genres(name)
	VALUES ($1)
	RETURNING
		id,
		name;
	`

	var genre models.Genre
	row := r.db.QueryRow(ctx, sqlQuery, name)
	err := row.Scan(
		&genre.ID,
		&genre.Name,
	)

	if err != nil {
		return models.Genre{}, err
	}

	return genre, nil
}

func (r *GenreRepository) DeleteGenre(ctx context.Context, id int) error {
	sqlQuery := `
	DELETE FROM genres
	WHERE id=$1;
	`
	result, err := r.db.Exec(ctx, sqlQuery, id)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrGenreNotFound
	}

	return nil
}

func (r *GenreRepository) UpdateGenre(ctx context.Context, name string, id int) (models.Genre, error) {
	sqlQuery := `
	UPDATE genres
	SET name=$1
	WHERE id=$2
	RETURNING
		id,
		name;
	`
	var genre models.Genre

	row := r.db.QueryRow(ctx, sqlQuery, name, id)
	err := row.Scan(
		&genre.ID,
		&genre.Name,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return models.Genre{}, ErrGenreNotFound
	}

	if err != nil {
		return models.Genre{}, err
	}

	return genre, nil
}
