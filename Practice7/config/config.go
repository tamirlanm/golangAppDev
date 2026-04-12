package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL     string
	JWTSecret string
	Port      string
}

func New() (*Config, error) {
	_ = godotenv.Load()
	return &Config{
		DBURL:     os.Getenv("DATABASE_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Port:      os.Getenv("PORT"),
	}, nil
}
