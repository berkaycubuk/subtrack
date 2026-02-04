package utils

import (
	"fmt"
	"math"
	"time"
)

const dateFormat = "02-01-2006"

func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(dateFormat, dateStr)
}

func FormatDate(t time.Time) string {
	return t.Format(dateFormat)
}

func DaysUntil(t time.Time) int {
	now := time.Now()
	duration := t.Sub(now)
	hours := duration.Hours()
	if hours < 0 {
		return -1
	}
	days := int(math.Round(hours / 24))
	return days
}

func UpdatePaymentDate(paymentDate time.Time, cycle string) (time.Time, error) {
	switch cycle {
	case "monthly":
		return paymentDate.AddDate(0, 1, 0), nil
	case "yearly":
		return paymentDate.AddDate(1, 0, 0), nil
	default:
		return time.Time{}, fmt.Errorf("invalid cycle: %s", cycle)
	}
}

func CalculateNextPaymentDate(lastPaymentDate time.Time, cycle string) (time.Time, error) {
	nextDate := lastPaymentDate

	for nextDate.Before(time.Now()) {
		var err error
		nextDate, err = UpdatePaymentDate(nextDate, cycle)
		if err != nil {
			return time.Time{}, err
		}
	}

	return nextDate, nil
}
