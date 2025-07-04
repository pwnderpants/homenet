package main

import (
	"log"
	"os"
	"strings"

	"github.com/pwnderpants/homenet/internal/logger"
	"github.com/pwnderpants/homenet/internal/server"
)

func main() {
	// Configure logging level from environment variable
	logLevel := strings.ToUpper(os.Getenv("LOG_LEVEL"))

	switch logLevel {
	case "DEBUG":
		logger.SetGlobalLevel(logger.DEBUG)
		log.Println("Log level set to DEBUG")
	case "INFO":
		logger.SetGlobalLevel(logger.INFO)
		log.Println("Log level set to INFO")
	case "WARN":
		logger.SetGlobalLevel(logger.WARN)
		log.Println("Log level set to WARN")
	case "ERROR":
		logger.SetGlobalLevel(logger.ERROR)
		log.Println("Log level set to ERROR")
	default:
		logger.SetGlobalLevel(logger.INFO)
		log.Println("Log level set to INFO (default)")
	}

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
