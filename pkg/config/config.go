package config

import (
	"os"
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
		TCP_PORT: os.Getenv("TCP_PORT"),

		DB_HOST: os.Getenv("DB_HOST"),
		DB_PORT: os.Getenv("DB_PORT"),
		DB_USER: os.Getenv("DB_USER"),
		DB_PASS: os.Getenv("DB_PASS"),
		DB_NAME: os.Getenv("DB_NAME"),
	}
}
