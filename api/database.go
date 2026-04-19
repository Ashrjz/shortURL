package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "../urls.db?_busy_timeout=5000&_journal_mode=WAL")
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS urls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func closeDB() {
	db.Close()
}

func createURL(originalURL string) (*URL, error) {
	shortCode := generateShortCode()

	// Check if short code already exists, regenerate if needed
	for {
		var exists int
		err := db.QueryRow("SELECT COUNT(*) FROM urls WHERE short_code = ?", shortCode).Scan(&exists)
		if err != nil {
			return nil, err
		}
		if exists == 0 {
			break
		}
		shortCode = generateShortCode()
	}

	result, err := db.Exec(
		"INSERT INTO urls (url, short_code) VALUES (?, ?)",
		originalURL, shortCode,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// Fetch the created record with timestamps
	var url URL
	err = db.QueryRow(
		"SELECT id, url, short_code, created_at, updated_at FROM urls WHERE id = ?",
		id,
	).Scan(&url.ID, &url.URL, &url.ShortCode, &url.CreatedAt, &url.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func getURLByShortCode(shortCode string) (*URL, error) {
	var url URL
	err := db.QueryRow(
		"SELECT id, url, short_code, created_at, updated_at FROM urls WHERE short_code = ?",
		shortCode,
	).Scan(&url.ID, &url.URL, &url.ShortCode, &url.CreatedAt, &url.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err // Other error
	}

	return &url, nil
}
