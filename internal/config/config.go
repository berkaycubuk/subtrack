package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	TelegramChatID   string
	DBPath           string
	WebUsername       string
	WebPassword      string
	WebPort          string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}

	chatID := os.Getenv("TELEGRAM_CHAT_ID")
	if chatID == "" {
		return nil, fmt.Errorf("TELEGRAM_CHAT_ID is required")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "subtrack.db"
	}

	webUsername := os.Getenv("WEB_USERNAME")
	if webUsername == "" {
		return nil, fmt.Errorf("WEB_USERNAME is required")
	}

	webPassword := os.Getenv("WEB_PASSWORD")
	if webPassword == "" {
		return nil, fmt.Errorf("WEB_PASSWORD is required")
	}

	webPort := os.Getenv("WEB_PORT")
	if webPort == "" {
		webPort = "8080"
	}

	return &Config{
		TelegramBotToken: botToken,
		TelegramChatID:   chatID,
		DBPath:           dbPath,
		WebUsername:       webUsername,
		WebPassword:      webPassword,
		WebPort:          webPort,
	}, nil
}
