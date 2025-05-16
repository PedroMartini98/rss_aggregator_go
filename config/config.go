package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBURL string
}

func LoadConfig() (*Config, error) {
	if os.Getenv("DB_ULR") == "" {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("Failed to load .env file: %v", err)
		}
	}

	config := &Config{
		Port:  os.Getenv("PORT"),
		DBURL: os.Getenv("DB_URL"),
	}

	if config.DBURL == "" {
		return nil, fmt.Errorf("Failed to find DBURL in .env")
	}

	if config.Port == "" {
		return nil, fmt.Errorf("Failed to find Port in .env")
	}

	return config, nil

}
