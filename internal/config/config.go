package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	HostPort    string
	DatabaseUrl string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Env:         os.Getenv("ENV"),
		HostPort:    os.Getenv("HOST_PORT"),
		DatabaseUrl: os.Getenv("DATABASE_URL"),
	}
}
