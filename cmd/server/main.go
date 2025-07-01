package main

import (
	"log"

	"github.com/pwnderpants/homenet/internal/server"
)

func main() {
	// Create and configure server
	srv := server.New()
	srv.SetupRoutes()

	// Start the server
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
