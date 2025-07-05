package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	Ollama struct {
		Host      string `json:"host"`
		ModelName string `json:"model_name"`
	} `json:"ollama"`
	Logging struct {
		Level string `json:"level"`
	} `json:"logging"`
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
	Database struct {
		DataDir string `json:"data_dir"`
		DBName  string `json:"db_name"`
	} `json:"database"`
	Static struct {
		Dir string `json:"dir"`
	} `json:"static"`
	Fortune struct {
		Command     string `json:"command"`
		Args        string `json:"args"`
		FallbackMsg string `json:"fallback_msg"`
	} `json:"fortune"`
}

// LoadConfig loads the configuration from the user's home directory
func LoadConfig() (*Config, error) {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".config", "homenet")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Config file path
	configPath := filepath.Join(configDir, "config.json")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config file
		if err := createDefaultConfig(configPath); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
	}

	// Read and parse config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults for any missing values
	setDefaults(&config)

	return &config, nil
}

// createDefaultConfig creates a default configuration file
func createDefaultConfig(configPath string) error {
	defaultConfig := Config{}

	// Set default Ollama configuration
	defaultConfig.Ollama.Host = "http://chadgpt.gotpwnd.org:11434"
	defaultConfig.Ollama.ModelName = "llama3.2:latest"

	// Set default logging configuration
	defaultConfig.Logging.Level = "INFO"

	// Set default server configuration
	defaultConfig.Server.Host = "localhost"
	defaultConfig.Server.Port = "8080"

	// Set default database configuration
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "~" // Fallback if we can't get home directory
	}
	defaultConfig.Database.DataDir = filepath.Join(homeDir, ".local", "share", "homenet", "data")
	defaultConfig.Database.DBName = "homenet"

	// Set default static files configuration
	defaultConfig.Static.Dir = "web/static"

	// Set default fortune command configuration
	defaultConfig.Fortune.Command = "/usr/games/fortune"
	defaultConfig.Fortune.Args = "-s"
	defaultConfig.Fortune.FallbackMsg = "Hello World!"

	// Marshal to JSON with pretty formatting
	data, err := json.MarshalIndent(&defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}

	return nil
}

// setDefaults sets default values for any missing configuration fields
func setDefaults(config *Config) {
	if config.Ollama.Host == "" {
		config.Ollama.Host = "http://chadgpt.gotpwnd.org:11434"
	}

	if config.Ollama.ModelName == "" {
		config.Ollama.ModelName = "llama3.2:latest"
	}

	if config.Logging.Level == "" {
		config.Logging.Level = "INFO"
	}

	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}

	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}

	if config.Database.DataDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			homeDir = "~" // Fallback if we can't get home directory
		}
		config.Database.DataDir = filepath.Join(homeDir, ".local", "share", "homenet", "data")
	}

	if config.Database.DBName == "" {
		config.Database.DBName = "movies"
	}

	if config.Static.Dir == "" {
		config.Static.Dir = "web/static"
	}

	if config.Fortune.Command == "" {
		config.Fortune.Command = "/usr/games/fortune"
	}

	if config.Fortune.Args == "" {
		config.Fortune.Args = "-s"
	}

	if config.Fortune.FallbackMsg == "" {
		config.Fortune.FallbackMsg = "Built with ❤️ using HTMX, Go, and Tailwind CSS"
	}
}
