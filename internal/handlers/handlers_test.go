package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shortener/internal/generator"
	"shortener/internal/repository"
	"shortener/internal/service"
	"testing"

	"github.com/gorilla/mux"
)

func newTestHandler() *Handlers {
	repo := repository.NewMemoryRepository()
	gen := generator.NewMemoryGenerator()
	svc := service.NewService(repo, gen)

	return NewHandlers(svc)
}

func TestHandler_Shorten(t *testing.T) {
	h := newTestHandler()

	body := []byte(`{"url":"https://finance.ozon.ru"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	h.Shorten(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	if rec.Body.Len() == 0 {
		t.Fatal("expected non-empty response")
	}

	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp["shortUrl"] == "" {
		t.Fatal("expected shortened url")
	}
}

func TestHandler_Redirect(t *testing.T) {
	h := newTestHandler()

	body := []byte(`{"url":"https://finance.ozon.ru"}`)
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()
	h.Shorten(rec, req)

	short := "aaaaaaaaab"

	req2 := httptest.NewRequest(http.MethodGet, "/"+short, nil)
	req2 = mux.SetURLVars(req2, map[string]string{
		"short": short,
	})

	rec2 := httptest.NewRecorder()
	h.Redirect(rec2, req2)

	if rec2.Code != http.StatusFound {
		t.Fatalf("expected status 302, got %d", rec2.Code)
	}
}
