package configs

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var config *Config

func LoadConfigs() error {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Warning: Could not load .env file, relying on system environment variables")
	}

	config = &Config{}
	if err := env.Parse(config); err != nil {
		return err
	}
	log.Println("Config loaded successfully")
	return nil
}

func GetConfig() *Config {
	return config
}
