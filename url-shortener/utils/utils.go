package utils

import (
	"encoding/json"
	"net/http"
	"url-shortener/repository"
)

type Response struct {
	Message   string                     `json:"message,omitempty"`
	URL       *repository.URL            `json:"url,omitempty"`
	Clicks    []repository.Click         `json:"clicks,omitempty"`
	Timelines []repository.TimelineEntry `json:"timelines,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
