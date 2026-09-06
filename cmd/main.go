package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"ozinshe/internal/app"

	_ "ozinshe/docs"

	"github.com/joho/godotenv"
)

// @title Özinshe API
// @version 1.0
// @description REST API для приложения Özinshe.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	secretKey := os.Getenv("JWT_SECRET")

	mux, cleanup, err := app.New(ctx, secretKey)
	if err != nil {
		log.Fatal(err)
	}

	defer cleanup()

	log.Println("Сервер слушает на порту :8080...")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
