package services

import (
	"fmt"
	"log"
	"strconv"

	"github.com/berkaycubuk/subtrack/internal/database"
	"github.com/berkaycubuk/subtrack/internal/utils"
)

type TelegramNotifier interface {
	SendNotification(name string, price float64, currency string, days int, cycle string, paymentDate string) error
}

type SubscriptionService struct {
	db *database.DB
	tg TelegramNotifier
}

func NewSubscriptionService(db *database.DB, tg TelegramNotifier) *SubscriptionService {
	return &SubscriptionService{
		db: db,
		tg: tg,
	}
}

func (s *SubscriptionService) AddSubscription(name, priceStr, currency, cycle, paymentDateStr string) error {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fmt.Errorf("invalid price format: %w", err)
	}

	paymentDate, err := utils.ParseDate(paymentDateStr)
	if err != nil {
		return fmt.Errorf("invalid payment date format (use DD-MM-YYYY): %w", err)
	}

	if cycle != "monthly" && cycle != "yearly" {
		return fmt.Errorf("cycle must be 'monthly' or 'yearly'")
	}

	sub := &database.Subscription{
		Name:        name,
		Price:       price,
		Currency:    currency,
		Cycle:       cycle,
		PaymentDate: paymentDate,
	}

	return s.db.CreateSubscription(sub)
}

func (s *SubscriptionService) UpdateSubscription(id uint, name, priceStr, currency, cycle, paymentDateStr string) error {
	sub, err := s.db.GetSubscriptionByID(id)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	if name != "" {
		sub.Name = name
	}

	if priceStr != "" {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return fmt.Errorf("invalid price format: %w", err)
		}
		sub.Price = price
	}

	if currency != "" {
		sub.Currency = currency
	}

	if cycle != "" {
		if cycle != "monthly" && cycle != "yearly" {
			return fmt.Errorf("cycle must be 'monthly' or 'yearly'")
		}
		sub.Cycle = cycle
	}

	if paymentDateStr != "" {
		paymentDate, err := utils.ParseDate(paymentDateStr)
		if err != nil {
			return fmt.Errorf("invalid payment date format (use DD-MM-YYYY): %w", err)
		}
		sub.PaymentDate = paymentDate
	}

	return s.db.UpdateSubscription(sub)
}

func (s *SubscriptionService) DeleteSubscription(id uint) error {
	return s.db.DeleteSubscription(id)
}

func (s *SubscriptionService) ListSubscriptions() ([]database.Subscription, error) {
	return s.db.GetAllSubscriptions()
}

func (s *SubscriptionService) CheckUpcomingPayments() ([]database.Subscription, error) {
	subs, err := s.db.GetUpcomingPayments(5)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *SubscriptionService) SendNotifications(subs []database.Subscription) error {
	for _, sub := range subs {
		days := utils.DaysUntil(sub.PaymentDate)
		if days >= 0 && days < 5 {
			paymentDateStr := utils.FormatDate(sub.PaymentDate)
			err := s.tg.SendNotification(sub.Name, sub.Price, sub.Currency, days, sub.Cycle, paymentDateStr)
			if err != nil {
				log.Printf("Failed to send notification for %s: %v", sub.Name, err)
			} else {
				log.Printf("Sent notification for %s (payment in %d days)", sub.Name, days)
			}
		}
	}
	return nil
}

func (s *SubscriptionService) UpdatePastDuePayments() error {
	subs, err := s.db.GetPastDuePayments()
	if err != nil {
		return err
	}

	for _, sub := range subs {
		newPaymentDate, err := utils.CalculateNextPaymentDate(sub.PaymentDate, sub.Cycle)
		if err != nil {
			log.Printf("Failed to calculate next payment date for %s: %v", sub.Name, err)
			continue
		}

		sub.PaymentDate = newPaymentDate
		if err := s.db.UpdateSubscription(&sub); err != nil {
			log.Printf("Failed to update payment date for %s: %v", sub.Name, err)
		} else {
			log.Printf("Updated payment date for %s to %s", sub.Name, utils.FormatDate(newPaymentDate))
		}
	}

	return nil
}
