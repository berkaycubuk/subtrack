package utils

import (
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:  "valid date",
			input: "15-02-2025",
			want:  time.Date(2025, time.February, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "valid date 31-12-2025",
			input: "31-12-2025",
			want:  time.Date(2025, time.December, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "invalid format",
			input:   "2025-02-15",
			wantErr: true,
		},
		{
			name:    "invalid date",
			input:   "31-02-2025",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ParseDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name  string
		input time.Time
		want  string
	}{
		{
			name:  "format date",
			input: time.Date(2025, time.February, 15, 0, 0, 0, 0, time.UTC),
			want:  "15-02-2025",
		},
		{
			name:  "format date with time",
			input: time.Date(2025, time.February, 15, 14, 30, 45, 0, time.UTC),
			want:  "15-02-2025",
		},
		{
			name:  "single digit day and month",
			input: time.Date(2025, time.January, 5, 0, 0, 0, 0, time.UTC),
			want:  "05-01-2025",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDate(tt.input)
			if got != tt.want {
				t.Errorf("FormatDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaysUntil(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name string
		t    time.Time
		want int
	}{
		{
			name: "future date - same day",
			t:    now.Add(2 * time.Hour),
			want: 0,
		},
		{
			name: "future date - 1 day",
			t:    now.Add(24 * time.Hour),
			want: 1,
		},
		{
			name: "future date - 5 days",
			t:    now.Add(5 * 24 * time.Hour),
			want: 5,
		},
		{
			name: "past date",
			t:    now.Add(-24 * time.Hour),
			want: -1,
		},
		{
			name: "past date - few hours",
			t:    now.Add(-2 * time.Hour),
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DaysUntil(tt.t)
			if got != tt.want {
				t.Errorf("DaysUntil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdatePaymentDate(t *testing.T) {
	baseDate := time.Date(2025, time.February, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		date    time.Time
		cycle   string
		want    time.Time
		wantErr bool
	}{
		{
			name:  "monthly cycle",
			date:  baseDate,
			cycle: "monthly",
			want:  time.Date(2025, time.March, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "yearly cycle",
			date:  baseDate,
			cycle: "yearly",
			want:  time.Date(2026, time.February, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "invalid cycle",
			date:    baseDate,
			cycle:   "weekly",
			wantErr: true,
		},
		{
			name:    "empty cycle",
			date:    baseDate,
			cycle:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdatePaymentDate(tt.date, tt.cycle)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePaymentDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("UpdatePaymentDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateNextPaymentDate(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		date    time.Time
		cycle   string
		wantErr bool
	}{
		{
			name:  "future date - no update needed",
			date:  now.Add(30 * 24 * time.Hour),
			cycle: "monthly",
		},
		{
			name:  "past date - single monthly update",
			date:  now.Add(-15 * 24 * time.Hour),
			cycle: "monthly",
		},
		{
			name:  "past date - single yearly update",
			date:  now.Add(-400 * 24 * time.Hour),
			cycle: "yearly",
		},
		{
			name:    "past date - invalid cycle",
			date:    now.Add(-15 * 24 * time.Hour),
			cycle:   "weekly",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateNextPaymentDate(tt.date, tt.cycle)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateNextPaymentDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Before(now) {
					t.Errorf("CalculateNextPaymentDate() = %v, should be after now", got)
				}
			}
		})
	}
}
