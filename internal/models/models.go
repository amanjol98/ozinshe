package models

import "time"

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Episode struct {
	ID            int    `json:"id"`
	EpisodeNumber int    `json:"episode_number"`
	VideoID       string `json:"video_id"`
}

type Favorite struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id"`
	MovieID int `json:"movie_id"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

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
	Movie            Movie             `json:"movie"`
	Genres           []Genre           `json:"genres"`
	Categories       []Category        `json:"categories"`
	Seasons          []Season          `json:"seasons"`
	MovieScreenshots []MovieScreenshot `json:"movie_screenshots"`
}

type CreateMovieRequest struct {
	Movie       Movie `json:"movie"`
	GenreIDs    []int `json:"genre_ids"`
	CategoryIDs []int `json:"category_ids"`
}

type CreateMovieModel struct {
	Movie      Movie      `json:"movie"`
	Genres     []Genre    `json:"genres"`
	Categories []Category `json:"categories"`
}

type MovieCategory struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MovieGenre struct {
	MovieID int `json:"movie_id"`
	GenreID int `json:"genre_id"`
}

type MovieScreenshot struct {
	ID       int     `json:"id"`
	ImageURL *string `json:"image_url"`
}

type Season struct {
	ID           int       `json:"id"`
	SeasonNumber int       `json:"season_number"`
	Episodes     []Episode `json:"episodes"`
}

type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	PhoneNumber string    `json:"phone_number"`
	BornAt      time.Time `json:"born_at"`
}
