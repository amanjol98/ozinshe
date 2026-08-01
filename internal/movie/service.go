package movie

import (
	"context"
	"fmt"
	"ozinshe/internal/episodes"
	"ozinshe/internal/movie_categories"
	"ozinshe/internal/movie_genres"
	"ozinshe/internal/movie_screenshots"
	"ozinshe/internal/seasons"
	"strings"
)

type MovieService struct {
	repo *MovieRepository

	genreService    *movie_genres.MovieGenreService
	categoryService *movie_categories.MovieCategoryService

	seasonService           *seasons.SeasonService
	episodeService          *episodes.EpisodeService
	movieScreenshotsService *movie_screenshots.MovieScreenshotsService
}

func NewMovieService(
	repo *MovieRepository,
	genreService *movie_genres.MovieGenreService,
	categoryService *movie_categories.MovieCategoryService,
	seasonService *seasons.SeasonService,
	episodeService *episodes.EpisodeService,
	movieScreenshotsService *movie_screenshots.MovieScreenshotsService,
) *MovieService {
	return &MovieService{
		repo:                    repo,
		genreService:            genreService,
		categoryService:         categoryService,
		seasonService:           seasonService,
		episodeService:          episodeService,
		movieScreenshotsService: movieScreenshotsService,
	}
}

func (s *MovieService) GetAll(ctx context.Context, search string, limit, offset int) ([]Movie, error) {
	search = strings.TrimSpace(search)
	return s.repo.GetAll(ctx, search, limit, offset)
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

	screenshots, err := s.movieScreenshotsService.GetScreenshots(ctx, id)
	if err != nil {
		return MovieResponse{}, err
	}

	response := MovieResponse{
		Movie:            movie,
		Genres:           genres,
		Categories:       categories,
		Seasons:          seasons,
		MovieScreenshots: screenshots,
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
	genreIDs,
	categoryIDs []int,
) (CreateMovieModel, error) {
	if strings.TrimSpace(title) == "" {
		return CreateMovieModel{}, fmt.Errorf("Ввели пустое значение")
	}

	createdMovie, err := s.repo.CreateMovie(
		ctx,
		title,
		releaseYear,
		description,
		duration,
		posterURL,
		director,
		producer,
		videoID,
	)

	if err != nil {
		return CreateMovieModel{}, err
	}

	for _, genreID := range genreIDs {
		err := s.genreService.AddGenreToMovie(ctx, createdMovie.ID, genreID)
		if err != nil {
			return CreateMovieModel{}, err
		}
	}

	for _, categoryID := range categoryIDs {
		err := s.categoryService.AddCategoryToMovie(ctx, createdMovie.ID, categoryID)
		if err != nil {
			return CreateMovieModel{}, err
		}
	}

	genres, err := s.genreService.GetGenresOfMovie(ctx, createdMovie.ID)
	if err != nil {
		return CreateMovieModel{}, err
	}
	categories, err := s.categoryService.GetMovieCategories(ctx, createdMovie.ID)
	if err != nil {
		return CreateMovieModel{}, err
	}

	return CreateMovieModel{
		Movie:      createdMovie,
		Genres:     genres,
		Categories: categories,
	}, nil

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
