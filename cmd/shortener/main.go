package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"shortener/internal/app"
	"shortener/internal/handlers"
	"shortener/internal/service"
)

func main() {
	ctx := context.Background()

	repo, gen, close, err := app.BuildStorage(ctx, os.Getenv("STORAGE"), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer close()

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
	log.Fatal(http.ListenAndServe(":8080", r))
}
