package database

import (
	"fmt"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Subscription struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Price       float64   `gorm:"not null" json:"price"`
	Currency    string    `gorm:"not null" json:"currency"`
	Cycle       string    `gorm:"not null" json:"cycle"`
	PaymentDate time.Time `gorm:"not null" json:"payment_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DB struct {
	*gorm.DB
}

func New(dbPath string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&Subscription{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) CreateSubscription(sub *Subscription) error {
	return db.Create(sub).Error
}

func (db *DB) GetSubscriptionByID(id uint) (*Subscription, error) {
	var sub Subscription
	err := db.First(&sub, id).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (db *DB) GetAllSubscriptions() ([]Subscription, error) {
	var subs []Subscription
	err := db.Find(&subs).Error
	return subs, err
}

func (db *DB) UpdateSubscription(sub *Subscription) error {
	return db.Save(sub).Error
}

func (db *DB) DeleteSubscription(id uint) error {
	return db.Delete(&Subscription{}, id).Error
}

func (db *DB) GetUpcomingPayments(days int) ([]Subscription, error) {
	var subs []Subscription
	now := time.Now()
	cutoff := now.AddDate(0, 0, days)
	err := db.Where("payment_date >= ? AND payment_date <= ?", now, cutoff).Find(&subs).Error
	return subs, err
}

func (db *DB) GetPastDuePayments() ([]Subscription, error) {
	var subs []Subscription
	err := db.Where("payment_date < ?", time.Now()).Find(&subs).Error
	return subs, err
}
