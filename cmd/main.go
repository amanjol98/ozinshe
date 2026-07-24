package main

import (
	"context"
	"log"
	"net/http"
	"ozinshe/internal/database"
	"ozinshe/internal/genres"
	"ozinshe/internal/movie"
	"ozinshe/internal/movie_genres"

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

	movieRepo := movie.NewMovieRepository(conn)

	movieService := movie.NewMovieService(movieRepo)

	movieHandler := movie.NewMovieHandler(movieService)

	http.HandleFunc("GET /movies", movieHandler.GetAll)
	http.HandleFunc("GET /movies/{id}", movieHandler.GetByID)
	http.HandleFunc("POST /movies", movieHandler.CreateMovie)
	http.HandleFunc("DELETE /movies/{id}", movieHandler.DeleteMovie)
	http.HandleFunc("PUT /movies/{id}", movieHandler.UpdateMovie)

	genreRepo := genres.NewGenreRepository(conn)
	genreService := genres.NewGenreService(genreRepo)
	genreHandler := genres.NewGenreHandler(genreService)

	http.HandleFunc("GET /genres", genreHandler.GetAll)
	http.HandleFunc("GET /genres/{id}", genreHandler.GetGenreByID)
	http.HandleFunc("POST /genres", genreHandler.CreateGenre)
	http.HandleFunc("DELETE /genres/{id}", genreHandler.DeleteGenre)
	http.HandleFunc("PATCH /genres/{id}", genreHandler.UpdateGenre)

	movieGenresRepo := movie_genres.NewMovieGenreRepository(conn)
	movieGenresService := movie_genres.NewMovieGenreService(movieGenresRepo)
	movieGenresHandler := movie_genres.NewMovieGenreHandler(movieGenresService)

	http.HandleFunc("POST /movies/{id}/genres", movieGenresHandler.AddGenreToMovie)
	http.HandleFunc("GET /movies/{id}/genres", movieGenresHandler.GetGenresOfMovie)
	http.HandleFunc("DELETE /movies/{id}/genres", movieGenresHandler.DeleteGenreFromMovie)

	log.Println("Сервер слушает на порту :8080...")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
