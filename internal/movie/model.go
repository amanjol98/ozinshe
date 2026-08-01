package movie

import (
	"ozinshe/internal/categories"
	"ozinshe/internal/genres"
	"ozinshe/internal/movie_screenshots"
	"ozinshe/internal/seasons"
)

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"release_year"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	PosterURL   string `json:"poster_url"`
	Director    string `json:"director"`
	Producer    string `json:"producer"`
	VideoID     string `json:"video_id"`
}

type MovieResponse struct {
	Movie            Movie                               `json:"movie"`
	Genres           []genres.Genre                      `json:"genres"`
	Categories       []categories.Category               `json:"categories"`
	Seasons          []seasons.Season                    `json:"seasons"`
	MovieScreenshots []movie_screenshots.MovieScreenshot `json:"movie_screenshots"`
}

type CreateMovieRequest struct {
	Movie       Movie `json:"movie"`
	GenreIDs    []int `json:"genre_ids"`
	CategoryIDs []int `json:"category_ids"`
}

type CreateMovieModel struct {
	Movie      Movie                 `json:"movie"`
	Genres     []genres.Genre        `json:"genres"`
	Categories []categories.Category `json:"categories"`
}
