package diffnotification

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/google/uuid"
)

type NotificationService interface {
	List(ctx context.Context, filters database.NotificationFilters, pagination domain.Pagination) ([]domain.Notification, int, error)
}

type Sender interface {
	Send(notification string, conf domain.Notification) error
}

type UseCase struct {
	service *NotificationService
	sender  Sender
}

func NewUseCase(service *NotificationService) *UseCase {
	return &UseCase{
		service: service,
	}
}

func (u *UseCase) Send(ctx context.Context, websiteID uuid.UUID) error {
	return nil
}
