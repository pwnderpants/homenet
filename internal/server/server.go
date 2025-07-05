package server

import (
	"fmt"
	"net/http"

	"github.com/pwnderpants/homenet/internal/config"
	"github.com/pwnderpants/homenet/internal/database"
	"github.com/pwnderpants/homenet/internal/handlers"
	"github.com/pwnderpants/homenet/internal/logger"
)

// Server represents the HTTP server
type Server struct {
	addr   string
	config *config.Config
}

// New creates a new server instance
func New(cfg *config.Config) *Server {
	return &Server{
		addr:   ":" + cfg.Server.Port,
		config: cfg,
	}
}

// SetupRoutes configures all the routes for the server
func (s *Server) SetupRoutes() {
	// Serve static files
	fs := http.FileServer(http.Dir(s.config.Static.Dir))

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
	http.HandleFunc("/fortune", s.createFortuneHandler())

	// AI routes
	http.HandleFunc("/ai", handlers.AiHandler)
	http.HandleFunc("/ai/query", s.createAIQueryHandler())
}

// createAIQueryHandler creates a handler that uses the server's configuration
func (s *Server) createAIQueryHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		handlers.AIQueryHandlerWithConfig(w, r, s.config)
	}
}

// createFortuneHandler creates a handler that uses the server's configuration
func (s *Server) createFortuneHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlers.FortuneHandlerWithConfig(w, r, s.config)
	}
}

// StartServer initializes and starts the HTTP server
func StartServer(port string) error {
	// Load configuration
	cfg, err := config.LoadConfig()

	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize database with configuration
	if err := database.InitDB(cfg.Database.DataDir, cfg.Database.DBName); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create server instance and set up routes
	server := New(cfg)

	server.SetupRoutes()

	// Start server
	logger.Info("Server starting on http://%s:%s", cfg.Server.Host, port)

	return http.ListenAndServe(server.addr, nil)
}
