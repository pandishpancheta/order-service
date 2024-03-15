package config

import (
	"os"
	"strings"
)

type Config struct {
	TCP_PORT string

	DB_HOST string
	DB_PORT string
	DB_USER string
	DB_PASS string
	DB_NAME string
}

func LoadConfig() *Config {
	return &Config{
		TCP_PORT: strings.TrimSpace(os.Getenv("TCP_PORT")),

		DB_HOST: strings.TrimSpace(os.Getenv("DB_HOST")),
		DB_PORT: strings.TrimSpace(os.Getenv("DB_PORT")),
		DB_USER: strings.TrimSpace(os.Getenv("DB_USER")),
		DB_PASS: strings.TrimSpace(os.Getenv("DB_PASS")),
		DB_NAME: strings.TrimSpace(os.Getenv("DB_NAME")),
	}
}
