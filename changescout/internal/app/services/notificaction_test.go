package services_test

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/gelleson/changescout/changescout/internal/app/services"
	"github.com/gelleson/changescout/changescout/internal/app/services/mocks"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type NotificationServiceSuite struct {
	suite.Suite
	service  *services.NotificationService
	mockRepo *mocks.NotificationRepository
}

func (suite *NotificationServiceSuite) SetupTest() {
	suite.mockRepo = new(mocks.NotificationRepository)
	suite.service = services.NewNotificationService(suite.mockRepo)
}

func (suite *NotificationServiceSuite) TestCreate() {
	notification := domain.Notification{ID: uuid.New()}
	suite.mockRepo.On("CreateNotification", mock.Anything, notification).Return(notification, nil)

	result, err := suite.service.Create(context.Background(), notification)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), notification, result)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *NotificationServiceSuite) TestGetByID() {
	id := uuid.New()
	expectedNotification := domain.Notification{ID: id}
	suite.mockRepo.On("GetNotificationByID", mock.Anything, id).Return(expectedNotification, nil)

	result, err := suite.service.GetByID(context.Background(), id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedNotification, result)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *NotificationServiceSuite) TestGetByWebsite() {
	websiteID := uuid.New()
	expectedNotification := domain.Notification{ID: uuid.New()}
	suite.mockRepo.On("GetNotificationByWebsite", mock.Anything, websiteID).Return(expectedNotification, nil)

	result, err := suite.service.GetByWebsite(context.Background(), websiteID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedNotification, result)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *NotificationServiceSuite) TestList() {
	filters := database.NotificationFilters{}
	pagination := domain.Pagination{}
	expectedNotifications := []domain.Notification{{ID: uuid.New()}}
	expectedCount := 5
	suite.mockRepo.On("ListNotifications", mock.Anything, filters, pagination).Return(expectedNotifications, expectedCount, nil)

	notifications, count, err := suite.service.List(context.Background(), filters, pagination)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedNotifications, notifications)
	assert.Equal(suite.T(), expectedCount, count)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *NotificationServiceSuite) TestUpdate() {
	notification := domain.Notification{ID: uuid.New()}
	updatedNotification := notification
	updatedNotification.Name = "Updated Name"
	suite.mockRepo.On("UpdateNotification", mock.Anything, notification).Return(updatedNotification, nil)

	result, err := suite.service.Update(context.Background(), notification)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedNotification, result)

	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *NotificationServiceSuite) TestDelete() {
	id := uuid.New()
	suite.mockRepo.On("DeleteNotification", mock.Anything, id).Return(nil)

	err := suite.service.Delete(context.Background(), id)
	assert.NoError(suite.T(), err)

	suite.mockRepo.AssertExpectations(suite.T())
}

func TestNotificationServiceSuite(t *testing.T) {
	suite.Run(t, new(NotificationServiceSuite))
}
