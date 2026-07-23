package movie

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type MovieRepository struct {
	db *pgx.Conn
}

func NewMovieRepository(db *pgx.Conn) *MovieRepository {
	return &MovieRepository{db: db}
}

var ErrMovieNotFound = errors.New("Нет фильма с таким ID")

func (r *MovieRepository) GetAll(ctx context.Context) ([]Movie, error) {
	sqlQuery := `
	SELECT 
	id,
	title,
	release_year,
	description,
	duration,
	poster_url,
	director,
	producer,
	video_id
	FROM movies;
	`
	rows, err := r.db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var movie Movie
		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Description,
			&movie.Duration,
			&movie.PosterURL,
			&movie.Director,
			&movie.Producer,
			&movie.VideoID,
		); err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil

}

func (r *MovieRepository) GetByID(ctx context.Context, id int) (Movie, error) {
	sqlQuery := `
	SELECT 
	id,
	title,
	release_year,
	description,
	duration,
	poster_url,
	director,
	producer,
	video_id
	FROM movies
	WHERE id=$1;
	`

	var movie Movie

	err := r.db.QueryRow(ctx, sqlQuery, id).Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseYear,
		&movie.Description,
		&movie.Duration,
		&movie.PosterURL,
		&movie.Director,
		&movie.Producer,
		&movie.VideoID,
	)

	if err != nil {
		return Movie{}, nil
	}

	return movie, nil
}

func (r *MovieRepository) CreateMovie(
	ctx context.Context,
	title string,
	releaseYear int,
	description string,
	duration int,
	posterURL string,
	director string,
	producer string,
	videoID string,
) (Movie, error) {
	sqlQuery := `
	INSERT INTO movies (
	title,
	release_year,
	description,
	duration,
	poster_url,
	director,
	producer,
	video_id
	)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING
	id,
	title,
	release_year,
	description,
	duration,
	poster_url,
	director,
	producer,
	video_id;
	`
	var movie Movie

	err := r.db.QueryRow(
		ctx,
		sqlQuery,
		title,
		releaseYear,
		description,
		duration,
		posterURL,
		director,
		producer,
		videoID,
	).Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseYear,
		&movie.Description,
		&movie.Duration,
		&movie.PosterURL,
		&movie.Director,
		&movie.Producer,
		&movie.VideoID,
	)

	return movie, err
}

func (r *MovieRepository) DeleteMovie(ctx context.Context, id int) error {
	sqlQuery := `
	DELETE FROM movies WHERE id=$1;
	`
	result, err := r.db.Exec(ctx, sqlQuery, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrMovieNotFound
	}

	return nil
}

func (r *MovieRepository) UpdateMovie(
	ctx context.Context,
	movie Movie,
	id int,
) (Movie, error) {
	sqlQuery := `
	UPDATE movies
	SET 
		title=$1,
		release_year=$2,
		description=$3,
		duration=$4,
		poster_url=$5,
		director=$6,
		producer=$7,
		video_id=$8
	WHERE id=$9
	RETURNING
		id,
		title,
		release_year,
		description,
		duration,
		poster_url,
		director,
		producer,
		video_id;
	`
	var updatedMovie Movie

	err := r.db.QueryRow(
		ctx,
		sqlQuery,
		movie.Title,
		movie.ReleaseYear,
		movie.Description,
		movie.Duration,
		movie.PosterURL,
		movie.Director,
		movie.Producer,
		movie.VideoID,
		id,
	).Scan(
		&updatedMovie.ID,
		&updatedMovie.Title,
		&updatedMovie.ReleaseYear,
		&updatedMovie.Description,
		&updatedMovie.Duration,
		&updatedMovie.PosterURL,
		&updatedMovie.Director,
		&updatedMovie.Producer,
		&updatedMovie.VideoID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Movie{}, ErrMovieNotFound
	}

	if err != nil {
		return Movie{}, err
	}

	return updatedMovie, err
}
