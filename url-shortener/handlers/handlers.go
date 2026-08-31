package handler

import (
	"encoding/json"
	"net/http"
	"url-shortener/repository"

	"github.com/jackc/pgx/v5"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type Response struct {
	Message string         `json:"message,omitempty"`
	URL     repository.URL `json:"url,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ShortenHandler(conn *pgx.Conn, clicks chan repository.Click) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var s ShortenRequest

		var err error
		var result *repository.URL
		if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
			writeJSON(w, http.StatusBadRequest, &Response{Message: "Url is required"})
			return
		}

		if s.URL == "" {
			writeJSON(w, http.StatusBadRequest, &Response{Message: "Url cannot be empty"})
			return
		}

		code := repository.GenerateShortCode()

		if result, err = repository.CreateUrl(conn, code, s.URL); err != nil {
			writeJSON(w, http.StatusInternalServerError, &Response{Message: "Error creating Url"})
			return
		}

		writeJSON(w, http.StatusCreated, &Response{Message: "Key set successfully", URL: *result})

	}
}
