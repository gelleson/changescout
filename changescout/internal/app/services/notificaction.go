package services

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/google/uuid"
)

//go:generate mockery --name NotificationRepository
type NotificationRepository interface {
	database.NotificationRepository
}

type NotificationService struct {
	repo NotificationRepository
}

func NewNotificationService(repo NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (n *NotificationService) Create(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	return n.repo.CreateNotification(ctx, notification)
}

func (n *NotificationService) GetByID(ctx context.Context, id uuid.UUID) (domain.Notification, error) {
	return n.repo.GetNotificationByID(ctx, id)
}

func (n *NotificationService) GetByWebsite(ctx context.Context, websiteID uuid.UUID) (domain.Notification, error) {
	return n.repo.GetNotificationByWebsite(ctx, websiteID)
}

func (n *NotificationService) List(ctx context.Context, filters database.NotificationFilters, pagination domain.Pagination) ([]domain.Notification, int, error) {
	return n.repo.ListNotifications(ctx, filters, pagination)
}

func (n *NotificationService) Update(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	return n.repo.UpdateNotification(ctx, notification)
}

func (n *NotificationService) Delete(ctx context.Context, id uuid.UUID) error {
	return n.repo.DeleteNotification(ctx, id)
}
