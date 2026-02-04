package scheduler

import (
	"testing"
	"time"

	"github.com/berkaycubuk/subtrack/internal/database"
	"github.com/berkaycubuk/subtrack/internal/services"
)

func setupScheduler(t *testing.T) (*Scheduler, *database.DB, *services.MockTelegramService) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	mockTg := &services.MockTelegramService{}
	subSvc := services.NewSubscriptionService(db, mockTg)
	sched := NewScheduler(subSvc)

	return sched, db, mockTg
}

func TestNewScheduler(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	mockTg := &services.MockTelegramService{}
	subSvc := services.NewSubscriptionService(db, mockTg)
	sched := NewScheduler(subSvc)

	if sched == nil {
		t.Fatal("NewScheduler() returned nil")
	}

	if sched.subSvc != subSvc {
		t.Error("NewScheduler() did not set subSvc")
	}
}

func TestScheduler_runCheck(t *testing.T) {
	sched, db, mockTg := setupScheduler(t)

	now := time.Now()

	subs := []*database.Subscription{
		{
			Name:        "Due in 2 days",
			Price:       10.00,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: now.Add(2 * 24 * time.Hour),
		},
		{
			Name:        "Due in 4 days",
			Price:       20.00,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: now.Add(4 * 24 * time.Hour),
		},
		{
			Name:        "Past due",
			Price:       30.00,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: now.Add(-5 * 24 * time.Hour),
		},
		{
			Name:        "Far in future",
			Price:       40.00,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: now.Add(30 * 24 * time.Hour),
		},
	}

	for _, sub := range subs {
		if err := db.CreateSubscription(sub); err != nil {
			t.Fatalf("failed to create test subscription: %v", err)
		}
	}

	notificationsSent := 0
	mockTg.SendNotificationFunc = func(name string, price float64, currency string, days int, cycle string, paymentDate string) error {
		notificationsSent++
		return nil
	}

	sched.runCheck()

	if notificationsSent != 2 {
		t.Errorf("runCheck() sent %d notifications, want 2", notificationsSent)
	}

	allSubs, err := db.GetAllSubscriptions()
	if err != nil {
		t.Fatalf("failed to get subscriptions: %v", err)
	}

	for _, sub := range allSubs {
		if sub.Name == "Past due" && sub.PaymentDate.Before(now) {
			t.Errorf("runCheck() did not update past due payment date for %s", sub.Name)
		}
	}
}

func TestScheduler_StopGracefully(t *testing.T) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	mockTg := &services.MockTelegramService{}
	subSvc := services.NewSubscriptionService(db, mockTg)
	sched := NewScheduler(subSvc)

	done := make(chan bool)
	go func() {
		time.Sleep(100 * time.Millisecond)
		sched.StopGracefully()
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Error("StopGracefully() did not complete within timeout")
	}
}
