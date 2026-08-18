package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=postgres dbname=shorturl sslmode=disable"
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Connection pool settings (Postgres handles concurrency well)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	createTables()
}

func createTables() {
	createURLsTable := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		url TEXT NOT NULL,
		short_code TEXT UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	createStatsTable := `
	CREATE TABLE IF NOT EXISTS stats (
		id SERIAL PRIMARY KEY,
		short_code TEXT NOT NULL,
		accessed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (short_code) REFERENCES urls(short_code) ON DELETE CASCADE
	);`

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := db.Exec(createURLsTable); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(createStatsTable); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(createUsersTable); err != nil {
		log.Fatal(err)
	}
}

func closeDB() {
	db.Close()
}

func createUser(username, passwordHash string) (int, error) {
	var userID int
	err := db.QueryRow(
		"INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id",
		username, passwordHash,
	).Scan(&userID)

	return userID, err
}

func getUserByUsername(username string) (int, string, error) {
	var userID int
	var passwordHash string
	err := db.QueryRow(
		"SELECT id, password_hash FROM users WHERE username = $1",
		username,
	).Scan(&userID, &passwordHash)

	return userID, passwordHash, err
}

func createURL(originalURL string) (*URL, error) {
	shortCode := generateShortCode()

	for {
		var exists int
		err := db.QueryRow("SELECT COUNT(*) FROM urls WHERE short_code = $1", shortCode).Scan(&exists)
		if err != nil {
			return nil, err
		}
		if exists == 0 {
			break
		}
		shortCode = generateShortCode()
	}

	var url URL
	err := db.QueryRow(
		`INSERT INTO urls (url, short_code) 
		 VALUES ($1, $2) 
		 RETURNING id, url, short_code, created_at, updated_at`,
		originalURL, shortCode,
	).Scan(&url.ID, &url.URL, &url.ShortCode, &url.CreatedAt, &url.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func getURLByShortCode(shortCode string) (*URL, error) {
	var url URL
	err := db.QueryRow(
		"SELECT id, url, short_code, created_at, updated_at FROM urls WHERE short_code = $1",
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

func updateURL(shortCode string, newURL string) (*URL, error) {
	var url URL
	err := db.QueryRow(
		`UPDATE urls SET url = $1, updated_at = CURRENT_TIMESTAMP WHERE short_code = $2
		RETURNING id, url, short_code, created_at, updated_at`,
		newURL, shortCode,
	).Scan(&url.ID, &url.URL, &url.ShortCode, &url.CreatedAt, &url.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}

	return &url, nil
}

func deleteURL(shortCode string) (bool, error) {
	result, err := db.Exec("DELETE FROM urls WHERE short_code = $1", shortCode)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, nil // Not found
	}

	return true, nil // Successfully deleted
}

func recordAccess(shortCode string) error {
	_, err := db.Exec(
		"INSERT INTO stats (short_code) VALUES ($1)",
		shortCode,
	)
	return err
}

func getURLStats(shortCode string) (*URLStats, error) {
	var stats URLStats

	// Get URL data and count stats
	err := db.QueryRow(`
		SELECT 
			u.id, 
			u.url, 
			u.short_code, 
			u.created_at, 
			u.updated_at,
			COUNT(s.id) as access_count
		FROM urls u
		LEFT JOIN stats s ON u.short_code = s.short_code
		WHERE u.short_code = $1
		GROUP BY u.id
	`, shortCode).Scan(
		&stats.ID,
		&stats.URL,
		&stats.ShortCode,
		&stats.CreatedAt,
		&stats.UpdatedAt,
		&stats.AccessCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
