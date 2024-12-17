package scheduler

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/scheduler/mocks"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/pkg/tests"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type UseCaseTestSuite struct {
	suite.Suite
	useCase        *UseCase
	mockPublisher  *mocks.Publisher
	mockWebService *mocks.WebsiteService
	testWebsites   []domain.Website
	ctx            context.Context
}

func TestUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UseCaseTestSuite))
}

func (s *UseCaseTestSuite) SetupTest() {
	s.mockPublisher = new(mocks.Publisher)
	s.mockWebService = new(mocks.WebsiteService)
	s.useCase = NewUseCase(s.mockPublisher, s.mockWebService, time.Microsecond)
	s.useCase.checkInterval = 100 * time.Millisecond
	ctx, _ := context.WithTimeout(context.Background(), 250*time.Millisecond)
	s.ctx = ctx
	s.useCase.logger = tests.TestLogger(s.T())

	s.testWebsites = []domain.Website{
		{
			ID:      uuid.New(),
			Name:    "Test Site",
			URL:     "https://example.com",
			Enabled: true,
			Mode:    domain.ModePlain,
			Cron:    "*/15 * * * *",
		},
	}
}

func (s *UseCaseTestSuite) TearDownTest() {
	s.mockPublisher.AssertExpectations(s.T())
	s.mockWebService.AssertExpectations(s.T())
}

func (s *UseCaseTestSuite) TestCheckDueWebsites_Success() {
	// Arrange
	s.mockWebService.On("GetDueForCheck", mock.Anything, domain.Pagination{}).
		Return(s.testWebsites, nil)

	s.mockPublisher.On("Publish", topic, mock.AnythingOfType("*message.Message")).
		Run(func(args mock.Arguments) {
			msgs := args.Get(1).(*message.Message)
			var website domain.Website
			err := json.Unmarshal(msgs.Payload, &website)
			s.NoError(err)
			s.Equal(s.testWebsites[0].ID, website.ID)
		}).
		Return(nil)

	// Act
	err := s.useCase.checkDueWebsites(s.ctx)

	// Assert
	s.NoError(err)
}

func (s *UseCaseTestSuite) TestCheckDueWebsites_WebServiceError() {
	// Arrange
	s.mockWebService.On("GetDueForCheck", mock.Anything, domain.Pagination{}).
		Return(nil, assert.AnError)

	// Act
	err := s.useCase.checkDueWebsites(s.ctx)

	// Assert
	s.Error(err)
}

func (s *UseCaseTestSuite) TestCheckDueWebsites_PublisherError() {
	// Arrange
	s.mockWebService.On("GetDueForCheck", mock.Anything, domain.Pagination{}).
		Return(s.testWebsites, nil)

	s.mockPublisher.On("Publish", topic, mock.AnythingOfType("*message.Message")).
		Return(assert.AnError)

	// Act
	err := s.useCase.checkDueWebsites(s.ctx)

	// Assert
	s.Error(err)
}

func (s *UseCaseTestSuite) TestRun_SuccessfulAndGracefulShutdown() {
	// Arrange
	s.mockWebService.On("GetDueForCheck", mock.Anything, domain.Pagination{}).
		Return([]domain.Website{}, nil).Maybe()

	s.mockPublisher.On("Publish", topic, mock.Anything).
		Return(nil).Maybe()

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	// Act
	err := s.useCase.Run(ctx)

	// Assert
	s.NoError(err)
}
