package main

import (
	"log"
	"net/http"
	"shortener/internal/generator"
	"shortener/internal/handlers"
	"shortener/internal/repository"
	"shortener/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	repo := repository.NewMemoryRepository()
	gen := generator.NewMemoryGenerator()

	srv := service.NewService(repo, gen)

	h := handlers.NewHandlers(srv)
	r := mux.NewRouter()

	r.HandleFunc("/shorten", h.Shorten).Methods(http.MethodPost)
	r.HandleFunc("/{short}", h.Redirect).Methods(http.MethodGet)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
