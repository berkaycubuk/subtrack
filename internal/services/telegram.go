package services

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramService struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramService(botToken string, chatID int64) (*TelegramService, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &TelegramService{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func (t *TelegramService) SendNotification(name string, price float64, currency string, days int, cycle string, paymentDate string) error {
	message := fmt.Sprintf("📢 Subscription Alert: %s\n💰 Price: %.2f %s\n📅 Payment in: %d days\n🔄 Cycle: %s\n📆 Next payment: %s",
		name, price, currency, days, cycle, paymentDate)

	msg := tgbotapi.NewMessage(t.chatID, message)
	msg.ParseMode = "Markdown"

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (t *TelegramService) SendMessage(text string) error {
	msg := tgbotapi.NewMessage(t.chatID, text)
	msg.ParseMode = "Markdown"

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (t *TelegramService) HealthCheck() error {
	_, err := t.bot.GetMe()
	if err != nil {
		return fmt.Errorf("bot health check failed: %w", err)
	}
	return nil
}
