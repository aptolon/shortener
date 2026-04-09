package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"shortener/internal/generator"
	"shortener/internal/handlers"
	"shortener/internal/repository"
	"shortener/internal/service"
)

func main() {
	ctx := context.Background()

	storage := os.Getenv("STORAGE")
	if storage == "" {
		storage = "memory"
	}

	var (
		repo repository.Repository
		gen  generator.Generator
	)

	switch storage {
	case "memory":
		repo = repository.NewMemoryRepository()
		gen = generator.NewMemoryGenerator()

	case "postgres":
		databaseURL := os.Getenv("DATABASE_URL")
		if databaseURL == "" {
			log.Fatal("DATABASE_URL is required for postgres storage")
		}

		db, err := pgxpool.New(ctx, databaseURL)
		if err != nil {
			log.Fatalf("failed to connect postgres: %v", err)
		}
		defer db.Close()

		if err := db.Ping(ctx); err != nil {
			log.Fatalf("failed to ping postgres: %v", err)
		}

		repo = repository.NewPostgresRepository(db)
		gen = generator.NewPostgresGenerator(db)

	default:
		log.Fatalf("unknown storage: %s", storage)
	}

	srv := service.NewService(repo, gen)
	h := handlers.NewHandlers(srv)

	r := mux.NewRouter()
	r.HandleFunc("/shorten", h.Shorten).Methods(http.MethodPost)
	r.HandleFunc("/{short}", h.Redirect).Methods(http.MethodGet)

	port := os.Getenv("SERV_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server started on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
