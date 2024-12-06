package crons

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"time"
)

type CronExpression string

type Schedule interface {
	Next(t time.Time) time.Time
}

type Scheduler struct {
	parser cron.Parser
}

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

func (s *Scheduler) Parse(expr CronExpression) (Schedule, error) {
	return s.parser.Parse(string(expr))
}

func (s *Scheduler) Validate(expr CronExpression) error {
	_, err := s.parser.Parse(string(expr))
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}
	return nil
}

func (s *Scheduler) NextRun(expr CronExpression, from time.Time) (time.Time, error) {
	schedule, err := s.parser.Parse(string(expr))
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse cron expression: %w", err)
	}

	return schedule.Next(from), nil
}

func (s *Scheduler) NextNRuns(expr CronExpression, from time.Time, n int) ([]time.Time, error) {
	if n < 1 {
		return nil, fmt.Errorf("n must be greater than 0")
	}

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

func (s *Scheduler) IsOverdue(expr CronExpression, lastRun, now time.Time) (bool, error) {
	nextAfterLast, err := s.NextRun(expr, lastRun)
	if err != nil {
		return false, err
	}
	return now.After(nextAfterLast), nil
}

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
