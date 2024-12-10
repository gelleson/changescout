package database

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/google/uuid"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetTotalUsers(ctx context.Context) (int, error)
	UpdateUser(ctx context.Context, user domain.User) error
}

type WebsiteRepository interface {
	CreateWebsite(ctx context.Context, website domain.Website) (domain.Website, error)
	GetWebsiteByID(ctx context.Context, id uuid.UUID) (domain.Website, error)
	GetWebsiteByURL(ctx context.Context, url string) (domain.Website, error)
	GetWebsiteByUserID(ctx context.Context, userID uuid.UUID, pagination domain.Pagination) ([]domain.Website, error)
	UpdateWebsite(ctx context.Context, website domain.Website) (domain.Website, error)

	ListWebsites(ctx context.Context, filters WebsiteFilters, pagination domain.Pagination) ([]domain.Website, int, error)
	GetWebsitesByStatus(ctx context.Context, enabled bool, pagination domain.Pagination) ([]domain.Website, error)
	SearchWebsites(ctx context.Context, query string, pagination domain.Pagination) ([]domain.Website, error)
	GetWebsitesDueForCheck(ctx context.Context, pagination domain.Pagination) ([]domain.Website, error)
	UpdateStatusWebsite(ctx context.Context, siteID uuid.UUID, enabled bool) (domain.Website, error)

	DeleteWebsite(ctx context.Context, id uuid.UUID) error
}

type CheckRepository interface {
	CreateCheck(ctx context.Context, check domain.Check) (domain.Check, error)
	GetCheckByID(ctx context.Context, id uuid.UUID) (domain.Check, error)
	GetLatestCheckByWebsite(ctx context.Context, websiteID uuid.UUID) (domain.Check, error)
	ClearChecksByWebsite(ctx context.Context, websiteID uuid.UUID) error
	ListChecks(ctx context.Context, filters CheckFilters, pagination domain.Pagination) ([]domain.Check, int, error)
	UpdateCheck(ctx context.Context, check domain.Check) (domain.Check, error)
}

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error)
	GetNotificationByID(ctx context.Context, id uuid.UUID) (domain.Notification, error)
	GetNotificationByWebsite(ctx context.Context, websiteID uuid.UUID) (domain.Notification, error)
	ListNotifications(ctx context.Context, filters NotificationFilters, pagination domain.Pagination) ([]domain.Notification, int, error)
	UpdateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error)
	DeleteNotification(ctx context.Context, id uuid.UUID) error
}

// WebsiteFilters represents the filtering options for websites
type WebsiteFilters struct {
	UserID   *uuid.UUID
	Enabled  *bool
	Mode     *domain.Mode
	URLQuery *string
	FromDate *time.Time
	ToDate   *time.Time
}

// CheckFilters represents the filtering options for checks
type CheckFilters struct {
	WebsiteID *uuid.UUID
	FromDate  *time.Time
	ToDate    *time.Time
}

type NotificationFilters struct {
	WebsiteID *uuid.UUID
	UserID    *uuid.UUID
	FromDate  *time.Time
	ToDate    *time.Time
}
