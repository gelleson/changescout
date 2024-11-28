package domain

import (
	"github.com/google/uuid"
	"time"
)

type NotificationType string

const (
	TelegramNotificationType NotificationType = "telegram"
)

func (t NotificationType) Values() []string {
	return []string{
		string(TelegramNotificationType),
	}
}

type Notification struct {
	ID          uuid.UUID        `json:"id"`
	Name        string           `json:"name"`
	Type        NotificationType `json:"type"`
	UserID      uuid.UUID        `json:"user_id"`
	Token       *string          `json:"token"`
	Destination *string          `json:"destination"`
	WebsiteID   *uuid.UUID       `json:"website_id"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
