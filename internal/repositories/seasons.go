package repositories

import (
	"context"
	"errors"
	"ozinshe/internal/models"

	"github.com/jackc/pgx/v5"
)

type SeasonRepository struct {
	db *pgx.Conn
}

func NewSeasonRepository(db *pgx.Conn) *SeasonRepository {
	return &SeasonRepository{db: db}
}

var ErrSeasonNotFound = errors.New("Нет сезона с таким ID")

func (r *SeasonRepository) GetSeasons(ctx context.Context, movieID int) ([]models.Season, error) {
	sqlQuery := `
	SELECT
	id,
	season_number
	FROM seasons
	WHERE seasons.movie_id=$1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seasons []models.Season

	for rows.Next() {
		var season models.Season
		err := rows.Scan(
			&season.ID,
			&season.SeasonNumber,
		)

		if err != nil {
			return nil, err
		}

		seasons = append(seasons, season)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return seasons, nil
}

func (r *SeasonRepository) AddSeasonToMovie(ctx context.Context, movieID, seasonNumber int) error {
	sqlQuery := `
	INSERT INTO seasons (movie_id, season_number)
	VALUES($1,$2);
	`

	_, err := r.db.Exec(ctx, sqlQuery, movieID, seasonNumber)

	return err
}

func (r *SeasonRepository) DeleteSeasonFromMovie(ctx context.Context, seasonID int) error {
	sqlQuery := `
	DELETE FROM seasons
	WHERE id=$1;
	`

	result, err := r.db.Exec(ctx, sqlQuery, seasonID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrSeasonNotFound
	}

	return nil
}
