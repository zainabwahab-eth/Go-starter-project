package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type URL struct {
	ID          int       `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type Click struct {
	ID        int       `json:"id"`
	ShortCode string    `json:"short_code"`
	ClickedAt time.Time `json:"clicked_at"`
	Device    string    `json:"device"`
	Browser   string    `json:"browser"`
	OS        string    `json:"os"`
	Referrer  string    `json:"referrer"`
}

func CreateUrl(conn *pgx.Conn, shortCode, originalURL string) (*URL, error) {
	query := `
			INSERT INTO urls (short_code, original_url)
			VALUES ($1, $2)
			RETURNING id, short_code, original_url, created_at
	`

	var url URL

	err := conn.QueryRow(context.Background(), query, shortCode, originalURL).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &url, nil
}

func GetUrl(conn *pgx.Conn, shortCode string) (*URL, error) {
	query := `
		SELECT id, short_code, original_url, created_at FROM urls WHERE short_code = $1
	`

	var url URL
	err := conn.QueryRow(context.Background(), query, shortCode).Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &url, nil
}

func SaveClick(conn *pgx.Conn, click Click) error {
	query := `
			INSERT INTO clicks (short_code, device, browser, os, referrer)
			VALUES ($1, $2, $3, $4, $5)
	`

	tag, err := conn.Exec(context.Background(), query, click.ShortCode, click.Device, click.Browser, click.OS, click.Referrer)

	if err != nil {
		return fmt.Errorf("failed to save click: %w", err)
	}

	rowsInserted := tag.RowsAffected()

	fmt.Printf("Successfully inserted %d row(s)\n", rowsInserted)

	return nil
}

func GetClickStats(conn *pgx.Conn, shortCode string) ([]Click, error) {
	query := `
			SELECT id, short_code, clicked_at, device, browser, os, referrer FROM clicks WHERE short_code = $1
	`

	rows, err := conn.Query(context.Background(), query, shortCode)

	var clicks = []Click{}
	var click Click

	if err != nil {
		return clicks, fmt.Errorf("failed to get click stat: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&click.ID, &click.ShortCode, &click.ClickedAt, &click.Device, &click.Browser, &click.OS, &click.Referrer)

		if err != nil {
			fmt.Errorf("failed to get click stat: %w", err)
		}

		clicks = append(clicks, click)
	}

	if err := rows.Err(); err != nil {
		return clicks, fmt.Errorf("failed to get click stat: %w", err)
	}

	return clicks, nil

}
