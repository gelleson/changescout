package ent

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/gelleson/changescout/changescout/internal/utils/testdb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type CheckRepositoryTestSuite struct {
	suite.Suite
	checkRepo   *CheckRepository
	websiteRepo *WebsiteRepository
	userRepo    *UserRepository
	ctx         context.Context
	website     domain.Website // Store reference website
	user        domain.User    // Store reference user
}

func TestCheckRepository(t *testing.T) {
	suite.Run(t, new(CheckRepositoryTestSuite))
}

func (s *CheckRepositoryTestSuite) SetupTest() {
	client := testdb.NewEntClient()
	s.checkRepo = NewCheckRepository(client)
	s.websiteRepo = NewWebsiteRepository(client)
	s.userRepo = NewUserRepository(client)
	s.ctx = context.Background()

	// Create a test user
	user := domain.User{
		Email:     "test@example.com",
		Role:      domain.RoleUser,
		IsActive:  true,
		Password:  "password",
		FirstName: "Test",
		LastName:  "User",
	}

	usr, err := s.userRepo.CreateUser(s.ctx, user)
	if err != nil {
		s.T().Fatal(err)
	}
	s.user = usr

	// Create a test website
	website := domain.Website{
		Name:    "Test Website",
		URL:     "https://example.com",
		Enabled: true,
		Mode:    domain.ModePlain,
		UserID:  s.user.ID,
		Cron:    "*/15 * * * *",
	}

	created, err := s.websiteRepo.CreateWebsite(s.ctx, website)
	if err != nil {
		s.T().Fatal(err)
	}
	s.website = created
}

func (s *CheckRepositoryTestSuite) TestCreateCheck() {
	tests := []struct {
		name    string
		check   domain.Check
		wantErr bool
	}{
		{
			name: "create valid check",
			check: domain.Check{
				WebsiteID: s.website.ID,
				Result:    []byte("test result"),
			},
			wantErr: false,
		},
		{
			name: "create check with non-existent website",
			check: domain.Check{
				WebsiteID: uuid.New(),
				Result:    []byte("test result"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.checkRepo.CreateCheck(s.ctx, tt.check)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.NotEqual(s.T(), uuid.Nil, got.ID)
			assert.Equal(s.T(), tt.check.WebsiteID, got.WebsiteID)
			assert.Equal(s.T(), tt.check.Result, got.Result)
			assert.NotZero(s.T(), got.CreatedAt)
		})
	}
}

func (s *CheckRepositoryTestSuite) TestGetCheckByID() {
	// Create test check
	check := domain.Check{
		WebsiteID: s.website.ID,
		Result:    []byte("test result"),
	}

	created, err := s.checkRepo.CreateCheck(s.ctx, check)
	assert.NoError(s.T(), err)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "existing check",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "non-existing check",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.checkRepo.GetCheckByID(s.ctx, tt.id)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.Equal(s.T(), created.ID, got.ID)
			assert.Equal(s.T(), created.WebsiteID, got.WebsiteID)
			assert.Equal(s.T(), created.Result, got.Result)
		})
	}
}

func (s *CheckRepositoryTestSuite) TestListChecks() {
	// Create another website for testing filters
	otherWebsite, err := s.websiteRepo.CreateWebsite(s.ctx, domain.Website{
		Name:    "Other Website",
		URL:     "https://other-example.com",
		Enabled: true,
		Mode:    domain.ModePlain,
		Cron:    "*/15 * * * *",
		UserID:  s.user.ID,
	})
	assert.NoError(s.T(), err)

	// Create multiple checks
	checks := []domain.Check{
		{
			WebsiteID: s.website.ID,
			Result:    []byte("result 1"),
		},
		{
			WebsiteID: s.website.ID,
			Result:    []byte("result 2"),
		},
		{
			WebsiteID: otherWebsite.ID,
			Result:    []byte("result 3"),
		},
	}

	for _, c := range checks {
		_, err := s.checkRepo.CreateCheck(s.ctx, c)
		assert.NoError(s.T(), err)
	}

	tests := []struct {
		name       string
		filters    database.CheckFilters
		pagination domain.Pagination
		wantCount  int
	}{
		{
			name: "filter by website ID",
			filters: database.CheckFilters{
				WebsiteID: &s.website.ID,
			},
			pagination: domain.Pagination{
				Limit:  10,
				Offset: 0,
			},
			wantCount: 2,
		},
		{
			name: "filter by other website ID",
			filters: database.CheckFilters{
				WebsiteID: &otherWebsite.ID,
			},
			pagination: domain.Pagination{
				Limit:  10,
				Offset: 0,
			},
			wantCount: 1,
		},
		{
			name: "filter by date range",
			filters: database.CheckFilters{
				FromDate: &time.Time{},
				ToDate:   &time.Time{},
			},
			pagination: domain.Pagination{
				Limit:  10,
				Offset: 0,
			},
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, count, err := s.checkRepo.ListChecks(s.ctx, tt.filters, tt.pagination)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tt.wantCount, count)
			assert.Len(s.T(), got, tt.wantCount)
		})
	}
}

func (s *CheckRepositoryTestSuite) TestUpdateCheck() {
	// Create initial check
	check := domain.Check{
		WebsiteID: s.website.ID,
		Result:    []byte("initial result"),
	}

	created, err := s.checkRepo.CreateCheck(s.ctx, check)
	assert.NoError(s.T(), err)

	tests := []struct {
		name    string
		check   domain.Check
		wantErr bool
	}{
		{
			name: "update existing check",
			check: domain.Check{
				ID:        created.ID,
				WebsiteID: created.WebsiteID,
				Result:    []byte("updated result"),
			},
			wantErr: false,
		},
		{
			name: "update non-existent check",
			check: domain.Check{
				ID:        uuid.New(),
				WebsiteID: uuid.New(),
				Result:    []byte("should fail"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.checkRepo.UpdateCheck(s.ctx, tt.check)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.Equal(s.T(), tt.check.ID, got.ID)
			assert.Equal(s.T(), tt.check.Result, got.Result)
		})
	}
}

func (s *CheckRepositoryTestSuite) TestClearChecksByWebsite() {
	// Create multiple checks for the website
	checks := []domain.Check{
		{
			WebsiteID: s.website.ID,
			Result:    []byte("result 1"),
		},
		{
			WebsiteID: s.website.ID,
			Result:    []byte("result 2"),
		},
	}

	for _, c := range checks {
		_, err := s.checkRepo.CreateCheck(s.ctx, c)
		assert.NoError(s.T(), err)
	}

	// Clear checks
	err := s.checkRepo.ClearChecksByWebsite(s.ctx, s.website.ID)
	assert.NoError(s.T(), err)

	// Verify checks were cleared
	got, count, err := s.checkRepo.ListChecks(s.ctx, database.CheckFilters{
		WebsiteID: &s.website.ID,
	}, domain.Pagination{Limit: 10})
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, count)
	assert.Empty(s.T(), got)
}

func (s *CheckRepositoryTestSuite) TestGetLatestCheckByWebsite() {
	// Create multiple checks
	checks := []domain.Check{
		{
			WebsiteID: s.website.ID,
			Result:    []byte("old result"),
		},
		{
			WebsiteID: s.website.ID,
			Result:    []byte("latest result"),
		},
	}

	var latest domain.Check
	for _, c := range checks {
		created, err := s.checkRepo.CreateCheck(s.ctx, c)
		assert.NoError(s.T(), err)
		latest = created
		time.Sleep(time.Millisecond * 100) // Ensure different timestamps
	}

	got, err := s.checkRepo.GetLatestCheckByWebsite(s.ctx, s.website.ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), latest.ID, got.ID)
	assert.Equal(s.T(), latest.Result, got.Result)
}
