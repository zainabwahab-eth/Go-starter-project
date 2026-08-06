package main

import (
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

func NewStore() *Store {
	return &Store{
		items: make(map[string]Item),
	}
}

func setHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hitset")
	}
}

func getHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hitget")
	}
}

func deleteHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hitdelete")
	}
}

type SetRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

func main() {
	store := NewStore()

	go store.StartEviction()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /set", setHandler(store))
	mux.HandleFunc("GET /get", getHandler(store))
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
