package scheduler

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/platform/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

const topic = "websites.check"

//go:generate mockery --name Publisher
type Publisher interface {
	message.Publisher
}

//go:generate mockery --name WebsiteService
type WebsiteService interface {
	GetDueForCheck(ctx context.Context, pagination domain.Pagination) ([]domain.Website, error)
}

type UseCase struct {
	publisher      Publisher
	websiteService WebsiteService
	checkInterval  time.Duration
	logger         *zap.Logger
}

func NewUseCase(publisher Publisher, websiteService WebsiteService, checkInterval time.Duration) *UseCase {
	return &UseCase{
		publisher:      publisher,
		websiteService: websiteService,
		logger:         logger.L("scheduler"),
		checkInterval:  checkInterval,
	}
}

func (u UseCase) checkDueWebsites(ctx context.Context) error {
	websites, err := u.websiteService.GetDueForCheck(ctx, domain.Pagination{})
	if err != nil {
		return err
	}

	messages := make([]*message.Message, len(websites))
	for index, site := range websites {
		payload, err := json.Marshal(site)
		if err != nil {
			return err
		}
		messages[index] = &message.Message{
			UUID: uuid.New().String(),
			Metadata: message.Metadata{
				"website_id": site.ID.String(),
			},
			Payload: payload,
		}
	}

	u.logger.Debug("Publishing messages", zap.Int("count", len(messages)))

	if err := u.publisher.Publish(topic, messages...); err != nil {
		return err
	}

	return nil
}

func (u UseCase) Run(ctx context.Context) error {
	ticker := time.NewTicker(u.checkInterval)
	defer ticker.Stop()
	u.logger.Info("Starting scheduler")

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return nil
		case <-ticker.C:
			u.logger.Debug("Checking due websites")
			if err := u.checkDueWebsites(ctx); err != nil {
				u.logger.Error("Error checking due websites", zap.Error(err))
			}
			u.logger.Debug("Checking due websites complete")
		}
	}
}
