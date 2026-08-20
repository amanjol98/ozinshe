package repositories

import (
	"context"
	"errors"
	"ozinshe/internal/models"

	"github.com/jackc/pgx/v5"
)

type MovieGenreRepository struct {
	db *pgx.Conn
}

func NewMovieGenreRepository(db *pgx.Conn) *MovieGenreRepository {
	return &MovieGenreRepository{db: db}
}

var ErrMovieGenreNotFound = errors.New("Нет жанра с таким ID")

func (r *MovieGenreRepository) AddGenreToMovie(ctx context.Context, movieID, genreID int) error {
	sqlQuery := `
	INSERT INTO movie_genres (
	movie_id,
	genre_id
	)
	VALUES ($1, $2);
	`
	_, err := r.db.Exec(ctx, sqlQuery, movieID, genreID)

	return err
}

func (r *MovieGenreRepository) GetGenresOfMovie(ctx context.Context, movieID int) ([]models.Genre, error) {
	sqlQuery := `
	SELECT 
	genres.id,
	genres.name
	FROM movie_genres
	JOIN genres ON movie_genres.genre_id=genres.id
	WHERE movie_genres.movie_id=$1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, movieID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movieGenres []models.Genre

	for rows.Next() {
		var movieGenre models.Genre
		err := rows.Scan(
			&movieGenre.ID,
			&movieGenre.Name,
		)

		if err != nil {
			return nil, err
		}

		movieGenres = append(movieGenres, movieGenre)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movieGenres, nil
}

func (r *MovieGenreRepository) DeleteGenreFromMovie(ctx context.Context, movieID, genreID int) error {
	sqlQuery := `
	DELETE FROM movie_genres
	WHERE movie_id=$1
	AND genre_id=$2;
	`

	result, err := r.db.Exec(ctx, sqlQuery, movieID, genreID)
	if result.RowsAffected() == 0 {
		return ErrMovieGenreNotFound
	}
	return err
}
