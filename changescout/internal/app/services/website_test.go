package services

import (
	"context"
	"errors"
	"github.com/gelleson/changescout/changescout/internal/app/services/mocks"
	"github.com/gelleson/changescout/changescout/internal/utils/transform"
	"testing"
	"time"

	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WebsiteServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.WebsiteRepository
	service  *WebsiteService
	ctx      context.Context
}

func (s *WebsiteServiceTestSuite) SetupTest() {
	s.mockRepo = mocks.NewWebsiteRepository(s.T())
	s.service = NewWebsiteService(s.mockRepo)
	s.ctx = context.Background()
}

func TestWebsiteServiceSuite(t *testing.T) {
	suite.Run(t, new(WebsiteServiceTestSuite))
}

func (s *WebsiteServiceTestSuite) TestCreate() {
	baseTime := time.Now()
	tests := []struct {
		name    string
		website domain.Website
		mock    func()
		wantErr bool
	}{
		{
			name: "successful creation",
			website: domain.Website{
				URL:     "https://example.com",
				UserID:  uuid.New(),
				Mode:    domain.ModePlain,
				Enabled: true,
				Cron:    "*/15 * * * *",
			},
			mock: func() {
				s.mockRepo.On("CreateWebsite", s.ctx, mock.MatchedBy(func(w domain.Website) bool {
					// Verify NextCheckAt is set and is in the future
					return w.NextCheckAt != nil && w.NextCheckAt.After(baseTime)
				})).Return(domain.Website{
					ID:          uuid.New(),
					URL:         "https://example.com",
					UserID:      uuid.New(),
					Enabled:     true,
					Cron:        "*/15 * * * *",
					NextCheckAt: transform.ToPtr(baseTime.Add(15 * time.Minute)),
				}, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "creation failed",
			website: domain.Website{
				URL:     "https://example.com",
				UserID:  uuid.New(),
				Mode:    domain.ModePlain,
				Enabled: true,
				Cron:    "*/15 * * * *",
			},
			mock: func() {
				s.mockRepo.On("CreateWebsite", s.ctx, mock.AnythingOfType("domain.Website")).
					Return(domain.Website{}, errors.New("creation failed")).Once()
			},
			wantErr: true,
		},
		{
			name: "invalid url",
			website: domain.Website{
				URL:     "[]lol",
				UserID:  uuid.New(),
				Mode:    domain.ModePlain,
				Enabled: true,
				Cron:    "* * * * *",
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "invalid mode",
			website: domain.Website{
				URL:     "https://example.com",
				UserID:  uuid.New(),
				Mode:    domain.Mode("lol"),
				Enabled: true,
				Cron:    "* * * * *",
			},
			mock:    func() {},
			wantErr: true,
		},
		{
			name: "invalid cron expression",
			website: domain.Website{
				URL:     "https://example.com",
				UserID:  uuid.New(),
				Mode:    domain.ModePlain,
				Enabled: true,
				Cron:    "invalid",
			},
			mock:    func() {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			got, err := s.service.Create(s.ctx, tt.website)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.NotEmpty(got)
			if !tt.wantErr {
				s.NotNil(got.NextCheckAt)
				s.True(got.NextCheckAt.After(baseTime))
			}
		})
	}
}

func (s *WebsiteServiceTestSuite) TestUpdateLastCheck() {
	baseTime := time.Now()
	tests := []struct {
		name      string
		websiteID uuid.UUID
		mock      func()
		wantErr   bool
	}{
		{
			name:      "successful update",
			websiteID: uuid.New(),
			mock: func() {
				website := domain.Website{
					ID:      uuid.New(),
					URL:     "https://example.com",
					UserID:  uuid.New(),
					Enabled: true,
					Cron:    "*/15 * * * *",
				}
				s.mockRepo.On("GetWebsiteByID", s.ctx, mock.AnythingOfType("uuid.UUID")).
					Return(website, nil).Once()
				s.mockRepo.On("UpdateWebsite", s.ctx, mock.MatchedBy(func(w domain.Website) bool {
					return w.LastCheckAt != nil && w.NextCheckAt != nil &&
						w.LastCheckAt.After(baseTime.Add(-time.Second)) &&
						w.NextCheckAt.After(*w.LastCheckAt)
				})).Return(website, nil).Once()
			},
			wantErr: false,
		},
		{
			name:      "website not found",
			websiteID: uuid.New(),
			mock: func() {
				s.mockRepo.On("GetWebsiteByID", s.ctx, mock.AnythingOfType("uuid.UUID")).
					Return(domain.Website{}, errors.New("not found")).Once()
			},
			wantErr: true,
		},
		{
			name:      "invalid cron expression",
			websiteID: uuid.New(),
			mock: func() {
				website := domain.Website{
					ID:      uuid.New(),
					URL:     "https://example.com",
					UserID:  uuid.New(),
					Enabled: true,
					Cron:    "invalid",
				}
				s.mockRepo.On("GetWebsiteByID", s.ctx, mock.AnythingOfType("uuid.UUID")).
					Return(website, nil).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			err := s.service.UpdateLastCheck(s.ctx, tt.websiteID)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
		})
	}
}

// Rest of the test functions remain the same, but update the mock websites to include Cron field
func (s *WebsiteServiceTestSuite) TestGetByID() {
	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		wantErr bool
	}{
		{
			name: "successful get",
			id:   uuid.New(),
			mock: func() {
				expectedWebsite := domain.Website{
					ID:      uuid.New(),
					URL:     "https://example.com",
					UserID:  uuid.New(),
					Enabled: true,
					Cron:    "*/15 * * * *",
				}
				s.mockRepo.On("GetWebsiteByID", s.ctx, mock.AnythingOfType("uuid.UUID")).
					Return(expectedWebsite, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "website not found",
			id:   uuid.New(),
			mock: func() {
				s.mockRepo.On("GetWebsiteByID", s.ctx, mock.AnythingOfType("uuid.UUID")).
					Return(domain.Website{}, errors.New("not found")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			got, err := s.service.GetByID(s.ctx, tt.id)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.NotEmpty(got)
		})
	}
}

func (s *WebsiteServiceTestSuite) TestList() {
	pagination := domain.Pagination{
		Limit:  10,
		Offset: 1,
	}
	filters := database.WebsiteFilters{
		UserID: transform.ToPtr(uuid.New()),
	}

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "successful list",
			mock: func() {
				websites := []domain.Website{
					{
						ID:      uuid.New(),
						URL:     "https://example1.com",
						UserID:  uuid.New(),
						Enabled: true,
						Cron:    "*/15 * * * *",
					},
					{
						ID:      uuid.New(),
						URL:     "https://example2.com",
						UserID:  uuid.New(),
						Enabled: true,
						Cron:    "0 * * * *",
					},
				}
				s.mockRepo.On("ListWebsites", s.ctx, filters, pagination).
					Return(websites, len(websites), nil).Once()
			},
			wantErr: false,
		},
		{
			name: "list error",
			mock: func() {
				s.mockRepo.On("ListWebsites", s.ctx, filters, pagination).
					Return(nil, 0, errors.New("list error")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			got, count, err := s.service.List(s.ctx, filters, pagination)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.NotEmpty(got)
			s.Greater(count, 0)
		})
	}
}
