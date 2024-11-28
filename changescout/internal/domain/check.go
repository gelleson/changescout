package domain

import (
	"github.com/gelleson/changescout/changescout/internal/app/services/diff"
	"github.com/google/uuid"
	"time"
)

type Check struct {
	ID           uuid.UUID    `json:"id"`
	Cron         string       `json:"cron"`
	WebsiteID    uuid.UUID    `json:"website_id"`
	Result       []byte       `json:"result"`
	HasError     bool         `json:"has_error"`
	ErrorMessage string       `json:"error_message"`
	HasChanges   bool         `json:"has_diff"`
	DiffResult   *diff.Result `json:"diff_change"`
	CreatedAt    time.Time    `json:"created_at"`
}

type CheckResult struct {
	OldValue   []byte
	NewValue   []byte
	HasChanges bool
	Check      diff.Result
}
