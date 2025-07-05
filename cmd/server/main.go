package main

import (
	"log"
	"strings"

	"github.com/pwnderpants/homenet/internal/config"
	"github.com/pwnderpants/homenet/internal/logger"
	"github.com/pwnderpants/homenet/internal/server"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Configure logging level from config
	logLevel := strings.ToUpper(cfg.Logging.Level)

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

	// Start the server with configured port
	if err := server.StartServer(cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
