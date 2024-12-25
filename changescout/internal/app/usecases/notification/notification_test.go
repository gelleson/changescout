package notification

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gelleson/changescout/changescout/internal/app/services/diff"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/notification/mocks"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/pkg/clock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type NotificationTestSuite struct {
	suite.Suite
	useCase                 *UseCase
	mockSender              *mocks.Sender
	mockNotificationService *mocks.NotificationService
	mockWebsiteService      *mocks.WebsiteService
	fixedTime               time.Time
}

func (suite *NotificationTestSuite) SetupTest() {
	suite.mockSender = new(mocks.Sender)
	suite.mockNotificationService = new(mocks.NotificationService)
	suite.mockWebsiteService = new(mocks.WebsiteService)

	suite.fixedTime = time.Date(2024, time.December, 6, 23, 14, 57, 0, time.UTC)

	suite.useCase = NewUseCase(suite.mockSender, suite.mockWebsiteService, suite.mockNotificationService)
	suite.useCase.now = clock.NewFixedTime(suite.fixedTime)
}

func (suite *NotificationTestSuite) TestNotifyChangesDefaultTemplate() {
	siteID := uuid.New()
	changeResult := domain.CheckResult{
		Check: diff.Result{
			Diff:    "+ added line - removed line",
			Changes: []diff.Change{},
		},
	}

	site := domain.Website{
		ID:     siteID,
		Name:   "Example Site",
		Mode:   "Live",
		URL:    "http://example.com",
		UserID: uuid.New(),
		// Template is not set, using default
	}

	notifications := []domain.Notification{
		{ID: uuid.New()},
	}

	// Mocking the services
	suite.mockWebsiteService.On("GetByID", mock.Anything, siteID).Return(site, nil)
	suite.mockNotificationService.On("List", mock.Anything, mock.Anything, domain.Pagination{}).Return(notifications, 1, nil)

	// Setting expectation on the mock sender
	suite.mockSender.On("Send", mock.Anything, notifications[0]).Return(nil)

	err := suite.useCase.NotifyChanges(context.Background(), siteID, changeResult)
	suite.NoError(err)

	suite.mockWebsiteService.AssertExpectations(suite.T())
	suite.mockNotificationService.AssertExpectations(suite.T())
	suite.mockSender.AssertExpectations(suite.T())
}

func (suite *NotificationTestSuite) TestNotifyChangesCustomTemplate() {
	siteID := uuid.New()
	changeResult := domain.CheckResult{
		Check: diff.Result{
			Diff: "+ added line\n- removed line",
		},
	}

	customTemplate := "*Custom Website:* {{.Name}}\n*Mode*: {{.Mode}}\n*Details:*\n- Custom URL: {{.URL}}\n- Last Checked: {{.LastChecked}}\n"

	site := domain.Website{
		ID:     siteID,
		Name:   "Example Site",
		Mode:   "Live",
		URL:    "http://example.com",
		UserID: uuid.New(),
		Setting: domain.Setting{
			Template: &customTemplate,
		},
	}

	notifications := []domain.Notification{
		{ID: uuid.New()},
	}

	// Mocking the services
	suite.mockWebsiteService.On("GetByID", mock.Anything, siteID).Return(site, nil)
	suite.mockNotificationService.On("List", mock.Anything, mock.Anything, domain.Pagination{}).Return(notifications, 1, nil)

	// Expected message based on the custom template
	expectedMessage := fmt.Sprintf("*Custom Website:* %s\n*Mode*: %s\n*Details:*\n- Custom URL: %s\n- Last Checked: %s\n", site.Name, site.Mode, site.URL, suite.fixedTime.Format("2006-01-02 15:04:05"))

	// Setting expectation on the mock sender
	suite.mockSender.On("Send", expectedMessage, notifications[0]).Return(nil)

	err := suite.useCase.NotifyChanges(context.Background(), siteID, changeResult)
	suite.NoError(err)

	suite.mockWebsiteService.AssertExpectations(suite.T())
	suite.mockNotificationService.AssertExpectations(suite.T())
	suite.mockSender.AssertExpectations(suite.T())
}

func TestNotificationTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationTestSuite))
}
