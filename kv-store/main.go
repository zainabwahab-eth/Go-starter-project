package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Item struct {
	Value     string
	ExpiresAt time.Time
}

type Store struct {
	mu    sync.RWMutex
	items map[string]Item
}

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

type Response struct {
	Message string   `json:"message,omitempty"`
	Value   string   `json:"value,omitempty"`
	Keys     []string `json:"keys,omitempty"`
}

func (s *Store) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var it = Item{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}

	s.items[key] = it
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[key]

	if !ok {
		return "", false
	}

	if time.Now().After(item.ExpiresAt) {
		return "", false
	}

	return item.Value, true
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
}

func (s *Store) StartEviction() {
	for {
		time.Sleep(time.Second)
		s.mu.Lock()
		for k, v := range s.items {
			if time.Now().After(v.ExpiresAt) {
				delete(s.items, k)
			}
		}
		s.mu.Unlock()
	}
}

func (s *Store) GetAll() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []string

	for k, v := range s.items {
		if !time.Now().After(v.ExpiresAt) {
			keys = append(keys, k)
		}
	}
	return keys
}

func NewStore() *Store {
	return &Store{
		items: make(map[string]Item),
	}
}

func writeResponse(status int, res *Response, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func setHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m SetRequest

		if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
			writeResponse(http.StatusBadRequest, &Response{Message: "Error Decoding request. Please check request"}, w)
			return
		}

		if m.Key == "" || m.Value == "" || m.TTL <= 0 {
			writeResponse(http.StatusBadRequest, &Response{Message: "Key or value cannot be empty and ttl must be greater than zero"}, w)
			return
		}

		store.Set(m.Key, m.Value, time.Duration(m.TTL)*time.Second)
		writeResponse(http.StatusCreated, &Response{Message: "Key set successfully"}, w)
	}
}

func getHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var response *Response
		key := r.URL.Query().Get("key")

		if key == "" {
			writeResponse(http.StatusBadRequest, &Response{Message: "Key cannot be empty"}, w)
			return
		}

		val, ok := store.Get(key)
		if ok == false {
			writeResponse(http.StatusNotFound, &Response{Message: "Key not found or expired"}, w)
			return
		}

		writeResponse(http.StatusOK, &Response{Value: val}, w)
	}
}

func deleteHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")

		if key == "" {
			writeResponse(http.StatusBadRequest, &Response{Message: "Key cannot be empty"}, w)
			return
		}

		_, ok := store.Get(key)
		if ok == false {
			writeResponse(http.StatusNotFound, &Response{Message: "Key not found or expired"}, w)
			return
		}

		store.Delete(key)

		writeResponse(http.StatusOK, &Response{Message: "Key deleted successfully"}, w)

	}
}

func getAllKeysHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys := store.GetAll()
		writeResponse(http.StatusOK, &Response{Keys: keys}, w)
	}

}

func main() {
	store := NewStore()

	go store.StartEviction()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /set", setHandler(store))
	mux.HandleFunc("GET /get", getHandler(store))
	mux.HandleFunc("GET /keys", getAllKeysHandler(store))
	mux.HandleFunc("DELETE /delete", deleteHandler(store))

	server := &http.Server{
		Handler: mux,
		Addr:    ":8181",
	}

	fmt.Println("Server starting at 8181")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
