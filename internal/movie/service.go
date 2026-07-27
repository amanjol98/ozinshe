package movie

import (
	"context"
	"fmt"
	"ozinshe/internal/episodes"
	"ozinshe/internal/movie_categories"
	"ozinshe/internal/movie_genres"
	"ozinshe/internal/seasons"
	"strings"
)

type MovieService struct {
	repo *MovieRepository

	genreService    *movie_genres.MovieGenreService
	categoryService *movie_categories.MovieCategoryService

	seasonService  *seasons.SeasonService
	episodeService *episodes.EpisodeService
}

func NewMovieService(
	repo *MovieRepository,
	genreService *movie_genres.MovieGenreService,
	categoryService *movie_categories.MovieCategoryService,
	seasonService *seasons.SeasonService,
	episodeService *episodes.EpisodeService,
) *MovieService {
	return &MovieService{
		repo:            repo,
		genreService:    genreService,
		categoryService: categoryService,
		seasonService:   seasonService,
		episodeService:  episodeService,
	}
}

func (s *MovieService) GetAll(ctx context.Context) ([]Movie, error) {
	return s.repo.GetAll(ctx)
}

func (s *MovieService) GetByID(ctx context.Context, id int) (MovieResponse, error) {
	movie, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return MovieResponse{}, err
	}
	genres, err := s.genreService.GetGenresOfMovie(ctx, id)
	if err != nil {
		return MovieResponse{}, err
	}
	categories, err := s.categoryService.GetMovieCategories(ctx, id)
	if err != nil {
		return MovieResponse{}, err
	}

	seasons, err := s.seasonService.GetSeasons(ctx, id)
	if err != nil {
		return MovieResponse{}, err
	}

	for i := range seasons {
		episodes, err := s.episodeService.GetEpisodes(ctx, seasons[i].ID)
		if err != nil {
			return MovieResponse{}, err
		}

		seasons[i].Episodes = episodes
	}

	response := MovieResponse{
		Movie:      movie,
		Genres:     genres,
		Categories: categories,
		Seasons:    seasons,
	}
	return response, nil
}

func (s *MovieService) CreateMovie(
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
	if strings.TrimSpace(title) == "" {
		return Movie{}, fmt.Errorf("Ввели пустое значение")
	}

	movie := Movie{
		Title:       title,
		ReleaseYear: releaseYear,
		Description: description,
		Duration:    duration,
		PosterURL:   posterURL,
		Director:    director,
		Producer:    producer,
		VideoID:     videoID,
	}

	return s.repo.CreateMovie(
		ctx,
		movie.Title,
		movie.ReleaseYear,
		movie.Description,
		movie.Duration,
		movie.PosterURL,
		movie.Director,
		movie.Producer,
		movie.VideoID,
	)

}

func (s *MovieService) DeleteMovie(ctx context.Context, id int) error {
	return s.repo.DeleteMovie(ctx, id)
}

func (s *MovieService) UpdateMovie(ctx context.Context, movie Movie, id int) (Movie, error) {
	if strings.TrimSpace(movie.Title) == "" {
		return Movie{}, fmt.Errorf("Ввели пустое значение")
	}

	return s.repo.UpdateMovie(ctx, movie, id)
}
