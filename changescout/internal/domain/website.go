package domain

import (
	"github.com/gelleson/changescout/changescout/pkg/crons"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Mode string

const (
	ModePlain    Mode = "plain"
	ModeRenderer Mode = "renderer"
)

// Setting represents the options for a website. It can be used to configure the HTTP request, extract text from the response, and handle errors.
type Setting struct {
	Headers   http.Header `json:"headers"`
	UserAgent string      `json:"user_agent"`
	Referer   string      `json:"referer"`
	Method    string      `json:"method"`
	// Template is a Go template to render notifications
	Template *string `json:"template"`
	// RenderedOption is setting for the rendered mode
	RenderedOption RenderedOption `json:"rendered_option"`

	// Selectors is a list of CSS selectors to extract text from the HTML content or xpath expressions to extract text from the XML content.
	Selectors []string `json:"selectors"`
	// Deduplication is a boolean flag to enable or disable deduplication of websites.
	Deduplication bool `json:"deduplication"`
	// Sort alphabetically
	Sort bool `json:"sort"`
	// Trim whitespace
	Trim bool `json:"trim"`
	// JSONPath is a list of JSONPath expressions to extract text from the JSON content.
	JSONPath []string `json:"json_path"`
}

// Website represents a website to be monitored.
type Website struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	URL         string               `json:"url" binding:"required,url"`
	Enabled     bool                 `json:"enabled"`
	Mode        Mode                 `json:"mode"`
	Cron        crons.CronExpression `json:"cron"`
	Setting     Setting              `json:"setting"`
	UserID      uuid.UUID            `json:"user_id"` // Foreign key to User
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	LastCheckAt *time.Time           `json:"last_check_at"`
	NextCheckAt *time.Time           `json:"next_check_at"`
}

// RenderedOption represents settings for the rendered mode.
// The struct is used to configure options like waiting for selectors to appear and timeout intervals.
type RenderedOption struct {
	// WaitForTimeout specifies how long to wait for the selector (in seconds).
	// This is beneficial to set a time frame after which the wait will timeout if the selector doesn’t appear.
	WaitForTimeout *int `json:"wait_for_timeout"`
}
