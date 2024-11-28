// internal/app/crons/scheduler.go
package crons

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

// CronExpression represents a cron expression with validation and parsing capabilities
type CronExpression string

// Schedule defines the interface for cron scheduling operations
type Schedule interface {
	Next(t time.Time) time.Time
}

// Scheduler handles cron expression parsing and next execution time calculations
type Scheduler struct {
	parser cron.Parser
}

// NewScheduler creates a new Scheduler instance
func NewScheduler() *Scheduler {
	parser := cron.NewParser(
		cron.SecondOptional |
			cron.Minute |
			cron.Hour |
			cron.Dom |
			cron.Month |
			cron.Dow,
	)

	return &Scheduler{
		parser: parser,
	}
}

// Parse parses a cron expression and returns a Schedule
func (s *Scheduler) Parse(expr CronExpression) (Schedule, error) {
	return s.parser.Parse(string(expr))
}

// Validate checks if a cron expression is valid
func (s *Scheduler) Validate(expr CronExpression) error {
	_, err := s.parser.Parse(string(expr))
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}
	return nil
}

// NextRun calculates the next execution time based on the cron expression
func (s *Scheduler) NextRun(expr CronExpression, from time.Time) (time.Time, error) {
	schedule, err := s.parser.Parse(string(expr))
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse cron expression: %w", err)
	}

	return schedule.Next(from), nil
}

// NextNRuns returns the next n execution times
func (s *Scheduler) NextNRuns(expr CronExpression, from time.Time, n int) ([]time.Time, error) {
	schedule, err := s.parser.Parse(string(expr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse cron expression: %w", err)
	}

	times := make([]time.Time, n)
	next := from

	for i := 0; i < n; i++ {
		next = schedule.Next(next)
		times[i] = next
	}

	return times, nil
}

// IsOverdue checks if the next execution time has passed
func (s *Scheduler) IsOverdue(expr CronExpression, lastRun, now time.Time) (bool, error) {
	nextAfterLast, err := s.NextRun(expr, lastRun)
	if err != nil {
		return false, err
	}
	return now.After(nextAfterLast), nil
}

// Common cron expression helpers
func (s *Scheduler) EveryMinute() CronExpression {
	return "* * * * *"
}

func (s *Scheduler) EveryNMinutes(n int) CronExpression {
	return CronExpression(fmt.Sprintf("*/%d * * * *", n))
}

func (s *Scheduler) EveryHour() CronExpression {
	return "0 * * * *"
}

func (s *Scheduler) EveryNHours(n int) CronExpression {
	return CronExpression(fmt.Sprintf("0 */%d * * *", n))
}

func (s *Scheduler) Daily() CronExpression {
	return "0 0 * * *"
}

func (s *Scheduler) Weekly() CronExpression {
	return "0 0 * * 0"
}

func (s *Scheduler) Monthly() CronExpression {
	return "0 0 1 * *"
}

func (s *Scheduler) FromInterval(duration time.Duration) CronExpression {
	minutes := int(duration.Minutes())
	if minutes < 1 {
		minutes = 1
	}
	return CronExpression(fmt.Sprintf("*/%d * * * *", minutes))
}
