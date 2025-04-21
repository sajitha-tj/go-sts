package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// getConfigFromEnv retrieves the value of the given key from environment variables.
func getConfigFromEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s not set", key)
	}
	return value
}

// getConfigFromEnvWithDefault retrieves the value of the given key from environment variables, or returns the default value if not set.
func getConfigFromEnvWithDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func readConfigFromFile(filePath string) string {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", filePath, err)
	}
	return string(data)
}

// Config holds the configuration values for the application.
type Config struct {
	PORT string
	SIGNING_SECRET string
}
// newConfig initializes a new Config instance with values from environment variables.
func newConfig() *Config {
	return &Config{
		PORT: getConfigFromEnvWithDefault("PORT", "8080"),
		SIGNING_SECRET: readConfigFromFile(getConfigFromEnv("SIGNING_SECRET_FILE_PATH")),
	}
}

var config *Config
// GetConfigInstance returns the singleton instance of Config.
func GetConfigInstance() *Config {
	if config == nil {
		config = newConfig()
	}
	return config
}