package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pwnderpants/homenet/internal/logger"
)

type Movie struct {
	ID           int
	Title        string
	Year         int
	Genre        string
	Streaming    string
	Notes        string
	IMDBLink     string
	AvailableNow bool
}

type TVShow struct {
	ID           int
	Title        string
	Year         int
	Genre        string
	Streaming    string
	Notes        string
	IMDBLink     string
	ActiveSeason bool
}

var db *sql.DB

// InitDB initializes the database connection and creates tables if they don't exist
func InitDB(dataDir, dbName string) error {
	logger.Info("Initializing database...")

	// Create data directory
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		logger.ErrorWithErr("Failed to create data directory", err)

		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Database file path
	dbPath := filepath.Join(dataDir, dbName+".db")

	logger.Info("Database path: %s", dbPath)

	// Open database
	var err error

	db, err = sql.Open("sqlite3", dbPath)

	if err != nil {
		logger.ErrorWithErr("Failed to open database", err)

		return fmt.Errorf("failed to open database: %w", err)
	}

	// Create movies table if it doesn't exist
	createMoviesTableSQL := `
	CREATE TABLE IF NOT EXISTS movies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		year INTEGER,
		genre TEXT,
		streaming TEXT,
		notes TEXT,
		imdb_link TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		available_now INTEGER DEFAULT 0
	);`

	// Create tv_shows table if it doesn't exist
	createTVShowsTableSQL := `
	CREATE TABLE IF NOT EXISTS tv_shows (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		year INTEGER,
		genre TEXT,
		streaming TEXT,
		notes TEXT,
		imdb_link TEXT,
		active_season INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	// Add streaming column if it doesn't exist (for existing databases)
	addStreamingColumnSQL := `ALTER TABLE movies ADD COLUMN streaming TEXT;`
	addTVShowsStreamingColumnSQL := `ALTER TABLE tv_shows ADD COLUMN streaming TEXT;`
	addActiveSeasonColumnSQL := `ALTER TABLE tv_shows ADD COLUMN active_season INTEGER DEFAULT 0;`

	db.Exec(addStreamingColumnSQL)        // Ignore error if column already exists
	db.Exec(addTVShowsStreamingColumnSQL) // Ignore error if column already exists
	db.Exec(addActiveSeasonColumnSQL)     // Ignore error if column already exists

	_, err = db.Exec(createMoviesTableSQL)

	if err != nil {
		logger.ErrorWithErr("Failed to create movies table", err)

		return fmt.Errorf("failed to create movies table: %w", err)
	}

	_, err = db.Exec(createTVShowsTableSQL)

	if err != nil {
		logger.ErrorWithErr("Failed to create tv_shows table", err)

		return fmt.Errorf("failed to create tv_shows table: %w", err)
	}

	logger.Info("Database initialized successfully")

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
	rows, err := db.Query("SELECT id, title, year, genre, streaming, notes, imdb_link, available_now FROM movies ORDER BY available_now DESC, year DESC, created_at DESC")

	if err != nil {
		return nil, fmt.Errorf("failed to query movies: %w", err)
	}

	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var movie Movie
		var availableNowInt int

		err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genre, &movie.Streaming, &movie.Notes, &movie.IMDBLink, &availableNowInt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan movie: %w", err)
		}

		movie.AvailableNow = availableNowInt == 1
		movies = append(movies, movie)
	}

	return movies, nil
}

// AddMovie adds a new movie to the database and returns the ID
func AddMovie(movie Movie) (int, error) {
	logger.Info("Adding movie: %s (%d)", movie.Title, movie.Year)

	query := `
	INSERT INTO movies (title, year, genre, streaming, notes, imdb_link, available_now)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, movie.Title, movie.Year, movie.Genre, movie.Streaming, movie.Notes, movie.IMDBLink, boolToInt(movie.AvailableNow))

	if err != nil {
		logger.ErrorWithErr("Failed to insert movie", err)

		return 0, fmt.Errorf("failed to insert movie: %w", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		logger.ErrorWithErr("Failed to get last insert id", err)
		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	logger.Info("Movie added successfully with ID: %d", id)

	return int(id), nil
}

// DeleteMovie deletes a movie by ID
func DeleteMovie(id int) error {
	logger.Info("Deleting movie with ID: %d", id)

	query := "DELETE FROM movies WHERE id = ?"
	_, err := db.Exec(query, id)

	if err != nil {
		logger.ErrorWithErr("Failed to delete movie", err)

		return fmt.Errorf("failed to delete movie: %w", err)
	}

	logger.Info("Movie deleted successfully")

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
	logger.Info("Updating movie with ID: %d, title: %s", movie.ID, movie.Title)

	query := `
	UPDATE movies 
	SET title = ?, year = ?, genre = ?, streaming = ?, notes = ?, imdb_link = ?, available_now = ?
	WHERE id = ?`

	_, err := db.Exec(query, movie.Title, movie.Year, movie.Genre, movie.Streaming, movie.Notes, movie.IMDBLink, boolToInt(movie.AvailableNow), movie.ID)

	if err != nil {
		logger.ErrorWithErr("Failed to update movie", err)

		return fmt.Errorf("failed to update movie: %w", err)
	}

	logger.Info("Movie updated successfully")

	return nil
}

// GetRandomMovie retrieves a random movie from the database
func GetRandomMovie() (*Movie, error) {
	query := "SELECT id, title, year, genre, streaming, notes, imdb_link FROM movies WHERE available_now = 1 ORDER BY RANDOM() LIMIT 1"

	var movie Movie

	err := db.QueryRow(query).Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Genre, &movie.Streaming, &movie.Notes, &movie.IMDBLink)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No movies found
		}

		return nil, fmt.Errorf("failed to get random movie: %w", err)
	}

	return &movie, nil
}

// GetAllTVShows retrieves all TV shows from the database
func GetAllTVShows() ([]TVShow, error) {
	rows, err := db.Query("SELECT id, title, year, genre, streaming, notes, imdb_link, active_season FROM tv_shows ORDER BY active_season DESC, year DESC, created_at DESC")

	if err != nil {
		return nil, fmt.Errorf("failed to query tv shows: %w", err)
	}

	defer rows.Close()

	var tvShows []TVShow

	for rows.Next() {
		var tvShow TVShow
		var activeSeasonInt int

		err := rows.Scan(&tvShow.ID, &tvShow.Title, &tvShow.Year, &tvShow.Genre, &tvShow.Streaming, &tvShow.Notes, &tvShow.IMDBLink, &activeSeasonInt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan tv show: %w", err)
		}

		tvShow.ActiveSeason = activeSeasonInt == 1
		tvShows = append(tvShows, tvShow)
	}

	return tvShows, nil
}

// AddTVShow adds a new TV show to the database and returns the ID
func AddTVShow(tvShow TVShow) (int, error) {
	logger.Info("Adding TV show: %s (%d)", tvShow.Title, tvShow.Year)

	query := `
	INSERT INTO tv_shows (title, year, genre, streaming, notes, imdb_link, active_season)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, tvShow.Title, tvShow.Year, tvShow.Genre, tvShow.Streaming, tvShow.Notes, tvShow.IMDBLink, boolToInt(tvShow.ActiveSeason))

	if err != nil {
		logger.ErrorWithErr("Failed to insert tv show", err)

		return 0, fmt.Errorf("failed to insert tv show: %w", err)
	}

	id, err := result.LastInsertId()

	if err != nil {
		logger.ErrorWithErr("Failed to get last insert id", err)

		return 0, fmt.Errorf("failed to get last insert id: %w", err)
	}

	logger.Info("TV show added successfully with ID: %d", id)

	return int(id), nil
}

// DeleteTVShow deletes a TV show by ID
func DeleteTVShow(id int) error {
	logger.Info("Deleting TV show with ID: %d", id)

	query := "DELETE FROM tv_shows WHERE id = ?"
	_, err := db.Exec(query, id)

	if err != nil {
		logger.ErrorWithErr("Failed to delete tv show", err)

		return fmt.Errorf("failed to delete tv show: %w", err)
	}

	logger.Info("TV show deleted successfully")

	return nil
}

// GetTVShowCount returns the total number of TV shows
func GetTVShowCount() (int, error) {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM tv_shows").Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("failed to get tv show count: %w", err)
	}

	return count, nil
}

// UpdateTVShow updates an existing TV show in the database
func UpdateTVShow(tvShow TVShow) error {
	logger.Info("Updating TV show with ID: %d, title: %s", tvShow.ID, tvShow.Title)

	query := `
	UPDATE tv_shows 
	SET title = ?, year = ?, genre = ?, streaming = ?, notes = ?, imdb_link = ?, active_season = ?
	WHERE id = ?`

	_, err := db.Exec(query, tvShow.Title, tvShow.Year, tvShow.Genre, tvShow.Streaming, tvShow.Notes, tvShow.IMDBLink, boolToInt(tvShow.ActiveSeason), tvShow.ID)

	if err != nil {
		logger.ErrorWithErr("Failed to update tv show", err)

		return fmt.Errorf("failed to update tv show: %w", err)
	}

	logger.Info("TV show updated successfully")

	return nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}

	return 0
}
