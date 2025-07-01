package main

import (
	"log"
	"os"

	"github.com/pwnderpants/homenet/internal/server"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	// Start the server
	if err := server.StartServer(port); err != nil {
		log.Fatal(err)
	}
}
