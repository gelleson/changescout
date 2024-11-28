package ent

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/notification"
	"github.com/google/uuid"
)

type NotificationRepository struct {
	client *ent.Client
}

func NewNotificationRepository(client *ent.Client) *NotificationRepository {
	return &NotificationRepository{
		client: client,
	}
}

func (n *NotificationRepository) CreateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	created, err := n.client.Notification.Create().
		SetType(notification.Type).
		SetName(notification.Name).
		SetUserID(notification.UserID).
		SetNillableToken(notification.Token).
		SetNillableDestination(notification.Destination).
		SetNillableWebsiteID(notification.WebsiteID).
		SetCreatedAt(notification.CreatedAt).
		Save(ctx)

	if err != nil {
		return domain.Notification{}, err
	}

	// Convert string to *string for Token and Destination
	var token, destination *string
	if created.Token != "" {
		token = &created.Token
	}
	if created.Destination != "" {
		destination = &created.Destination
	}

	return domain.Notification{
		ID:          created.ID,
		Name:        created.Name,
		Type:        created.Type,
		UserID:      created.UserID,
		Token:       token,
		Destination: destination,
		WebsiteID:   created.WebsiteID,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}, nil
}

func (n *NotificationRepository) GetNotificationByID(ctx context.Context, id uuid.UUID) (domain.Notification, error) {
	notification, err := n.client.Notification.Get(ctx, id)
	if err != nil {
		return domain.Notification{}, err
	}

	var token, destination *string
	if notification.Token != "" {
		token = &notification.Token
	}
	if notification.Destination != "" {
		destination = &notification.Destination
	}

	return domain.Notification{
		ID:          notification.ID,
		Name:        notification.Name,
		Type:        domain.NotificationType(notification.Type),
		UserID:      notification.UserID,
		Token:       token,
		Destination: destination,
		WebsiteID:   notification.WebsiteID,
		CreatedAt:   notification.CreatedAt,
		UpdatedAt:   notification.UpdatedAt,
	}, nil
}

func (n *NotificationRepository) GetNotificationByWebsite(ctx context.Context, websiteID uuid.UUID) (domain.Notification, error) {
	notification, err := n.client.Notification.Query().
		Where(notification.WebsiteID(websiteID)).
		Only(ctx)

	if err != nil {
		return domain.Notification{}, err
	}

	var token, destination *string
	if notification.Token != "" {
		token = &notification.Token
	}
	if notification.Destination != "" {
		destination = &notification.Destination
	}

	return domain.Notification{
		ID:          notification.ID,
		Name:        notification.Name,
		Type:        domain.NotificationType(notification.Type),
		UserID:      notification.UserID,
		Token:       token,
		Destination: destination,
		WebsiteID:   notification.WebsiteID,
		CreatedAt:   notification.CreatedAt,
		UpdatedAt:   notification.UpdatedAt,
	}, nil
}

func (n *NotificationRepository) ListNotifications(ctx context.Context, filters database.NotificationFilters, pagination domain.Pagination) ([]domain.Notification, int, error) {
	query := n.client.Notification.Query()

	if filters.WebsiteID != nil && filters.UserID == nil {
		query = query.Where(notification.WebsiteID(*filters.WebsiteID))
	}
	if filters.FromDate != nil {
		query = query.Where(notification.CreatedAtGTE(*filters.FromDate))
	}

	if filters.ToDate != nil {
		query = query.Where(notification.CreatedAtLTE(*filters.ToDate))
	}
	if filters.UserID != nil && filters.WebsiteID == nil {
		query = query.Where(notification.UserID(*filters.UserID))
	}
	if filters.UserID != nil && filters.WebsiteID != nil {
		query = query.Where(
			notification.Or(
				notification.UserID(*filters.UserID),
				notification.WebsiteID(*filters.WebsiteID),
			),
		)
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	notifications, err := query.
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	result := make([]domain.Notification, len(notifications))
	for i, n := range notifications {
		var token, destination *string
		if n.Token != "" {
			token = &n.Token
		}
		if n.Destination != "" {
			destination = &n.Destination
		}

		result[i] = domain.Notification{
			ID:          n.ID,
			Name:        n.Name,
			Type:        domain.NotificationType(n.Type),
			UserID:      n.UserID,
			Token:       token,
			Destination: destination,
			WebsiteID:   n.WebsiteID,
			CreatedAt:   n.CreatedAt,
			UpdatedAt:   n.UpdatedAt,
		}
	}

	return result, total, nil
}

func (n *NotificationRepository) UpdateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	query := n.client.Notification.UpdateOneID(notification.ID)

	if notification.Name != "" {
		query = query.SetName(notification.Name)
	}

	if notification.Token != nil {
		query = query.SetNillableToken(notification.Token)
	}
	if notification.Destination != nil {
		query = query.SetNillableDestination(notification.Destination)
	}
	if notification.WebsiteID != nil {
		query = query.SetNillableWebsiteID(notification.WebsiteID)
	}
	if notification.UserID != uuid.Nil {
		query = query.SetUserID(notification.UserID)
	}
	if notification.Type != "" {
		query = query.SetType(notification.Type)
	}

	updated, err := query.Save(ctx)

	if err != nil {
		return domain.Notification{}, err
	}

	var token, destination *string
	if updated.Token != "" {
		token = &updated.Token
	}
	if updated.Destination != "" {
		destination = &updated.Destination
	}

	return domain.Notification{
		ID:          updated.ID,
		Name:        updated.Name,
		Type:        updated.Type,
		UserID:      updated.UserID,
		Token:       token,
		Destination: destination,
		WebsiteID:   updated.WebsiteID,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}, nil
}

func (n *NotificationRepository) DeleteNotification(ctx context.Context, id uuid.UUID) error {
	err := n.client.Notification.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
