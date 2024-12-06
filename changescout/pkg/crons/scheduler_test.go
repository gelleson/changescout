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
		{"valid every minute", "* * * * *", false},
		{"valid every 15 minutes", "*/15 * * * *", false},
		{"valid business hours", "0 9-17 * * MON-FRI", false},
		{"invalid fields count", "* * *", true},
		{"invalid minute", "60 * * * *", true},
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
		{"every minute", scheduler.EveryMinute(), baseTime, baseTime.Add(time.Minute)},
		{"every 15 minutes", scheduler.EveryNMinutes(15), baseTime, time.Date(2024, 1, 1, 12, 15, 0, 0, time.UTC)},
		{"daily", scheduler.Daily(), baseTime, time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)},
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
			"next 3 minutes", scheduler.EveryMinute(), baseTime, 3,
			[]time.Time{
				baseTime.Add(1 * time.Minute),
				baseTime.Add(2 * time.Minute),
				baseTime.Add(3 * time.Minute),
			},
		},
		{
			"next 2 hourly runs", scheduler.EveryHour(), baseTime, 2,
			[]time.Time{
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

	t.Run("invalid n", func(t *testing.T) {
		_, err := scheduler.NextNRuns(scheduler.EveryMinute(), baseTime, 0)
		assert.Error(t, err)
	})
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
		{"not overdue", scheduler.EveryNMinutes(15), baseTime, baseTime.Add(14 * time.Minute), false},
		{"overdue", scheduler.EveryNMinutes(15), baseTime, baseTime.Add(16 * time.Minute), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isOverdue, err := scheduler.IsOverdue(tt.expr, tt.lastRun, tt.now)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, isOverdue)
		})
	}
}

func TestScheduler_CommonCronExpressions(t *testing.T) {
	scheduler := NewScheduler()

	tests := []struct {
		name string
		expr CronExpression
		want string
	}{
		{"every minute", scheduler.EveryMinute(), "* * * * *"},
		{"every 5 minutes", scheduler.EveryNMinutes(5), "*/5 * * * *"},
		{"every hour", scheduler.EveryHour(), "0 * * * *"},
		{"every 2 hours", scheduler.EveryNHours(2), "0 */2 * * *"},
		{"daily", scheduler.Daily(), "0 0 * * *"},
		{"weekly", scheduler.Weekly(), "0 0 * * 0"},
		{"monthly", scheduler.Monthly(), "0 0 1 * *"},
		{"from short interval", scheduler.FromInterval(30 * time.Second), "*/1 * * * *"},
		{"from interval", scheduler.FromInterval(5 * time.Minute), "*/5 * * * *"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, string(tt.expr))
		})
	}
}

func TestScheduler_ParseAndValidate(t *testing.T) {
	scheduler := NewScheduler()

	validExpression := CronExpression("0 0 * * *")
	invalidExpression := CronExpression("60 24 *")

	t.Run("parse valid expression", func(t *testing.T) {
		schedule, err := scheduler.Parse(validExpression)
		require.NoError(t, err)
		assert.NotNil(t, schedule)
	})

	t.Run("parse invalid expression", func(t *testing.T) {
		schedule, err := scheduler.Parse(invalidExpression)
		assert.Error(t, err)
		assert.Nil(t, schedule)
	})

	t.Run("validate valid expression", func(t *testing.T) {
		err := scheduler.Validate(validExpression)
		assert.NoError(t, err)
	})

	t.Run("validate invalid expression", func(t *testing.T) {
		err := scheduler.Validate(invalidExpression)
		assert.Error(t, err)
	})
}
