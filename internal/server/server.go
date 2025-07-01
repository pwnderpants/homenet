package server

import (
	"log"
	"net/http"
	"os"

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
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("Server starting on http://localhost%s", s.addr)

	return http.ListenAndServe(s.addr, nil)
}
