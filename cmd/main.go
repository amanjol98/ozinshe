package main

import (
	"context"
	"log"
	"net/http"
	"ozinshe/internal/database"
	"ozinshe/internal/genres"
	"ozinshe/internal/movie"

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

	repo := movie.NewMovieRepository(conn)

	service := movie.NewMovieService(repo)

	handler := movie.NewMovieHandler(service)

	http.HandleFunc("GET /movies", handler.GetAll)
	http.HandleFunc("GET /movies/{id}", handler.GetByID)
	http.HandleFunc("POST /movies", handler.CreateMovie)
	http.HandleFunc("DELETE /movies/{id}", handler.DeleteMovie)
	http.HandleFunc("PUT /movies/{id}", handler.UpdateMovie)

	genreRepo := genres.NewGenreRepository(conn)
	genreService := genres.NewGenreService(genreRepo)
	genreHandler := genres.NewGenreHandler(genreService)

	http.HandleFunc("GET /genres", genreHandler.GetAll)
	http.HandleFunc("POST /genres", genreHandler.CreateGenre)

	log.Println("Сервер слушает на порту :8080...")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
