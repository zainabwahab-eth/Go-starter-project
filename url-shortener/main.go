package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	handler "url-shortener/handlers"
	"url-shortener/repository"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func connectDB(connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	return conn, nil
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.URL.Path)
}

func main() {

	godotenv.Load()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL not set")
	}

	conn, err := connectDB(connStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DB successfully connected")
	defer conn.Close(context.Background())

	clicks := make(chan repository.Click, 100)

	go func() {
		for click := range clicks {
			repository.SaveClick(conn, click)
		}
	}()

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Post("/shorten", handler.ShortenHandler(conn, clicks))
	r.Get("/{code}", testHandler)
	r.Get("/stats/{code}", testHandler)
	r.Get("/stats/{code}/timeline", testHandler)

	server := &http.Server{
		Handler: r,
		Addr:    ":8181",
	}

	fmt.Println("Server starting at 8181")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
