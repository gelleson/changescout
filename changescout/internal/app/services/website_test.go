package services

import (
	"context"
	"errors"
	"github.com/gelleson/changescout/changescout/internal/app/services/mocks"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/gelleson/changescout/changescout/internal/utils/transform"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
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

// Тестирование метода Create
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

// Тестирование метода UpdateLastCheck
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

// Тестирование метода GetByID
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

// Тестирование метода List
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

// Добавленный тест для GetByURL
func (s *WebsiteServiceTestSuite) TestGetByURL() {
	tests := []struct {
		name    string
		url     string
		mock    func()
		wantErr bool
	}{
		{
			name: "successful get",
			url:  "https://example.com",
			mock: func() {
				expectedWebsite := domain.Website{
					ID:      uuid.New(),
					URL:     "https://example.com",
					UserID:  uuid.New(),
					Enabled: true,
					Cron:    "*/15 * * * *",
				}
				s.mockRepo.On("GetWebsiteByURL", s.ctx, "https://example.com").
					Return(expectedWebsite, nil).Once()
			},
			wantErr: false,
		},
		{
			name: "website not found",
			url:  "https://notfound.com",
			mock: func() {
				s.mockRepo.On("GetWebsiteByURL", s.ctx, "https://notfound.com").
					Return(domain.Website{}, errors.New("not found")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			got, err := s.service.GetByURL(s.ctx, tt.url)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.NotEmpty(got)
		})
	}
}

// Добавленный тест для GetByUserID
func (s *WebsiteServiceTestSuite) TestGetByUserID() {
	pagination := domain.Pagination{
		Limit:  10,
		Offset: 0,
	}
	tests := []struct {
		name    string
		userID  uuid.UUID
		mock    func()
		wantErr bool
	}{
		{
			name:   "successful get by user id",
			userID: uuid.New(),
			mock: func() {
				websites := []domain.Website{
					{
						ID:      uuid.New(),
						URL:     "https://example1.com",
						UserID:  uuid.New(),
						Enabled: true,
						Cron:    "*/15 * * * *",
					},
				}
				s.mockRepo.On("GetWebsiteByUserID", s.ctx, mock.AnythingOfType("uuid.UUID"), pagination).
					Return(websites, nil).Once()
			},
			wantErr: false,
		},
		{
			name:   "user has no websites",
			userID: uuid.New(),
			mock: func() {
				s.mockRepo.On("GetWebsiteByUserID", s.ctx, mock.AnythingOfType("uuid.UUID"), pagination).
					Return([]domain.Website{}, nil).Once()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			got, err := s.service.GetByUserID(s.ctx, tt.userID, pagination)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.NotNil(got)
		})
	}
}

// Добавленный тест для GetDueForCheck
func (s *WebsiteServiceTestSuite) TestGetDueForCheck() {
	pagination := domain.Pagination{
		Limit:  10,
		Offset: 0,
	}
	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "successful retrieve due websites",
			mock: func() {
				websites := []domain.Website{
					{
						ID:          uuid.New(),
						URL:         "https://exampledue.com",
						UserID:      uuid.New(),
						Enabled:     true,
						NextCheckAt: transform.ToPtr(time.Now().Add(-time.Hour)),
					},
				}
				s.mockRepo.On("GetWebsitesDueForCheck", s.ctx, pagination).
					Return(websites, nil).Once()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			got, err := s.service.GetDueForCheck(s.ctx, pagination)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.NotNil(got)
		})
	}
}

// Добавленный тест для GetByStatus
func (s *WebsiteServiceTestSuite) TestGetByStatus() {
	pagination := domain.Pagination{
		Limit:  10,
		Offset: 0,
	}
	tests := []struct {
		name    string
		enabled bool
		mock    func()
		wantErr bool
	}{
		{
			name:    "successful get by status",
			enabled: true,
			mock: func() {
				websites := []domain.Website{
					{
						ID:      uuid.New(),
						URL:     "https://examplestatus.com",
						UserID:  uuid.New(),
						Enabled: true,
						Cron:    "0 * * * *",
					},
				}
				s.mockRepo.On("GetWebsitesByStatus", s.ctx, true, pagination).
					Return(websites, nil).Once()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			got, err := s.service.GetByStatus(s.ctx, tt.enabled, pagination)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.NotNil(got)
		})
	}
}

// Добавленный тест для UpdateStatus
func (s *WebsiteServiceTestSuite) TestUpdateStatus() {
	tests := []struct {
		name    string
		id      uuid.UUID
		enabled bool
		mock    func()
		wantErr bool
	}{
		{
			name:    "successful status update",
			id:      uuid.New(),
			enabled: false,
			mock: func() {
				updatedWebsite := domain.Website{
					ID:      uuid.New(),
					URL:     "https://exampleupdatestatus.com",
					UserID:  uuid.New(),
					Enabled: false,
					Cron:    "0 * * * *",
				}
				s.mockRepo.On("UpdateStatusWebsite", s.ctx, mock.AnythingOfType("uuid.UUID"), false).
					Return(updatedWebsite, nil).Once()
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			got, err := s.service.UpdateStatus(s.ctx, tt.id, tt.enabled)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.Equal(tt.enabled, got.Enabled)
		})
	}
}

// Добавленный тест для Delete
func (s *WebsiteServiceTestSuite) TestDelete() {
	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		wantErr bool
	}{
		{
			name: "successful delete",
			id:   uuid.New(),
			mock: func() {
				s.mockRepo.On("DeleteWebsite", s.ctx, mock.AnythingOfType("uuid.UUID")).
					Return(nil).Once()
			},
			wantErr: false,
		},
		{
			name: "delete error not found",
			id:   uuid.New(),
			mock: func() {
				s.mockRepo.On("DeleteWebsite", s.ctx, mock.AnythingOfType("uuid.UUID")).
					Return(errors.New("not found")).Once()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			err := s.service.Delete(s.ctx, tt.id)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
		})
	}
}
