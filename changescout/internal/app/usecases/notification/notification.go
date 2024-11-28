package notification

import (
	"context"
	"fmt"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/google/uuid"
	"strings"
	"time"
)

type Sender interface {
	Send(notification string, conf domain.Notification) error
}

type NotificationService interface {
	List(ctx context.Context, filters database.NotificationFilters, pagination domain.Pagination) ([]domain.Notification, int, error)
}

type WebsiteService interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.Website, error)
}

type UseCase struct {
	sender              Sender
	websiteService      WebsiteService
	notificationService NotificationService
}

func NewUseCase(sender Sender, websiteService WebsiteService, notificationService NotificationService) *UseCase {
	return &UseCase{sender: sender, websiteService: websiteService, notificationService: notificationService}
}

var ansi = map[string]string{
	"+": "\033[32m+\033[0m",
	"-": "\033[31m-\033[0m",
}

const template = `
*Website:* %s
*Mode*: %s

*Details:*
- URL: %s
- Last Checked: %s
- Diff:
%s
%s
%s
`

var diffWithColor = strings.NewReplacer(
	"+", ansi["+"],
	"-", ansi["-"],
)

func (c UseCase) NotifyChanges(ctx context.Context, siteID uuid.UUID, change domain.CheckResult) error {
	site, err := c.websiteService.GetByID(ctx, siteID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(template,
		site.Name,
		site.Mode,
		site.URL,
		time.Now().Format("2006-01-02 15:04:05"),
		"```diff",
		change.Check.Diff,
		"```",
	)

	senders, _, err := c.notificationService.List(ctx, database.NotificationFilters{
		WebsiteID: &siteID,
		UserID:    &site.UserID,
	}, domain.Pagination{})

	for _, conf := range senders {
		if err := c.sender.Send(msg, conf); err != nil {
			return err
		}
	}

	return nil
}
