package app

import (
	"context"
	"net/http"
	"ozinshe/internal/database"
	"ozinshe/internal/handlers"
	"ozinshe/internal/middleware"
	"ozinshe/internal/repositories"
	"ozinshe/internal/routes"
	"ozinshe/internal/services"
)

func New(ctx context.Context, secretKey string) (*http.ServeMux, func(), error) {
	conn, err := database.Connect(ctx)
	if err != nil {
		return nil, nil, err
	}

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

	movieScreenshotsRepo := repositories.NewMovieScreenshotsRepository(conn)
	movieScreenshotsService := services.NewMovieScreenshotsService(movieScreenshotsRepo)
	movieScreenshotsHandler := handlers.NewMovieScreenshotsHandler(movieScreenshotsService)

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
		movieScreenshotsService,
	)
	movieHandler := handlers.NewMovieHandler(movieService)

	mux := http.NewServeMux()

	routes.Register(
		mux,
		categoriesHandler,
		episodeHandler,
		genreHandler,
		movieCategoriesHandler,
		movieGenresHandler,
		movieScreenshotsHandler,
		movieHandler,
		seasonHandler,
		userHandler,
		authMiddleware,
	)

	cleanup := func() {
		conn.Close(ctx)
	}

	return mux, cleanup, nil
}
