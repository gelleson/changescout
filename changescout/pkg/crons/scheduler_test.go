package crons

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestScheduler_Validate(t *testing.T) {
	scheduler := NewScheduler()

	tests := []struct {
		name    string
		expr    CronExpression
		wantErr bool
	}{
		{
			name:    "valid every minute",
			expr:    "* * * * *",
			wantErr: false,
		},
		{
			name:    "valid every 15 minutes",
			expr:    "*/15 * * * *",
			wantErr: false,
		},
		{
			name:    "valid business hours",
			expr:    "0 9-17 * * MON-FRI",
			wantErr: false,
		},
		{
			name:    "invalid fields count",
			expr:    "* * *",
			wantErr: true,
		},
		{
			name:    "invalid minute",
			expr:    "60 * * * *",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := scheduler.Validate(tt.expr)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestScheduler_NextRun(t *testing.T) {
	scheduler := NewScheduler()
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		expr     CronExpression
		from     time.Time
		expected time.Time
	}{
		{
			name:     "every minute",
			expr:     scheduler.EveryMinute(),
			from:     baseTime,
			expected: baseTime.Add(time.Minute),
		},
		{
			name:     "every 15 minutes",
			expr:     scheduler.EveryNMinutes(15),
			from:     baseTime,
			expected: time.Date(2024, 1, 1, 12, 15, 0, 0, time.UTC),
		},
		{
			name:     "daily",
			expr:     scheduler.Daily(),
			from:     baseTime,
			expected: time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next, err := scheduler.NextRun(tt.expr, tt.from)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, next)
		})
	}
}

func TestScheduler_NextNRuns(t *testing.T) {
	scheduler := NewScheduler()
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		expr     CronExpression
		from     time.Time
		n        int
		expected []time.Time
	}{
		{
			name: "next 3 minutes",
			expr: scheduler.EveryMinute(),
			from: baseTime,
			n:    3,
			expected: []time.Time{
				baseTime.Add(1 * time.Minute),
				baseTime.Add(2 * time.Minute),
				baseTime.Add(3 * time.Minute),
			},
		},
		{
			name: "next 2 hourly runs",
			expr: scheduler.EveryHour(),
			from: baseTime,
			n:    2,
			expected: []time.Time{
				time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC),
				time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			times, err := scheduler.NextNRuns(tt.expr, tt.from, tt.n)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, times)
		})
	}
}

func TestScheduler_IsOverdue(t *testing.T) {
	scheduler := NewScheduler()
	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		expr     CronExpression
		lastRun  time.Time
		now      time.Time
		expected bool
	}{
		{
			name:     "not overdue",
			expr:     scheduler.EveryNMinutes(15),
			lastRun:  baseTime,
			now:      baseTime.Add(14 * time.Minute),
			expected: false,
		},
		{
			name:     "overdue",
			expr:     scheduler.EveryNMinutes(15),
			lastRun:  baseTime,
			now:      baseTime.Add(16 * time.Minute),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isOverdue, err := scheduler.IsOverdue(tt.expr, tt.lastRun, tt.now)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, isOverdue)
		})
	}
}

func TestScheduler_CommonExpressions(t *testing.T) {
	scheduler := NewScheduler()

	tests := []struct {
		name     string
		expr     CronExpression
		validate func(t *testing.T, expr CronExpression)
	}{
		{
			name: "every minute",
			expr: scheduler.EveryMinute(),
			validate: func(t *testing.T, expr CronExpression) {
				assert.Equal(t, CronExpression("* * * * *"), expr)
			},
		},
		{
			name: "every hour",
			expr: scheduler.EveryHour(),
			validate: func(t *testing.T, expr CronExpression) {
				assert.Equal(t, CronExpression("0 * * * *"), expr)
			},
		},
		{
			name: "daily",
			expr: scheduler.Daily(),
			validate: func(t *testing.T, expr CronExpression) {
				assert.Equal(t, CronExpression("0 0 * * *"), expr)
			},
		},
		{
			name: "weekly",
			expr: scheduler.Weekly(),
			validate: func(t *testing.T, expr CronExpression) {
				assert.Equal(t, CronExpression("0 0 * * 0"), expr)
			},
		},
		{
			name: "monthly",
			expr: scheduler.Monthly(),
			validate: func(t *testing.T, expr CronExpression) {
				assert.Equal(t, CronExpression("0 0 1 * *"), expr)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, tt.expr)
			assert.NoError(t, scheduler.Validate(tt.expr))
		})
	}
}
