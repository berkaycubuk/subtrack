package database

import (
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *DB {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	return db
}

func TestNew(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if db == nil {
		t.Fatal("New() returned nil DB")
	}
}

func TestCreateSubscription(t *testing.T) {
	db := setupTestDB(t)

	sub := &Subscription{
		Name:        "Netflix",
		Price:       15.99,
		Currency:    "USD",
		Cycle:       "monthly",
		PaymentDate: time.Now().Add(5 * 24 * time.Hour),
	}

	err := db.CreateSubscription(sub)
	if err != nil {
		t.Fatalf("CreateSubscription() error = %v", err)
	}

	if sub.ID == 0 {
		t.Error("CreateSubscription() did not set ID")
	}
}

func TestGetSubscriptionByID(t *testing.T) {
	db := setupTestDB(t)

	sub := &Subscription{
		Name:        "Netflix",
		Price:       15.99,
		Currency:    "USD",
		Cycle:       "monthly",
		PaymentDate: time.Now().Add(5 * 24 * time.Hour),
	}

	if err := db.CreateSubscription(sub); err != nil {
		t.Fatalf("failed to create test subscription: %v", err)
	}

	got, err := db.GetSubscriptionByID(sub.ID)
	if err != nil {
		t.Fatalf("GetSubscriptionByID() error = %v", err)
	}

	if got == nil {
		t.Fatal("GetSubscriptionByID() returned nil")
	}

	if got.Name != sub.Name {
		t.Errorf("GetSubscriptionByID() Name = %v, want %v", got.Name, sub.Name)
	}

	if got.Price != sub.Price {
		t.Errorf("GetSubscriptionByID() Price = %v, want %v", got.Price, sub.Price)
	}
}

func TestGetSubscriptionByID_NotFound(t *testing.T) {
	db := setupTestDB(t)

	_, err := db.GetSubscriptionByID(999)
	if err == nil {
		t.Error("GetSubscriptionByID() expected error for non-existent ID")
	}
}

func TestGetAllSubscriptions(t *testing.T) {
	db := setupTestDB(t)

	subs := []*Subscription{
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

	got, err := db.GetAllSubscriptions()
	if err != nil {
		t.Fatalf("GetAllSubscriptions() error = %v", err)
	}

	if len(got) != len(subs) {
		t.Errorf("GetAllSubscriptions() returned %d subscriptions, want %d", len(got), len(subs))
	}
}

func TestUpdateSubscription(t *testing.T) {
	db := setupTestDB(t)

	sub := &Subscription{
		Name:        "Netflix",
		Price:       15.99,
		Currency:    "USD",
		Cycle:       "monthly",
		PaymentDate: time.Now().Add(5 * 24 * time.Hour),
	}

	if err := db.CreateSubscription(sub); err != nil {
		t.Fatalf("failed to create test subscription: %v", err)
	}

	sub.Name = "Netflix Premium"
	sub.Price = 19.99

	if err := db.UpdateSubscription(sub); err != nil {
		t.Fatalf("UpdateSubscription() error = %v", err)
	}

	updated, err := db.GetSubscriptionByID(sub.ID)
	if err != nil {
		t.Fatalf("failed to get updated subscription: %v", err)
	}

	if updated.Name != "Netflix Premium" {
		t.Errorf("UpdateSubscription() Name = %v, want Netflix Premium", updated.Name)
	}

	if updated.Price != 19.99 {
		t.Errorf("UpdateSubscription() Price = %v, want 19.99", updated.Price)
	}
}

func TestDeleteSubscription(t *testing.T) {
	db := setupTestDB(t)

	sub := &Subscription{
		Name:        "Netflix",
		Price:       15.99,
		Currency:    "USD",
		Cycle:       "monthly",
		PaymentDate: time.Now().Add(5 * 24 * time.Hour),
	}

	if err := db.CreateSubscription(sub); err != nil {
		t.Fatalf("failed to create test subscription: %v", err)
	}

	if err := db.DeleteSubscription(sub.ID); err != nil {
		t.Fatalf("DeleteSubscription() error = %v", err)
	}

	_, err := db.GetSubscriptionByID(sub.ID)
	if err == nil {
		t.Error("DeleteSubscription() subscription still exists after deletion")
	}
}

func TestGetUpcomingPayments(t *testing.T) {
	db := setupTestDB(t)

	now := time.Now()

	subs := []*Subscription{
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
			Name:        "Due in 6 days",
			Price:       30.00,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: now.Add(6 * 24 * time.Hour),
		},
		{
			Name:        "Past due",
			Price:       40.00,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: now.Add(-2 * 24 * time.Hour),
		},
	}

	for _, sub := range subs {
		if err := db.CreateSubscription(sub); err != nil {
			t.Fatalf("failed to create test subscription: %v", err)
		}
	}

	got, err := db.GetUpcomingPayments(5)
	if err != nil {
		t.Fatalf("GetUpcomingPayments() error = %v", err)
	}

	expected := 2
	if len(got) != expected {
		t.Errorf("GetUpcomingPayments() returned %d subscriptions, want %d", len(got), expected)
	}
}

func TestGetPastDuePayments(t *testing.T) {
	db := setupTestDB(t)

	now := time.Now()

	subs := []*Subscription{
		{
			Name:        "Past due 1 day",
			Price:       10.00,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: now.Add(-1 * 24 * time.Hour),
		},
		{
			Name:        "Past due 5 days",
			Price:       20.00,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: now.Add(-5 * 24 * time.Hour),
		},
		{
			Name:        "Future due",
			Price:       30.00,
			Currency:    "USD",
			Cycle:       "monthly",
			PaymentDate: now.Add(2 * 24 * time.Hour),
		},
	}

	for _, sub := range subs {
		if err := db.CreateSubscription(sub); err != nil {
			t.Fatalf("failed to create test subscription: %v", err)
		}
	}

	got, err := db.GetPastDuePayments()
	if err != nil {
		t.Fatalf("GetPastDuePayments() error = %v", err)
	}

	expected := 2
	if len(got) != expected {
		t.Errorf("GetPastDuePayments() returned %d subscriptions, want %d", len(got), expected)
	}
}
