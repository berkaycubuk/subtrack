package services

import (
	"testing"
	"time"

	"github.com/berkaycubuk/subtrack/internal/database"
)

func setupSubscriptionService(t *testing.T) (*SubscriptionService, *database.DB, *MockTelegramService) {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	mockTg := &MockTelegramService{}
	subSvc := NewSubscriptionService(db, mockTg)

	return subSvc, db, mockTg
}

func TestSubscriptionService_AddSubscription(t *testing.T) {
	subSvc, _, _ := setupSubscriptionService(t)

	tests := []struct {
		name        string
		price       string
		currency    string
		cycle       string
		paymentDate string
		wantErr     bool
	}{
		{
			name:        "valid subscription",
			price:       "15.99",
			currency:    "USD",
			cycle:       "monthly",
			paymentDate: "15-02-2025",
			wantErr:     false,
		},
		{
			name:        "invalid price",
			price:       "invalid",
			currency:    "USD",
			cycle:       "monthly",
			paymentDate: "15-02-2025",
			wantErr:     true,
		},
		{
			name:        "invalid date",
			price:       "15.99",
			currency:    "USD",
			cycle:       "monthly",
			paymentDate: "2025-02-15",
			wantErr:     true,
		},
		{
			name:        "invalid cycle",
			price:       "15.99",
			currency:    "USD",
			cycle:       "weekly",
			paymentDate: "15-02-2025",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := subSvc.AddSubscription(tt.name, tt.price, tt.currency, tt.cycle, tt.paymentDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscriptionService_UpdateSubscription(t *testing.T) {
	subSvc, db, _ := setupSubscriptionService(t)

	sub := &database.Subscription{
		Name:        "Netflix",
		Price:       15.99,
		Currency:    "USD",
		Cycle:       "monthly",
		PaymentDate: time.Now().Add(5 * 24 * time.Hour),
	}

	if err := db.CreateSubscription(sub); err != nil {
		t.Fatalf("failed to create test subscription: %v", err)
	}

	tests := []struct {
		name        string
		id          uint
		price       string
		currency    string
		cycle       string
		paymentDate string
		wantErr     bool
	}{
		{
			name:        "update all fields",
			id:          sub.ID,
			price:       "19.99",
			currency:    "EUR",
			cycle:       "yearly",
			paymentDate: "20-03-2025",
			wantErr:     false,
		},
		{
			name:        "update name only",
			id:          sub.ID,
			price:       "",
			currency:    "",
			cycle:       "",
			paymentDate: "",
			wantErr:     false,
		},
		{
			name:    "invalid price",
			id:      sub.ID,
			price:   "invalid",
			wantErr: true,
		},
		{
			name:    "invalid cycle",
			id:      sub.ID,
			cycle:   "weekly",
			wantErr: true,
		},
		{
			name:    "non-existent ID",
			id:      999,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := subSvc.UpdateSubscription(tt.id, tt.name, tt.price, tt.currency, tt.cycle, tt.paymentDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSubscription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubscriptionService_DeleteSubscription(t *testing.T) {
	subSvc, db, _ := setupSubscriptionService(t)

	sub := &database.Subscription{
		Name:        "Netflix",
		Price:       15.99,
		Currency:    "USD",
		Cycle:       "monthly",
		PaymentDate: time.Now().Add(5 * 24 * time.Hour),
	}

	if err := db.CreateSubscription(sub); err != nil {
		t.Fatalf("failed to create test subscription: %v", err)
	}

	err := subSvc.DeleteSubscription(sub.ID)
	if err != nil {
		t.Errorf("DeleteSubscription() error = %v", err)
	}

	_, err = db.GetSubscriptionByID(sub.ID)
	if err == nil {
		t.Error("DeleteSubscription() subscription still exists")
	}
}

func TestSubscriptionService_ListSubscriptions(t *testing.T) {
	subSvc, db, _ := setupSubscriptionService(t)

	subs := []*database.Subscription{
		{
			Name:        "Netflix",
			Price:       15.99,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: time.Now().Add(5 * 24 * time.Hour),
		},
		{
			Name:        "Spotify",
			Price:       9.99,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: time.Now().Add(10 * 24 * time.Hour),
		},
	}

	for _, sub := range subs {
		if err := db.CreateSubscription(sub); err != nil {
			t.Fatalf("failed to create test subscription: %v", err)
		}
	}

	got, err := subSvc.ListSubscriptions()
	if err != nil {
		t.Fatalf("ListSubscriptions() error = %v", err)
	}

	if len(got) != len(subs) {
		t.Errorf("ListSubscriptions() returned %d subscriptions, want %d", len(got), len(subs))
	}
}

func TestSubscriptionService_CheckUpcomingPayments(t *testing.T) {
	subSvc, db, _ := setupSubscriptionService(t)

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
			Name:        "Due in 10 days",
			Price:       30.00,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: now.Add(10 * 24 * time.Hour),
		},
	}

	for _, sub := range subs {
		if err := db.CreateSubscription(sub); err != nil {
			t.Fatalf("failed to create test subscription: %v", err)
		}
	}

	got, err := subSvc.CheckUpcomingPayments()
	if err != nil {
		t.Fatalf("CheckUpcomingPayments() error = %v", err)
	}

	expected := 2
	if len(got) != expected {
		t.Errorf("CheckUpcomingPayments() returned %d subscriptions, want %d", len(got), expected)
	}
}

func TestSubscriptionService_SendNotifications(t *testing.T) {
	subSvc, db, mockTg := setupSubscriptionService(t)

	now := time.Now()

	subs := []database.Subscription{
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
	}

	for i := range subs {
		if err := db.CreateSubscription(&subs[i]); err != nil {
			t.Fatalf("failed to create test subscription: %v", err)
		}
	}

	notificationsSent := 0
	mockTg.SendNotificationFunc = func(name string, price float64, currency string, days int, cycle string, paymentDate string) error {
		notificationsSent++
		return nil
	}

	err := subSvc.SendNotifications(subs)
	if err != nil {
		t.Errorf("SendNotifications() error = %v", err)
	}

	if notificationsSent != 2 {
		t.Errorf("SendNotifications() sent %d notifications, want 2", notificationsSent)
	}
}

func TestSubscriptionService_UpdatePastDuePayments(t *testing.T) {
	subSvc, db, _ := setupSubscriptionService(t)

	now := time.Now()

	subs := []*database.Subscription{
		{
			Name:        "Past due 5 days",
			Price:       10.00,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: now.Add(-5 * 24 * time.Hour),
		},
		{
			Name:        "Past due 1 day",
			Price:       20.00,
			Currency:    "USD",
			Cycle:       "yearly",
			PaymentDate: now.Add(-1 * 24 * time.Hour),
		},
	}

	for _, sub := range subs {
		if err := db.CreateSubscription(sub); err != nil {
			t.Fatalf("failed to create test subscription: %v", err)
		}
	}

	err := subSvc.UpdatePastDuePayments()
	if err != nil {
		t.Errorf("UpdatePastDuePayments() error = %v", err)
	}

	updatedSubs, err := db.GetAllSubscriptions()
	if err != nil {
		t.Fatalf("failed to get updated subscriptions: %v", err)
	}

	for _, sub := range updatedSubs {
		if sub.PaymentDate.Before(now) {
			t.Errorf("UpdatePastDuePayments() payment date %v is still in the past", sub.PaymentDate)
		}
	}
}
