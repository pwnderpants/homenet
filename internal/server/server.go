package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/pwnderpants/homenet/internal/database"
	"github.com/pwnderpants/homenet/internal/handlers"
)

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

	// Handle main routes
	http.HandleFunc("/", handlers.HomeHandler)

	// Movie board routes
	http.HandleFunc("/movie-board", handlers.MovieBoardHandler)
	http.HandleFunc("/movie-board/add", handlers.AddMovieHandler)
	http.HandleFunc("/movie-board/edit", handlers.EditMovieHandler)
	http.HandleFunc("/movie-board/delete/", handlers.DeleteMovieHandler)
	http.HandleFunc("/movie-board/random", handlers.RandomMovieHandler)

	// TV Shows board routes
	http.HandleFunc("/tv-shows-board", handlers.TVShowBoardHandler)
	http.HandleFunc("/tv-shows-board/add", handlers.AddTVShowHandler)
	http.HandleFunc("/tv-shows-board/edit", handlers.EditTVShowHandler)
	http.HandleFunc("/tv-shows-board/delete/", handlers.DeleteTVShowHandler)

	// Fortune route
	http.HandleFunc("/fortune", handlers.FortuneHandler)

	// AI routes
	http.HandleFunc("/ai", handlers.AiHandler)
}

// StartServer initializes and starts the HTTP server
func StartServer(port string) error {
	// Initialize database
	if err := database.InitDB(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create server instance and set up routes
	server := New()
	server.SetupRoutes()

	// Start server
	log.Printf("Server starting on http://localhost:%s", port)
	return http.ListenAndServe(server.addr, nil)
}
