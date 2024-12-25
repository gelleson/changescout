package notification

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/app/services/diff"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/gelleson/changescout/changescout/internal/pkg/templates"
	"github.com/gelleson/changescout/changescout/pkg/clock"
	"github.com/google/uuid"
	"strings"
	"text/template"
)

type data struct {
	Name        string
	Mode        domain.Mode
	URL         string
	LastChecked string
	Result      diff.Result
}

//go:generate mockery --name Sender
type Sender interface {
	Send(notification string, conf domain.Notification) error
}

//go:generate mockery --name NotificationService
type NotificationService interface {
	List(ctx context.Context, filters database.NotificationFilters, pagination domain.Pagination) ([]domain.Notification, int, error)
}

//go:generate mockery --name WebsiteService
type WebsiteService interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.Website, error)
}

type UseCase struct {
	sender              Sender
	websiteService      WebsiteService
	notificationService NotificationService
	now                 *clock.Clock
}

func NewUseCase(sender Sender, websiteService WebsiteService, notificationService NotificationService) *UseCase {
	return &UseCase{sender: sender, websiteService: websiteService, notificationService: notificationService, now: clock.New()}
}

var ansi = map[string]string{
	"+": "\033[32m+\033[0m",
	"-": "\033[31m-\033[0m",
}

var diffWithColor = strings.NewReplacer(
	"+", ansi["+"],
	"-", ansi["-"],
)

func (c UseCase) NotifyChanges(ctx context.Context, siteID uuid.UUID, change domain.CheckResult) error {
	site, err := c.websiteService.GetByID(ctx, siteID)
	if err != nil {
		return err
	}
	// Determine which template to use
	tmplString := templates.DiffDefaultMessage
	if site.Setting.Template != nil {
		tmplString = *site.Setting.Template
	}

	tmpl, err := template.New("notification").Parse(tmplString)
	if err != nil {
		return err
	}

	data := data{
		Name:        site.Name,
		Mode:        site.Mode,
		URL:         site.URL,
		LastChecked: c.now.Now().Format("2006-01-02 15:04:05"),
		Result:      change.Check,
	}

	var msg strings.Builder
	if err := tmpl.Execute(&msg, data); err != nil {
		return err
	}

	senders, _, err := c.notificationService.List(ctx, database.NotificationFilters{
		WebsiteID: &siteID,
		UserID:    &site.UserID,
	}, domain.Pagination{})

	text := msg.String()
	for _, conf := range senders {
		if err := c.sender.Send(text, conf); err != nil {
			return err
		}
	}

	return nil
}
