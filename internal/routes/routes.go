package routes

import (
	"net/http"
	"ozinshe/internal/handlers"
	"ozinshe/internal/middleware"
)

func Register(
	mux *http.ServeMux,
	categoriesHandler *handlers.CategoryHandler,
	episodeHandler *handlers.EpisodeHandler,
	favoriteHandler *handlers.FavoriteHandler,
	genreHandler *handlers.GenreHandler,
	movieCategoriesHandler *handlers.MovieCategoryHandler,
	movieGenresHandler *handlers.MovieGenreHandler,
	movieScreenshotsHandler *handlers.MovieScreenshotsHandler,
	movieHandler *handlers.MovieHandler,
	seasonsHandler *handlers.SeasonHandler,
	userHandler *handlers.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) {

	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)

	mux.HandleFunc("GET /movies", movieHandler.GetAll)
	mux.HandleFunc("GET /movies/{id}", movieHandler.GetByID)
	mux.HandleFunc("GET /genres", genreHandler.GetAll)
	mux.HandleFunc("GET /genres/{id}", genreHandler.GetGenreByID)
	mux.HandleFunc("GET /movies/{id}/genres", movieGenresHandler.GetGenresOfMovie)
	mux.HandleFunc("GET /categories", categoriesHandler.GetCategories)
	mux.HandleFunc("GET /categories/{id}", categoriesHandler.GetCategoryByID)
	mux.HandleFunc("GET /movies/{id}/categories", movieCategoriesHandler.GetMovieCategories)
	mux.HandleFunc("GET /movies/{id}/screenshots", movieScreenshotsHandler.GetScreenshots)
	mux.HandleFunc("GET /home", movieHandler.GetHome)

	mux.Handle("GET /users/me", authMiddleware.Auth(http.HandlerFunc(userHandler.Me)))
	mux.Handle("PATCH /users/me", authMiddleware.Auth(http.HandlerFunc(userHandler.Update)))
	mux.Handle("PATCH /users/me/password", authMiddleware.Auth(http.HandlerFunc(userHandler.UpdatePassword)))

	mux.Handle(
		"POST /movies",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(movieHandler.CreateMovie),
			),
		),
	)

	mux.Handle(
		"PUT /movies/{id}",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(movieHandler.UpdateMovie),
			),
		),
	)

	mux.Handle(
		"DELETE /movies/{id}",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(movieHandler.DeleteMovie),
			),
		),
	)

	mux.Handle(
		"POST /genres",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(genreHandler.CreateGenre),
			),
		),
	)

	mux.Handle(
		"DELETE /genres/{id}",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(genreHandler.DeleteGenre),
			),
		),
	)

	mux.Handle(
		"PATCH /genres/{id}",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(genreHandler.UpdateGenre),
			),
		),
	)

	mux.Handle(
		"POST /movies/{id}/genres",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(movieGenresHandler.AddGenreToMovie),
			),
		),
	)

	mux.Handle(
		"DELETE /movies/{id}/genres",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(movieGenresHandler.DeleteGenreFromMovie),
			),
		),
	)

	mux.Handle(
		"POST /categories",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(categoriesHandler.CreateCategory),
			),
		),
	)

	mux.Handle(
		"DELETE /categories/{id}",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(categoriesHandler.DeleteCategory),
			),
		),
	)

	mux.Handle(
		"PATCH /categories/{id}",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(categoriesHandler.UpdateCategory),
			),
		),
	)

	mux.Handle(
		"POST /movies/{id}/categories",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(movieCategoriesHandler.AddCategoryToMovie),
			),
		),
	)

	mux.Handle(
		"DELETE /movies/{id}/categories",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(movieCategoriesHandler.DeleteCategoryFromMovie),
			),
		),
	)

	mux.Handle(
		"POST /movies/{movieID}/seasons",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(seasonsHandler.AddSeasonToMovie),
			),
		),
	)

	mux.Handle(
		"DELETE /seasons/{id}",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(seasonsHandler.DeleteSeasonFromMovie),
			),
		),
	)

	mux.Handle(
		"POST /seasons/{id}/episodes",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(episodeHandler.AddEpisodeToSeason),
			),
		),
	)

	mux.Handle(
		"DELETE /episodes/{id}",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(episodeHandler.DeleteEpisode),
			),
		),
	)

	mux.Handle(
		"POST /movies/{id}/screenshots",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(movieScreenshotsHandler.AddScreenshotsToMovie),
			),
		),
	)

	mux.Handle(
		"DELETE /screenshots/{id}",
		authMiddleware.Auth(
			authMiddleware.RequireRole(
				"admin",
				http.HandlerFunc(movieScreenshotsHandler.DeleteScreenshot),
			),
		),
	)

	mux.Handle(
		"GET /favorites",
		authMiddleware.Auth(
			http.HandlerFunc(favoriteHandler.GetFavoriteMovies),
		),
	)

	mux.Handle(
		"POST /favorites",
		authMiddleware.Auth(
			http.HandlerFunc(favoriteHandler.AddFavoriteMovieToUser),
		),
	)

	mux.Handle(
		"DELETE /favorites/{movie_id}",
		authMiddleware.Auth(
			http.HandlerFunc(favoriteHandler.DeleteFavoriteMovieFromUser),
		),
	)
}
