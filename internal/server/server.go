package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pwnderpants/homenet/internal/database"
	"github.com/pwnderpants/homenet/internal/handlers"
)

// Server represents the HTTP server
type Server struct {
	addr string
}

// New creates a new server instance
func New() *Server {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	return &Server{
		addr: ":" + port,
	}
}

// SetupRoutes configures all the routes for the server
func (s *Server) SetupRoutes() {
	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/movie-board", handlers.MovieBoardHandler)
	http.HandleFunc("/movie-board/add", handlers.AddMovieHandler)
	http.HandleFunc("/movie-board/edit", handlers.EditMovieHandler)
	http.HandleFunc("/movie-board/delete/", handlers.DeleteMovieHandler)

	// TV Shows Board routes
	http.HandleFunc("/tv-shows-board", handlers.TVShowBoardHandler)
	http.HandleFunc("/tv-shows-board/add", handlers.AddTVShowHandler)
	http.HandleFunc("/tv-shows-board/edit", handlers.EditTVShowHandler)
	http.HandleFunc("/tv-shows-board/delete/", handlers.DeleteTVShowHandler)
}

// StartServer initializes and starts the HTTP server
func StartServer(port string) error {
	// Initialize database
	if err := database.InitDB(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Set up routes
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/movie-board", handlers.MovieBoardHandler)
	http.HandleFunc("/movie-board/add", handlers.AddMovieHandler)
	http.HandleFunc("/movie-board/edit", handlers.EditMovieHandler)
	http.HandleFunc("/movie-board/delete/", handlers.DeleteMovieHandler)

	// TV Shows Board routes
	http.HandleFunc("/tv-shows-board", handlers.TVShowBoardHandler)
	http.HandleFunc("/tv-shows-board/add", handlers.AddTVShowHandler)
	http.HandleFunc("/tv-shows-board/edit", handlers.EditTVShowHandler)
	http.HandleFunc("/tv-shows-board/delete/", handlers.DeleteTVShowHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	log.Printf("Server starting on http://localhost:%s", port)
	return http.ListenAndServe(":"+port, nil)
}
