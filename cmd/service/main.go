package main

import (
	"log"
	"os"
	"strconv"

	"github.com/berkaycubuk/subtrack/internal/config"
	"github.com/berkaycubuk/subtrack/internal/database"
	"github.com/berkaycubuk/subtrack/internal/scheduler"
	"github.com/berkaycubuk/subtrack/internal/services"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	chatID, err := strconv.ParseInt(cfg.TelegramChatID, 10, 64)
	if err != nil {
		log.Fatalf("Invalid TELEGRAM_CHAT_ID: %v", err)
	}

	tgSvc, err := services.NewTelegramService(cfg.TelegramBotToken, chatID)
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot: %v", err)
	}

	subSvc := services.NewSubscriptionService(db, tgSvc)

	sched := scheduler.NewScheduler(subSvc)

	if err := sched.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}

	os.Exit(0)
}
