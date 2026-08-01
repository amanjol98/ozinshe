package main

import (
	"context"
	"log"
	"net/http"
	"ozinshe/internal/categories"
	"ozinshe/internal/database"
	"ozinshe/internal/episodes"
	"ozinshe/internal/genres"
	"ozinshe/internal/movie"
	"ozinshe/internal/movie_categories"
	"ozinshe/internal/movie_genres"
	"ozinshe/internal/movie_screenshots"
	"ozinshe/internal/seasons"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	conn, err := database.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close(ctx)

	genreRepo := genres.NewGenreRepository(conn)
	genreService := genres.NewGenreService(genreRepo)
	genreHandler := genres.NewGenreHandler(genreService)

	movieGenresRepo := movie_genres.NewMovieGenreRepository(conn)
	movieGenresService := movie_genres.NewMovieGenreService(movieGenresRepo)
	movieGenresHandler := movie_genres.NewMovieGenreHandler(movieGenresService)

	categoriesRepo := categories.NewCategoryRepository(conn)
	categoriesService := categories.NewCategoryService(categoriesRepo)
	categoriesHandler := categories.NewCategoryHandler(categoriesService)

	movieCategoriesRepo := movie_categories.NewMovieCategoryRepositry(conn)
	movieCategoriesService := movie_categories.NewMovieCategoryService(movieCategoriesRepo)
	movieCategoriesHandler := movie_categories.NewMovieCategoryHandler(movieCategoriesService)

	episodeRepo := episodes.NewEpisodeRepository(conn)
	episodeService := episodes.NewEpisodeService(episodeRepo)
	episodeHandler := episodes.NewEpisodeHandler(episodeService)

	seasonRepo := seasons.NewSeasonRepository(conn)
	seasonService := seasons.NewSeasonService(seasonRepo)
	seasonHandler := seasons.NewSeasonHandler(seasonService)

	screensotRepo := movie_screenshots.NewMovieScreenshotsRepository(conn)
	screensotService := movie_screenshots.NewMovieScreenshotsService(screensotRepo)
	screensotHandler := movie_screenshots.NewMovieScreenshotsHandler(screensotService)

	movieRepo := movie.NewMovieRepository(conn)
	movieService := movie.NewMovieService(
		movieRepo,
		movieGenresService,
		movieCategoriesService,
		seasonService,
		episodeService,
		screensotService,
	)
	movieHandler := movie.NewMovieHandler(movieService)

	http.HandleFunc("GET /movies", movieHandler.GetAll)
	http.HandleFunc("GET /movies/{id}", movieHandler.GetByID)
	http.HandleFunc("POST /movies", movieHandler.CreateMovie)
	http.HandleFunc("DELETE /movies/{id}", movieHandler.DeleteMovie)
	http.HandleFunc("PUT /movies/{id}", movieHandler.UpdateMovie)

	http.HandleFunc("GET /genres", genreHandler.GetAll)
	http.HandleFunc("GET /genres/{id}", genreHandler.GetGenreByID)
	http.HandleFunc("POST /genres", genreHandler.CreateGenre)
	http.HandleFunc("DELETE /genres/{id}", genreHandler.DeleteGenre)
	http.HandleFunc("PATCH /genres/{id}", genreHandler.UpdateGenre)

	http.HandleFunc("POST /movies/{id}/genres", movieGenresHandler.AddGenreToMovie)
	http.HandleFunc("GET /movies/{id}/genres", movieGenresHandler.GetGenresOfMovie)
	http.HandleFunc("DELETE /movies/{id}/genres", movieGenresHandler.DeleteGenreFromMovie)

	http.HandleFunc("GET /categories", categoriesHandler.GetCategories)
	http.HandleFunc("GET /categories/{id}", categoriesHandler.GetCategoryByID)
	http.HandleFunc("POST /categories", categoriesHandler.CreateCategory)
	http.HandleFunc("DELETE /categories/{id}", categoriesHandler.DeleteCategory)
	http.HandleFunc("PATCH /categories/{id}", categoriesHandler.UpdateCategory)

	http.HandleFunc("POST /movies/{id}/categories", movieCategoriesHandler.AddCategoryToMovie)
	http.HandleFunc("GET /movies/{id}/categories", movieCategoriesHandler.GetMovieCategories)
	http.HandleFunc("DELETE /movies/{id}/categories", movieCategoriesHandler.DeleteCategoryFromMovie)

	http.HandleFunc("POST /movies/{movieID}/seasons", seasonHandler.AddSeasonToMovie)
	http.HandleFunc("DELETE /seasons/{id}", seasonHandler.DeleteSeasonFromMovie)

	http.HandleFunc("POST /seasons/{id}/episodes", episodeHandler.AddEpisodeToSeason)
	http.HandleFunc("DELETE /episodes/{id}", episodeHandler.DeleteEpisode)

	http.HandleFunc("POST /movies/{id}/screenshots", screensotHandler.AddScreenshotsToMovie)
	http.HandleFunc("GET /movies/{id}/screenshots", screensotHandler.GetScreenshots)
	http.HandleFunc("DELETE /screenshots/{id}", screensotHandler.DeleteScreenshot)

	log.Println("Сервер слушает на порту :8080...")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
