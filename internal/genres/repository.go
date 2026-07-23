package genres

import (
	"context"

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

func (r *GenreRepository) GetAll(ctx context.Context) ([]Genre, error) {
	sqlQuery := `
	SELECT *FROM genres;
	`

	rows, err := r.db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	var genres []Genre

	for rows.Next() {
		var genre Genre
		err := rows.Scan(
			&genre.ID,
			&genre.Name,
		)
		if err != nil {
			return nil, err
		}

		genres = append(genres, genre)
	}

	return genres, nil
}

func (r *GenreRepository) CreateGenre(ctx context.Context, name string) (Genre, error) {
	sqlQuery := `
	INSERT INTO genres(name)
	VALUES ($1)
	RETURNING
		id,
		name;
	`

	var genre Genre
	row := r.db.QueryRow(ctx, sqlQuery, name)
	err := row.Scan(
		&genre.ID,
		&genre.Name,
	)

	if err != nil {
		return Genre{}, err
	}

	return genre, nil

}
