package utils

import (
	"log"
	"os"
)

type Config struct {
	AuthServerURL string
	TodoAppURL    string
}

func LoadConfig() *Config {
	authServerURL := os.Getenv("AUTH_SERVER_URL")
	todoAppURL := os.Getenv("TODO_APP_URL")

	if authServerURL == "" || todoAppURL == "" {
		log.Fatal("請設置環境變數 AUTH_SERVER_URL 和 TODO_APP_URL")
	}

	return &Config{
		AuthServerURL: authServerURL,
		TodoAppURL:    todoAppURL,
	}
}
