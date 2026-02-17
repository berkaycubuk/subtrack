package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/berkaycubuk/subtrack/internal/config"
	"github.com/berkaycubuk/subtrack/internal/database"
	"github.com/berkaycubuk/subtrack/internal/scheduler"
	"github.com/berkaycubuk/subtrack/internal/services"
	"github.com/berkaycubuk/subtrack/internal/web"
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
	if err := sched.StartCron(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}

	srv := web.NewServer(subSvc, cfg.WebUsername, cfg.WebPassword)

	go func() {
		if err := srv.Start(":" + cfg.WebPort); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Web server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	sched.Stop()
	if err := srv.Shutdown(); err != nil {
		log.Printf("Web server shutdown error: %v", err)
	}

	os.Exit(0)
}
