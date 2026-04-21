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
	);
	CREATE TABLE IF NOT EXISTS stats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_code TEXT NOT NULL,
		accessed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (short_code) REFERENCES urls(short_code) ON DELETE CASCADE
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

func updateURL(shortCode string, newURL string) (*URL, error) {
	// Check if short code exists
	var exists int
	err := db.QueryRow("SELECT COUNT(*) FROM urls WHERE short_code = ?", shortCode).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, nil // Not found
	}

	// Update the URL
	_, err = db.Exec(
		"UPDATE urls SET url = ?, updated_at = CURRENT_TIMESTAMP WHERE short_code = ?",
		newURL, shortCode,
	)
	if err != nil {
		return nil, err
	}

	// Fetch updated record
	var url URL
	err = db.QueryRow(
		"SELECT id, url, short_code, created_at, updated_at FROM urls WHERE short_code = ?",
		shortCode,
	).Scan(&url.ID, &url.URL, &url.ShortCode, &url.CreatedAt, &url.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &url, nil
}

func deleteURL(shortCode string) (bool, error) {
	result, err := db.Exec("DELETE FROM urls WHERE short_code = ?", shortCode)
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
		"INSERT INTO stats (short_code) VALUES (?)",
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
		WHERE u.short_code = ?
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


