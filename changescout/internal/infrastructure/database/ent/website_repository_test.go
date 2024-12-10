package ent

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent"
	"github.com/gelleson/changescout/changescout/internal/utils/testdb"
	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

func getUserID(client *ent.Client) uuid.UUID {
	usr := domain.User{}
	if err := faker.FakeData(&usr); err != nil {
		panic(err)
	}
	usr.Role = domain.RoleUser

	user, err := NewUserRepository(client).CreateUser(context.Background(), usr)
	if err != nil {
		panic(err)
	}
	return user.ID
}

type WebsiteRepositoryTestSuite struct {
	suite.Suite
	repo   *WebsiteRepository
	ctx    context.Context
	userID uuid.UUID
}

func TestWebsiteRepository(t *testing.T) {
	suite.Run(t, new(WebsiteRepositoryTestSuite))
}

func (s *WebsiteRepositoryTestSuite) SetupTest() {
	client := testdb.NewEntClient()
	s.repo = NewWebsiteRepository(client)
	s.ctx = context.Background()
	s.userID = getUserID(client)
}

func (s *WebsiteRepositoryTestSuite) TestCreateWebsite() {
	tests := []struct {
		name    string
		website domain.Website
		wantErr bool
	}{
		{
			name: "successful creation",
			website: domain.Website{
				Name:    "Test Website",
				URL:     "https://example.com",
				Enabled: true,
				Mode:    domain.ModePlain,
				Cron:    "*/15 * * * *",
				Setting: domain.Setting{},
				UserID:  s.userID,
			},
			wantErr: false,
		},
		{
			name: "duplicate URL",
			website: domain.Website{
				Name:    "Test Website 2",
				URL:     "https://example.com", // Same URL as above
				Enabled: true,
				Mode:    domain.ModePlain,
				UserID:  s.userID,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.CreateWebsite(s.ctx, tt.website)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.NotEqual(s.T(), uuid.Nil, got.ID)
			assert.Equal(s.T(), tt.website.Name, got.Name)
			assert.Equal(s.T(), tt.website.URL, got.URL)
			assert.Equal(s.T(), tt.website.Enabled, got.Enabled)
			assert.Equal(s.T(), tt.website.Mode, got.Mode)
			assert.Equal(s.T(), tt.website.UserID, got.UserID)
		})
	}
}

func (s *WebsiteRepositoryTestSuite) TestGetWebsiteByID() {
	// Create test website
	website := domain.Website{
		Name:    "Test Website",
		URL:     "https://example.com",
		Enabled: true,
		Mode:    domain.ModePlain,
		UserID:  s.userID,
		Cron:    "*/15 * * * *",
	}

	created, err := s.repo.CreateWebsite(s.ctx, website)
	assert.NoError(s.T(), err)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "existing website",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "non-existing website",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.GetWebsiteByID(s.ctx, tt.id)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.Equal(s.T(), created.ID, got.ID)
			assert.Equal(s.T(), created.Name, got.Name)
			assert.Equal(s.T(), created.URL, got.URL)
		})
	}
}

func (s *WebsiteRepositoryTestSuite) TestGetWebsiteByURL() {
	website := domain.Website{
		Name:    "Test Website",
		URL:     "https://example.com",
		Enabled: true,
		Mode:    domain.ModePlain,
		UserID:  s.userID,
		Cron:    "*/15 * * * *",
	}

	created, err := s.repo.CreateWebsite(s.ctx, website)
	assert.NoError(s.T(), err)

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "existing URL",
			url:     created.URL,
			wantErr: false,
		},
		{
			name:    "non-existing URL",
			url:     "https://nonexistent.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.GetWebsiteByURL(s.ctx, tt.url)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.Equal(s.T(), created.URL, got.URL)
		})
	}
}

func (s *WebsiteRepositoryTestSuite) TestGetWebsiteByUserID() {
	otherUserID := getUserID(s.repo.client)
	websites := []domain.Website{
		{
			Name:    "Website 1",
			URL:     "https://example1.com",
			Enabled: true,
			Mode:    domain.ModePlain,
			UserID:  s.userID,
			Cron:    "*/15 * * * *",
		},
		{
			Name:    "Website 2",
			URL:     "https://example2.com",
			Enabled: true,
			Mode:    domain.ModePlain,
			UserID:  s.userID,
			Cron:    "*/15 * * * *",
		},
		{
			Name:    "Other User Website",
			URL:     "https://example3.com",
			Enabled: true,
			Mode:    domain.ModePlain,
			UserID:  otherUserID,
			Cron:    "*/15 * * * *",
		},
	}

	for _, w := range websites {
		_, err := s.repo.CreateWebsite(s.ctx, w)
		assert.NoError(s.T(), err)
	}

	tests := []struct {
		name       string
		userID     uuid.UUID
		pagination domain.Pagination
		wantCount  int
	}{
		{
			name:       "get user websites",
			userID:     s.userID,
			pagination: domain.Pagination{Limit: 10, Offset: 0},
			wantCount:  2,
		},
		{
			name:       "get other user websites",
			userID:     otherUserID,
			pagination: domain.Pagination{Limit: 10, Offset: 0},
			wantCount:  1,
		},
		{
			name:       "get with pagination",
			userID:     s.userID,
			pagination: domain.Pagination{Limit: 1, Offset: 0},
			wantCount:  1,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.GetWebsiteByUserID(s.ctx, tt.userID, tt.pagination)
			assert.NoError(s.T(), err)
			assert.Len(s.T(), got, tt.wantCount)
		})
	}
}

func (s *WebsiteRepositoryTestSuite) TestUpdateWebsite() {
	// Create initial website
	website := domain.Website{
		Name:    "Test Website",
		URL:     "https://example.com",
		Enabled: true,
		Mode:    domain.ModePlain,
		UserID:  s.userID,
		Cron:    "*/15 * * * *",
	}

	created, err := s.repo.CreateWebsite(s.ctx, website)
	assert.NoError(s.T(), err)

	// Update test cases
	tests := []struct {
		name    string
		update  domain.Website
		wantErr bool
	}{
		{
			name: "update name",
			update: domain.Website{
				ID:   created.ID,
				Name: "Updated Name",
			},
			wantErr: false,
		},
		{
			name: "update multiple fields",
			update: domain.Website{
				ID:          created.ID,
				Name:        "New Name",
				URL:         "https://updated.com",
				Enabled:     false,
				LastCheckAt: &time.Time{},
			},
			wantErr: false,
		},
		{
			name: "update non-existent website",
			update: domain.Website{
				ID:   uuid.New(),
				Name: "Should Fail",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			_, err := s.repo.UpdateWebsite(s.ctx, tt.update)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)

			// Verify the update
			updated, err := s.repo.GetWebsiteByID(s.ctx, tt.update.ID)
			assert.NoError(s.T(), err)

			if tt.update.Name != "" {
				assert.Equal(s.T(), tt.update.Name, updated.Name)
			}
			if tt.update.URL != "" {
				assert.Equal(s.T(), tt.update.URL, updated.URL)
			}
		})
	}
}

func (s *WebsiteRepositoryTestSuite) TestListWebsites() {
	// Create test websites
	websites := []domain.Website{
		{
			Name:    "Website 1",
			URL:     "https://example1.com",
			Enabled: true,
			Mode:    domain.ModePlain,
			UserID:  s.userID,
			Cron:    "*/15 * * * *",
		},
		{
			Name:    "Website 2",
			URL:     "https://example2.com",
			Enabled: false,
			Mode:    domain.ModePlain,
			UserID:  s.userID,
			Cron:    "*/15 * * * *",
		},
	}

	for _, w := range websites {
		_, err := s.repo.CreateWebsite(s.ctx, w)
		assert.NoError(s.T(), err)
	}

	tests := []struct {
		name       string
		filters    database.WebsiteFilters
		pagination domain.Pagination
		wantCount  int
	}{
		{
			name:       "no filters",
			filters:    database.WebsiteFilters{},
			pagination: domain.Pagination{Limit: 10, Offset: 0},
			wantCount:  2,
		},
		{
			name: "filter by enabled",
			filters: database.WebsiteFilters{
				Enabled: ptr(true),
			},
			pagination: domain.Pagination{Limit: 10, Offset: 0},
			wantCount:  1,
		},
		{
			name: "filter by user ID",
			filters: database.WebsiteFilters{
				UserID: &s.userID,
			},
			pagination: domain.Pagination{Limit: 10, Offset: 0},
			wantCount:  2,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, total, err := s.repo.ListWebsites(s.ctx, tt.filters, tt.pagination)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tt.wantCount, len(got))
			assert.Equal(s.T(), tt.wantCount, total)
		})
	}
}

func (s *WebsiteRepositoryTestSuite) TestGetWebsitesByStatus() {
	websites := []domain.Website{
		{
			Name:    "Enabled Website",
			URL:     "https://enabled.com",
			Enabled: true,
			Mode:    domain.ModePlain,
			UserID:  s.userID,
			Cron:    "*/15 * * * *",
		},
		{
			Name:    "Disabled Website",
			URL:     "https://disabled.com",
			Enabled: false,
			Mode:    domain.ModePlain,
			UserID:  s.userID,
			Cron:    "*/15 * * * *",
		},
	}

	for _, w := range websites {
		_, err := s.repo.CreateWebsite(s.ctx, w)
		assert.NoError(s.T(), err)
	}

	tests := []struct {
		name       string
		enabled    bool
		pagination domain.Pagination
		wantCount  int
	}{
		{
			name:       "get enabled websites",
			enabled:    true,
			pagination: domain.Pagination{Limit: 10, Offset: 0},
			wantCount:  1,
		},
		{
			name:       "get disabled websites",
			enabled:    false,
			pagination: domain.Pagination{Limit: 10, Offset: 0},
			wantCount:  1,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.GetWebsitesByStatus(s.ctx, tt.enabled, tt.pagination)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tt.wantCount, len(got))
		})
	}
}

func (s *WebsiteRepositoryTestSuite) TestSearchWebsites() {
	websites := []domain.Website{
		{
			Name:    "Test Website",
			URL:     "https://test.com",
			Enabled: true,
			Mode:    domain.ModePlain,
			UserID:  s.userID,
			Cron:    "*/15 * * * *",
		},
		{
			Name:    "Another Website",
			URL:     "https://example.com",
			Enabled: true,
			Mode:    domain.ModePlain,
			UserID:  s.userID,
			Cron:    "*/15 * * * *",
		},
	}

	for _, w := range websites {
		_, err := s.repo.CreateWebsite(s.ctx, w)
		assert.NoError(s.T(), err)
	}

	tests := []struct {
		name       string
		query      string
		pagination domain.Pagination
		wantCount  int
	}{
		{
			name:       "search by name",
			query:      "Test",
			pagination: domain.Pagination{Limit: 10, Offset: 0},
			wantCount:  1,
		},
		{
			name:       "search by URL",
			query:      "example",
			pagination: domain.Pagination{Limit: 10, Offset: 0},
			wantCount:  1,
		},
		{
			name:       "no results",
			query:      "nonexistent",
			pagination: domain.Pagination{Limit: 10, Offset: 0},
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.SearchWebsites(s.ctx, tt.query, tt.pagination)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tt.wantCount, len(got))
		})
	}
}

func (s *WebsiteRepositoryTestSuite) TestGetWebsitesDueForCheck() {
	// Use a fixed timestamp for more predictable testing
	baseTime := time.Now()

	websites := []domain.Website{
		{
			Name:        "Due Website 1",
			URL:         "https://due1.com",
			Enabled:     true,
			Mode:        domain.ModePlain,
			UserID:      s.userID,
			Cron:        "*/15 * * * *", // Every 15 minutes
			LastCheckAt: ptr(baseTime.Add(-20 * time.Minute)),
			NextCheckAt: ptr(baseTime.Add(-5 * time.Minute)), // Due 5 minutes ago
		},
		{
			Name:        "Due Website 2",
			URL:         "https://due2.com",
			Enabled:     true,
			Mode:        domain.ModePlain,
			UserID:      s.userID,
			Cron:        "0 * * * *", // Every hour
			LastCheckAt: ptr(baseTime.Add(-65 * time.Minute)),
			NextCheckAt: ptr(baseTime.Add(-5 * time.Minute)), // Due 5 minutes ago
		},
		{
			Name:        "Not Due Website",
			URL:         "https://notdue.com",
			Enabled:     true,
			Mode:        domain.ModePlain,
			UserID:      s.userID,
			Cron:        "0 */2 * * *", // Every 2 hours
			LastCheckAt: ptr(baseTime.Add(-30 * time.Minute)),
			NextCheckAt: ptr(baseTime.Add(30 * time.Minute)), // Due in 30 minutes
		},
		{
			Name:        "Disabled Website",
			URL:         "https://disabled.com",
			Enabled:     false,
			Mode:        domain.ModePlain,
			UserID:      s.userID,
			Cron:        "*/15 * * * *",
			LastCheckAt: ptr(baseTime.Add(-20 * time.Minute)),
			NextCheckAt: ptr(baseTime.Add(-5 * time.Minute)), // Even though due, it's disabled
		},
		{
			Name:    "Never Checked Website",
			URL:     "https://neverchecked.com",
			Enabled: true,
			Mode:    domain.ModePlain,
			UserID:  s.userID,
			Cron:    "*/5 * * * *", // Every 5 minutes
			// LastCheckAt is nil
			NextCheckAt: ptr(baseTime.Add(-1 * time.Minute)), // Due 1 minute ago
		},
	}

	// Create test websites
	for _, w := range websites {
		_, err := s.repo.CreateWebsite(s.ctx, w)
		assert.NoError(s.T(), err)
	}

	// Add a small delay to ensure all timestamps are properly set
	time.Sleep(100 * time.Millisecond)

	tests := []struct {
		name       string
		pagination domain.Pagination
		wantCount  int
		wantURLs   []string // Add expected URLs for more precise testing
	}{
		{
			name:       "get all due websites",
			pagination: domain.Pagination{Limit: 10, Offset: 0},
			wantCount:  3, // Due Website 1, Due Website 2, and Never Checked Website
			wantURLs: []string{
				"https://due1.com",
				"https://due2.com",
				"https://neverchecked.com",
			},
		},
		{
			name:       "get due websites with limit",
			pagination: domain.Pagination{Limit: 2, Offset: 0},
			wantCount:  2,
			wantURLs: []string{
				"https://due1.com",
				"https://due2.com",
			},
		},
		{
			name:       "get due websites with offset",
			pagination: domain.Pagination{Limit: 2, Offset: 1},
			wantCount:  2,
			wantURLs: []string{
				"https://due2.com",
				"https://neverchecked.com",
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.GetWebsitesDueForCheck(s.ctx, tt.pagination)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tt.wantCount, len(got))

			// Verify that all returned websites are actually due for check
			gotURLs := make([]string, len(got))
			for i, website := range got {
				gotURLs[i] = website.URL
				assert.True(s.T(), website.Enabled, "Website should be enabled")
				assert.NotNil(s.T(), website.NextCheckAt, "NextCheckAt should not be nil")
				assert.True(s.T(), website.NextCheckAt.Before(time.Now()),
					"Website should be due for check. Next check: %v",
					website.NextCheckAt)
			}

			// If specific URLs are expected, verify them
			if len(tt.wantURLs) > 0 {
				assert.ElementsMatch(s.T(), tt.wantURLs, gotURLs,
					"Returned websites should match expected URLs")
			}
		})
	}
}

func (s *WebsiteRepositoryTestSuite) TestDeleteWebsite() {
	now := time.Now()

	site := domain.Website{
		Name:        "Test Website",
		URL:         "https://example.com",
		Enabled:     true,
		Mode:        domain.ModePlain,
		UserID:      s.userID,
		Cron:        "*/15 * * * *",
		LastCheckAt: ptr(now.Add(-2 * time.Minute)), // Should be due (last check + interval < now)
	}

	created, err := s.repo.CreateWebsite(s.ctx, site)
	assert.NoError(s.T(), err)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "delete existing website",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "delete non-existing website",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repo.DeleteWebsite(s.ctx, tt.id)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)

			// Verify the delete
			_, err = s.repo.GetWebsiteByID(s.ctx, tt.id)
			assert.Error(s.T(), err)
		})
	}
}

// Helper function to create pointer
func ptr[T any](v T) *T {
	return &v
}
