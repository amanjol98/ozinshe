package repositories

import (
	"context"
	"errors"
	"ozinshe/internal/models"

	"github.com/jackc/pgx/v5"
)

type MovieScreenshotsRepository struct {
	db *pgx.Conn
}

func NewMovieScreenshotsRepository(db *pgx.Conn) *MovieScreenshotsRepository {
	return &MovieScreenshotsRepository{db: db}
}

var ErrScreenshotNotFound = errors.New("Нет скриншота с такой ID")

func (r *MovieScreenshotsRepository) AddScreenshotsToMovie(
	ctx context.Context,
	movieID int,
	imageURL string,
) error {
	sqlQuery := `
	INSERT INTO movie_screenshots (
	movie_id,
	image_url
	)
	VALUES($1, $2);
	`

	_, err := r.db.Exec(ctx, sqlQuery, movieID, imageURL)

	return err
}

func (r *MovieScreenshotsRepository) GetScreenshots(ctx context.Context, movieID int) ([]models.MovieScreenshot, error) {
	sqlQuery := `
	SELECT
	ms.id,
	ms.image_url
	FROM movie_screenshots ms
	JOIN movies m ON ms.movie_id=m.id
	WHERE ms.movie_id=$1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, movieID)
	if err != nil {
		return nil, err
	}

	var movieScreenshots []models.MovieScreenshot

	for rows.Next() {
		var movieScreenshot models.MovieScreenshot

		err := rows.Scan(
			&movieScreenshot.ID,
			&movieScreenshot.ImageURL,
		)
		if err != nil {
			return nil, err
		}

		movieScreenshots = append(movieScreenshots, movieScreenshot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movieScreenshots, nil
}

func (r *MovieScreenshotsRepository) DeleteScreenshot(ctx context.Context, screenshotID int) error {
	sqlQuery := `
	DELETE FROM movie_screenshots
	WHERE id=$1;
	`

	result, err := r.db.Exec(ctx, sqlQuery, screenshotID)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrScreenshotNotFound
	}

	return err
}
