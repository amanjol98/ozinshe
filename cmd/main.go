package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"ozinshe/internal/database"
	"ozinshe/internal/handlers"
	"ozinshe/internal/middleware"
	"ozinshe/internal/repositories"
	"ozinshe/internal/services"

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

	secretKey := os.Getenv("JWT_SECRET")

	defer conn.Close(ctx)

	genreRepo := repositories.NewGenreRepository(conn)
	genreService := services.NewGenreService(genreRepo)
	genreHandler := handlers.NewGenreHandler(genreService)

	movieGenresRepo := repositories.NewMovieGenreRepository(conn)
	movieGenresService := services.NewMovieGenreService(movieGenresRepo)
	movieGenresHandler := handlers.NewMovieGenreHandler(movieGenresService)

	categoriesRepo := repositories.NewCategoryRepository(conn)
	categoriesService := services.NewCategoryService(categoriesRepo)
	categoriesHandler := handlers.NewCategoryHandler(categoriesService)

	movieCategoriesRepo := repositories.NewMovieCategoryRepositry(conn)
	movieCategoriesService := services.NewMovieCategoryService(movieCategoriesRepo)
	movieCategoriesHandler := handlers.NewMovieCategoryHandler(movieCategoriesService)

	episodeRepo := repositories.NewEpisodeRepository(conn)
	episodeService := services.NewEpisodeService(episodeRepo)
	episodeHandler := handlers.NewEpisodeHandler(episodeService)

	seasonRepo := repositories.NewSeasonRepository(conn)
	seasonService := services.NewSeasonService(seasonRepo)
	seasonHandler := handlers.NewSeasonHandler(seasonService)

	screensotRepo := repositories.NewMovieScreenshotsRepository(conn)
	screensotService := services.NewMovieScreenshotsService(screensotRepo)
	screensotHandler := handlers.NewMovieScreenshotsHandler(screensotService)

	userRepo := repositories.NewUserRepository(conn)
	tokenService := middleware.NewTokenService(secretKey)
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	userService := services.NewUserService(userRepo, tokenService)
	userHandler := handlers.NewUserHandler(userService)

	movieRepo := repositories.NewMovieRepository(conn)
	movieService := services.NewMovieService(
		movieRepo,
		movieGenresService,
		movieCategoriesService,
		seasonService,
		episodeService,
		screensotService,
	)
	movieHandler := handlers.NewMovieHandler(movieService)

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

	http.HandleFunc("POST /register", userHandler.Register)
	http.HandleFunc("POST /login", userHandler.Login)

	http.Handle("GET /users/me", authMiddleware.Auth(http.HandlerFunc(userHandler.Me)))

	log.Println("Сервер слушает на порту :8080...")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
