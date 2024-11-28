package ent

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/gelleson/changescout/changescout/internal/utils/testdb"
	"github.com/gelleson/changescout/changescout/internal/utils/transform"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type NotificationRepositoryTestSuite struct {
	suite.Suite
	repo *NotificationRepository
	ctx  context.Context
}

func TestNotificationRepository(t *testing.T) {
	suite.Run(t, new(NotificationRepositoryTestSuite))
}

func (s *NotificationRepositoryTestSuite) SetupTest() {
	client := testdb.NewEntClient()
	s.repo = NewNotificationRepository(client)
	s.ctx = context.Background()
}

func (s *NotificationRepositoryTestSuite) TestCreateNotification() {
	userID := uuid.New()
	websiteID := transform.ToPtr(uuid.New())
	token := "test-token"
	destination := "test-destination"

	tests := []struct {
		name         string
		notification domain.Notification
		wantErr      bool
	}{
		{
			name: "create notification with all fields",
			notification: domain.Notification{
				ID:          uuid.New(), // Add ID
				Type:        domain.TelegramNotificationType,
				Name:        "Full Test Notification",
				UserID:      userID,
				Token:       &token,
				Destination: &destination,
				WebsiteID:   websiteID,
				CreatedAt:   time.Now(),
			},
			wantErr: false,
		},
		{
			name: "create notification without optional fields",
			notification: domain.Notification{
				ID:        uuid.New(), // Add ID
				Type:      domain.TelegramNotificationType,
				Name:      "Basic Test Notification",
				UserID:    userID,
				WebsiteID: websiteID,
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.CreateNotification(s.ctx, tt.notification)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.NotEqual(s.T(), uuid.Nil, got.ID)
			// TODO: Fix this test
			//assert.Equal(s.T(), tt.notification.ID, got.ID) // Compare with the set ID
			assert.Equal(s.T(), tt.notification.Type, got.Type)
			assert.Equal(s.T(), tt.notification.Name, got.Name)
			assert.Equal(s.T(), tt.notification.UserID, got.UserID)
			assert.Equal(s.T(), tt.notification.WebsiteID, got.WebsiteID)
			assert.Equal(s.T(), tt.notification.Token, got.Token)
			assert.Equal(s.T(), tt.notification.Destination, got.Destination)
			assert.NotZero(s.T(), got.CreatedAt)
			assert.NotZero(s.T(), got.UpdatedAt)
		})
	}
}

func (s *NotificationRepositoryTestSuite) TestGetNotificationByID() {
	userID := uuid.New()
	websiteID := transform.ToPtr(uuid.New())
	token := "test-token"
	notification := domain.Notification{
		Type:      domain.TelegramNotificationType,
		Name:      "Test Notification By ID",
		UserID:    userID,
		Token:     &token,
		WebsiteID: websiteID,
		CreatedAt: time.Now(),
	}

	created, err := s.repo.CreateNotification(s.ctx, notification)
	assert.NoError(s.T(), err)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "existing notification",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "non-existing notification",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.GetNotificationByID(s.ctx, tt.id)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.Equal(s.T(), created.ID, got.ID)
			assert.Equal(s.T(), created.Name, got.Name)
			assert.Equal(s.T(), created.Type, got.Type)
			assert.Equal(s.T(), created.UserID, got.UserID)
			assert.Equal(s.T(), created.WebsiteID, got.WebsiteID)
		})
	}
}

func (s *NotificationRepositoryTestSuite) TestGetNotificationByWebsite() {
	userID := uuid.New()
	websiteID := transform.ToPtr(uuid.New())
	notification := domain.Notification{
		Type:      domain.TelegramNotificationType,
		Name:      "Test Notification By Website",
		UserID:    userID,
		WebsiteID: websiteID,
		CreatedAt: time.Now(),
	}

	created, err := s.repo.CreateNotification(s.ctx, notification)
	assert.NoError(s.T(), err)

	tests := []struct {
		name      string
		websiteID *uuid.UUID
		wantErr   bool
	}{
		{
			name:      "existing website",
			websiteID: created.WebsiteID,
			wantErr:   false,
		},
		{
			name:      "non-existing website",
			websiteID: transform.ToPtr(uuid.New()),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.GetNotificationByWebsite(s.ctx, *tt.websiteID)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.Equal(s.T(), created.WebsiteID, got.WebsiteID)
			assert.Equal(s.T(), created.Name, got.Name)
		})
	}
}

func (s *NotificationRepositoryTestSuite) TestListNotifications() {
	userID := uuid.New()
	websiteID := uuid.New()
	now := time.Now()

	notifications := []domain.Notification{
		{
			ID:        uuid.New(), // Add unique ID
			Type:      domain.TelegramNotificationType,
			Name:      "First Test Notification",
			UserID:    userID,
			WebsiteID: &websiteID,
			CreatedAt: now,
		},
		{
			ID:        uuid.New(), // Add unique ID
			Type:      domain.TelegramNotificationType,
			Name:      "Second Test Notification",
			UserID:    userID,
			WebsiteID: &websiteID,
			CreatedAt: now.Add(time.Hour),
		},
	}

	// Create the notifications and store them
	createdNotifications := make([]domain.Notification, 0, len(notifications))
	for _, n := range notifications {
		created, err := s.repo.CreateNotification(s.ctx, n)
		assert.NoError(s.T(), err)
		createdNotifications = append(createdNotifications, created)
	}

	tests := []struct {
		name       string
		filters    database.NotificationFilters
		pagination domain.Pagination
		wantCount  int
		wantErr    bool
	}{
		{
			name: "filter by user ID",
			filters: database.NotificationFilters{
				UserID: &userID,
			},
			pagination: domain.Pagination{
				Limit:  10,
				Offset: 0,
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "filter by website ID",
			filters: database.NotificationFilters{
				WebsiteID: &websiteID,
			},
			pagination: domain.Pagination{
				Limit:  10,
				Offset: 0,
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name: "filter by date range",
			filters: database.NotificationFilters{
				FromDate: &now,
				ToDate:   transform.ToPtr(now.Add(2 * time.Hour)), // Extend the time range to include both notifications
			},
			pagination: domain.Pagination{
				Limit:  10,
				Offset: 0,
			},
			wantCount: 2,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, count, err := s.repo.ListNotifications(s.ctx, tt.filters, tt.pagination)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tt.wantCount, count)
			assert.Len(s.T(), got, tt.wantCount)

			// Verify that all notifications are returned with correct data
			if count > 0 {
				for _, n := range got {
					assert.NotEqual(s.T(), uuid.Nil, n.ID)
					assert.NotEmpty(s.T(), n.Name)
					assert.Equal(s.T(), userID, n.UserID)
					assert.Equal(s.T(), &websiteID, n.WebsiteID)
				}
			}
		})
	}
}
func (s *NotificationRepositoryTestSuite) TestUpdateNotification() {
	userID := uuid.New()
	websiteID := uuid.New()
	token := "initial-token"
	notification := domain.Notification{
		Type:      domain.TelegramNotificationType,
		Name:      "Initial Test Notification",
		UserID:    userID,
		Token:     &token,
		WebsiteID: &websiteID,
		CreatedAt: time.Now(),
	}

	created, err := s.repo.CreateNotification(s.ctx, notification)
	assert.NoError(s.T(), err)

	newToken := "updated-token"
	tests := []struct {
		name         string
		notification domain.Notification
		wantErr      bool
		verify       func(assert.TestingT, domain.Notification)
	}{
		{
			name: "update token and name",
			notification: domain.Notification{
				ID:        created.ID,
				Type:      domain.TelegramNotificationType,
				Name:      "Updated Test Notification",
				Token:     &newToken,
				WebsiteID: created.WebsiteID,
			},
			wantErr: false,
			verify: func(t assert.TestingT, n domain.Notification) {
				assert.Equal(t, "Updated Test Notification", n.Name)
				assert.Equal(t, newToken, *n.Token)
			},
		},
		{
			name: "update non-existent notification",
			notification: domain.Notification{
				ID:   uuid.New(),
				Type: domain.TelegramNotificationType,
				Name: "Should Not Update",
			},
			wantErr: true,
			verify:  func(t assert.TestingT, n domain.Notification) {},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.UpdateNotification(s.ctx, tt.notification)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			tt.verify(s.T(), got)
		})
	}
}
