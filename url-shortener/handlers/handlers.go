package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"url-shortener/repository"
	"url-shortener/utils"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/mssola/useragent"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

func ShortenHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var s ShortenRequest

		var err error
		var result *repository.URL
		if err = json.NewDecoder(r.Body).Decode(&s); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, &utils.Response{Message: "Url is required"})
			return
		}

		if s.URL == "" {
			utils.WriteJSON(w, http.StatusBadRequest, &utils.Response{Message: "Url cannot be empty"})
			return
		}

		code := repository.GenerateShortCode()

		if result, err = repository.CreateUrl(conn, code, s.URL); err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, &utils.Response{Message: "Error creating Url"})
			return
		}

		utils.WriteJSON(w, http.StatusCreated, &utils.Response{Message: "Key set successfully", URL: result})

	}
}

func RedirectHandler(conn *pgx.Conn, clicks chan repository.Click) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := chi.URLParam(r, "code")

		url, err := repository.GetUrl(conn, shortCode)

		if err != nil {
			utils.WriteJSON(w, http.StatusNotFound, &utils.Response{Message: "Url not found"})
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
			utils.WriteJSON(w, http.StatusNotFound, &utils.Response{Message: "Url not found"})
			return
		}

		clicks, err := repository.GetClickStats(conn, shortCode)

		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, &utils.Response{Message: "Something went wrong"})
			return
		}

		utils.WriteJSON(w, http.StatusOK, &utils.Response{URL: url, Clicks: clicks})
	}
}

func TimelineHandler(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortCode := chi.URLParam(r, "code")

		url, err := repository.GetUrl(conn, shortCode)

		if err != nil {
			utils.WriteJSON(w, http.StatusNotFound, &utils.Response{Message: "Url not found"})
			return
		}

		clicksTL, err := repository.GetClickTimeline(conn, shortCode)

		if err != nil {
			utils.WriteJSON(w, http.StatusInternalServerError, &utils.Response{Message: "Something went wrong"})
			return
		}

		utils.WriteJSON(w, http.StatusOK, &utils.Response{URL: url, Timelines: clicksTL})
	}
}
