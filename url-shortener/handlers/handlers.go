package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"url-shortener/repository"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/mssola/useragent"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type Response struct {
	Message   string                     `json:"message,omitempty"`
	URL       *repository.URL            `json:"url,omitempty"`
	Clicks    []repository.Click         `json:"clicks,omitempty"`
	Timelines []repository.TimelineEntry `json:"timelines,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func ShortenHandler(conn *pgx.Conn) http.HandlerFunc {
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

		writeJSON(w, http.StatusCreated, &Response{Message: "Key set successfully", URL: result})

	}
}

func RedirectHandler(conn *pgx.Conn, clicks chan repository.Click) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := chi.URLParam(r, "code")

		url, err := repository.GetUrl(conn, shortCode)

		if err != nil {
			writeJSON(w, http.StatusNotFound, &Response{Message: "Url not found"})
			return
		}

		ua := useragent.New(r.UserAgent())
		browser, _ := ua.Browser()
		device := "Desktop"
		if ua.Mobile() {
			device = "Mobile"
		}

		click := repository.Click{
			ShortCode: shortCode,
			ClickedAt: time.Now(),
			Referrer:  r.Referer(),
			Browser:   browser,
			OS:        ua.OS(),
			Device:    device,
		}

		clicks <- click

		http.Redirect(w, r, url.OriginalURL, http.StatusFound)
	}
}

func StatsHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := chi.URLParam(r, "code")

		url, err := repository.GetUrl(conn, shortCode)

		if err != nil {
			writeJSON(w, http.StatusNotFound, &Response{Message: "Url not found"})
			return
		}

		clicks, err := repository.GetClickStats(conn, shortCode)

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, &Response{Message: "Something went wrong"})
			return
		}

		writeJSON(w, http.StatusOK, &Response{URL: url, Clicks: clicks})
	}
}

func TimelineHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := chi.URLParam(r, "code")

		url, err := repository.GetUrl(conn, shortCode)

		if err != nil {
			writeJSON(w, http.StatusNotFound, &Response{Message: "Url not found"})
			return
		}

		clicksTL, err := repository.GetClickTimeline(conn, shortCode)

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, &Response{Message: "Something went wrong"})
			return
		}

		writeJSON(w, http.StatusOK, &Response{URL: url, Timelines: clicksTL})
	}
}
