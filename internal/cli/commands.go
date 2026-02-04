package cli

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/berkaycubuk/subtrack/internal/config"
	"github.com/berkaycubuk/subtrack/internal/database"
	"github.com/berkaycubuk/subtrack/internal/services"
	"github.com/berkaycubuk/subtrack/internal/utils"
)

type CLI struct {
	cfg    *config.Config
	db     *database.DB
	subSvc *services.SubscriptionService
	tgSvc  *services.TelegramService
}

func New() (*CLI, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	db, err := database.New(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	chatID, err := strconv.ParseInt(cfg.TelegramChatID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid TELEGRAM_CHAT_ID: %w", err)
	}

	tgSvc, err := services.NewTelegramService(cfg.TelegramBotToken, chatID)
	if err != nil {
		return nil, err
	}

	subSvc := services.NewSubscriptionService(db, tgSvc)

	return &CLI{
		cfg:    cfg,
		db:     db,
		subSvc: subSvc,
		tgSvc:  tgSvc,
	}, nil
}

func (c *CLI) Add(name, price, currency, cycle, paymentDate string) error {
	if err := c.subSvc.AddSubscription(name, price, currency, cycle, paymentDate); err != nil {
		return err
	}
	fmt.Println("✓ Subscription added successfully")
	return nil
}

func (c *CLI) List() error {
	subs, err := c.subSvc.ListSubscriptions()
	if err != nil {
		return err
	}

	if len(subs) == 0 {
		fmt.Println("No subscriptions found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "ID\tName\tPrice\tCurrency\tCycle\tPayment Date\n")
	fmt.Fprintf(w, "--\t----\t-----\t--------\t-----\t------------\n")

	for _, sub := range subs {
		paymentDateStr := utils.FormatDate(sub.PaymentDate)
		fmt.Fprintf(w, "%d\t%s\t%.2f\t%s\t%s\t%s\n",
			sub.ID, sub.Name, sub.Price, sub.Currency, sub.Cycle, paymentDateStr)
	}

	w.Flush()
	return nil
}

func (c *CLI) Update(idStr, name, price, currency, cycle, paymentDate string) error {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid subscription ID: %w", err)
	}

	if err := c.subSvc.UpdateSubscription(uint(id), name, price, currency, cycle, paymentDate); err != nil {
		return err
	}
	fmt.Println("✓ Subscription updated successfully")
	return nil
}

func (c *CLI) Delete(idStr string) error {
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid subscription ID: %w", err)
	}

	if err := c.subSvc.DeleteSubscription(uint(id)); err != nil {
		return err
	}
	fmt.Println("✓ Subscription deleted successfully")
	return nil
}

func (c *CLI) Check() error {
	fmt.Println("Checking upcoming payments...")

	if err := c.subSvc.UpdatePastDuePayments(); err != nil {
		fmt.Printf("Error updating past due payments: %v\n", err)
	}

	subs, err := c.subSvc.CheckUpcomingPayments()
	if err != nil {
		return err
	}

	if len(subs) == 0 {
		fmt.Println("No upcoming payments found")
		return nil
	}

	fmt.Printf("Found %d subscriptions with upcoming payments:\n\n", len(subs))

	for _, sub := range subs {
		days := utils.DaysUntil(sub.PaymentDate)
		paymentDateStr := utils.FormatDate(sub.PaymentDate)
		fmt.Printf("📢 %s\n", sub.Name)
		fmt.Printf("   💰 Price: %.2f %s\n", sub.Price, sub.Currency)
		fmt.Printf("   📅 Payment in: %d days\n", days)
		fmt.Printf("   🔄 Cycle: %s\n", sub.Cycle)
		fmt.Printf("   📆 Next payment: %s\n\n", paymentDateStr)
	}

	if err := c.subSvc.SendNotifications(subs); err != nil {
		fmt.Printf("Error sending notifications: %v\n", err)
	}

	return nil
}

func (c *CLI) Health() error {
	if err := c.tgSvc.HealthCheck(); err != nil {
		return fmt.Errorf("Telegram bot health check failed: %w", err)
	}
	fmt.Println("✓ Telegram bot is healthy")
	return nil
}
