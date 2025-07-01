package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type Movie struct {
	ID        int
	Title     string
	Year      int
	Genre     string
	Streaming string
	Notes     string
	IMDBLink  string
}

var db *sql.DB

// InitDB initializes the SQLite database
func InitDB() error {
	// Get user home directory
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create data directory
	dataDir := filepath.Join(homeDir, ".local", "share", "homenet", "data")

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Database file path
	dbPath := filepath.Join(dataDir, "movies.db")

	// Open database
	db, err = sql.Open("sqlite3", dbPath)

	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Create movies table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS movies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		year INTEGER,
		genre TEXT,
		streaming TEXT,
		notes TEXT,
		imdb_link TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Add streaming column if it doesn't exist (for existing databases)
	addStreamingColumnSQL := `ALTER TABLE movies ADD COLUMN streaming TEXT;`

	db.Exec(addStreamingColumnSQL) // Ignore error if column already exists

	_, err = db.Exec(createTableSQL)

	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}

	return nil
}

// GetAllMovies retrieves all movies from the database
func GetAllMovies() ([]Movie, error) {
	rows, err := db.Query("SELECT id, title, year, genre, streaming, notes, imdb_link FROM movies ORDER BY created_at DESC")

	if err != nil {
		return nil, fmt.Errorf("failed to query movies: %w", err)
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var movie Movie
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genre, &movie.Streaming, &movie.Notes, &movie.IMDBLink)
		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

// AddMovie adds a new movie to the database and returns the ID
func AddMovie(movie Movie) (int, error) {
	query := `
	INSERT INTO movies (title, year, genre, streaming, notes, imdb_link)
	VALUES (?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, movie.Title, movie.Year, movie.Genre, movie.Streaming, movie.Notes, movie.IMDBLink)

	if err != nil {
		return 0, fmt.Errorf("failed to insert movie: %w", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return int(id), nil
}

// DeleteMovie deletes a movie by ID
func DeleteMovie(id int) error {
	query := "DELETE FROM movies WHERE id = ?"
	_, err := db.Exec(query, id)

	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	return nil
}

// GetMovieCount returns the total number of movies
func GetMovieCount() (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM movies").Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get movie count: %w", err)
	}

	return count, nil
}

// UpdateMovie updates an existing movie in the database
func UpdateMovie(movie Movie) error {
	query := `
	UPDATE movies 
	SET title = ?, year = ?, genre = ?, streaming = ?, notes = ?, imdb_link = ?
	WHERE id = ?`

	_, err := db.Exec(query, movie.Title, movie.Year, movie.Genre, movie.Streaming, movie.Notes, movie.IMDBLink, movie.ID)

	if err != nil {
		return fmt.Errorf("failed to update movie: %w", err)
	}

	return nil
}
