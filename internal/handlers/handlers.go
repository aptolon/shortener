package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"shortener/internal/errs"
	"shortener/internal/service"

	"github.com/gorilla/mux"
)

type Handlers struct {
	service *service.Service
}

func NewHandlers(service *service.Service) *Handlers {
	return &Handlers{
		service: service,
	}
}

type shortenResponse struct {
	ShortUrl string `json:"shortUrl"`
}

type shortenRequest struct {
	Url string `json:"url"`
}

func (h *Handlers) Shorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if req.Url == "" || err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	short, err := h.service.Shorten(r.Context(), req.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := shortenResponse{
		ShortUrl: "http://localhost:8080/" + short,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) Redirect(w http.ResponseWriter, r *http.Request) {
	short := mux.Vars(r)["short"]
	if short == "" {
		http.Error(w, "short code is required", http.StatusBadRequest)
		return
	}

	original, err := h.service.GetOriginal(r.Context(), short)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, original, http.StatusFound)
}
