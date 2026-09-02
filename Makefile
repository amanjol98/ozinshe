run:
	go run ./cmd

migrate-up:
	migrate -path ./migrations -database "postgres://postgres:1111@localhost:5432/ozinshe?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:1111@localhost:5432/ozinshe?sslmode=disable" down