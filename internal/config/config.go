package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Dsn       string
	JwtSecret string
}

func LoadEnv() *Config {
	// In dev, .env exists and is loaded. In production (Render, etc.),
	// .env is not present and env vars come from the platform.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
	cfg := &Config{
		Port:      os.Getenv("PORT"),
		Dsn:       os.Getenv("DSN"),
		JwtSecret: os.Getenv("JWT_SECRET"),
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.Dsn == "" {
		log.Fatal("DSN environment variable is required")
	}
	if cfg.JwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
	return cfg
}
