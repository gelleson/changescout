package usecases

import (
	"context"
	"errors"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/mocks"
	"testing"

	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/pkg/crons"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WebsiteUseCaseTestSuite struct {
	suite.Suite
	useCase            *WebsiteUseCase
	mockWebsiteService *mocks.WebsiteService
	mockUserService    *mocks.UserService
	ctx                context.Context
}

func (suite *WebsiteUseCaseTestSuite) SetupTest() {
	suite.mockWebsiteService = new(mocks.WebsiteService)
	suite.mockUserService = new(mocks.UserService)
	suite.useCase = NewWebsiteUseCase(suite.mockWebsiteService, suite.mockUserService)
	suite.ctx = context.TODO()
}

func (suite *WebsiteUseCaseTestSuite) TestCreate_Successful() {
	userID := uuid.New()
	website := domain.Website{URL: "http://example.com"}

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("Create", suite.ctx, mock.AnythingOfType("domain.Website")).Return(website, nil).Once()

	createdWebsite, err := suite.useCase.Create(suite.ctx, userID, website)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), website, createdWebsite)
}

func (suite *WebsiteUseCaseTestSuite) TestCreate_UserNotFound() {
	userID := uuid.New()
	website := domain.Website{}

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, errors.New("user not found")).Once()

	_, err := suite.useCase.Create(suite.ctx, userID, website)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), domain.ErrInvalidToken, err)
}

func (suite *WebsiteUseCaseTestSuite) TestGetByID_Successful() {
	userID := uuid.New()
	websiteID := uuid.New()
	expectedWebsite := domain.Website{URL: "http://example.com"}

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByID", suite.ctx, websiteID).Return(expectedWebsite, nil).Once()

	website, err := suite.useCase.GetByID(suite.ctx, userID, websiteID)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedWebsite, website)
}

func (suite *WebsiteUseCaseTestSuite) TestGetByID_UserNotFound() {
	userID := uuid.New()
	websiteID := uuid.New()

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, errors.New("user not found")).Once()

	_, err := suite.useCase.GetByID(suite.ctx, userID, websiteID)

	assert.Error(suite.T(), err)
}

func (suite *WebsiteUseCaseTestSuite) TestGetByURL_Successful() {
	userID := uuid.New()
	url := "http://example.com"
	expectedWebsite := domain.Website{URL: url}

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByURL", suite.ctx, url).Return(expectedWebsite, nil).Once()

	website, err := suite.useCase.GetByURL(suite.ctx, userID, url)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedWebsite, website)
}

func (suite *WebsiteUseCaseTestSuite) TestGetByURL_UserNotFound() {
	userID := uuid.New()
	url := "http://example.com"

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, errors.New("user not found")).Once()

	_, err := suite.useCase.GetByURL(suite.ctx, userID, url)

	assert.Error(suite.T(), err)
}

func (suite *WebsiteUseCaseTestSuite) TestGetByUserID_Successful() {
	userID := uuid.New()
	pagination := domain.Pagination{Limit: 10, Offset: 0}
	expectedWebsites := []domain.Website{{URL: "http://example.com"}}

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByUserID", suite.ctx, userID, pagination).Return(expectedWebsites, nil).Once()

	websites, err := suite.useCase.GetByUserID(suite.ctx, userID, pagination)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedWebsites, websites)
}

func (suite *WebsiteUseCaseTestSuite) TestGetByUserID_UserNotFound() {
	userID := uuid.New()
	pagination := domain.Pagination{Limit: 10, Offset: 0}

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, errors.New("user not found")).Once()

	_, err := suite.useCase.GetByUserID(suite.ctx, userID, pagination)

	assert.Error(suite.T(), err)
}

func (suite *WebsiteUseCaseTestSuite) TestUpdate_Successful() {
	userID := uuid.New()
	websiteID := uuid.New()
	website := domain.Website{URL: "http://example.com", UserID: userID, Cron: "* * * * *"}
	expectedWebsite := website

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByID", suite.ctx, websiteID).Return(expectedWebsite, nil).Once()
	suite.mockWebsiteService.On("Update", suite.ctx, mock.AnythingOfType("domain.Website")).Return(expectedWebsite, nil).Once()

	updatedWebsite, err := suite.useCase.Update(suite.ctx, userID, websiteID, website)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedWebsite, updatedWebsite)
}

func (suite *WebsiteUseCaseTestSuite) TestUpdate_UserNotAuthorized() {
	userID := uuid.New()
	anotherUserID := uuid.New()
	websiteID := uuid.New()
	website := domain.Website{URL: "http://example.com"}

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByID", suite.ctx, websiteID).Return(domain.Website{UserID: anotherUserID}, nil).Once()

	_, err := suite.useCase.Update(suite.ctx, userID, websiteID, website)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "not authorized", err.Error())
}

func (suite *WebsiteUseCaseTestSuite) TestUpdate_ValidateCron() {
	userID := uuid.New()
	websiteID := uuid.New()
	invalidCron := crons.CronExpression("invalid cron")
	website := domain.Website{URL: "http://example.com", UserID: userID, Cron: invalidCron}

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByID", suite.ctx, websiteID).Return(website, nil).Once()

	scheduler := crons.NewScheduler()
	err := scheduler.Validate(invalidCron)

	assert.Error(suite.T(), err)
}

func (suite *WebsiteUseCaseTestSuite) TestUpdateStatus_Successful() {
	userID := uuid.New()
	websiteID := uuid.New()
	enabled := true

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByID", suite.ctx, websiteID).Return(domain.Website{UserID: userID}, nil).Once()
	suite.mockWebsiteService.On("UpdateStatus", suite.ctx, websiteID, enabled).Return(domain.Website{UserID: userID, Enabled: enabled}, nil).Once()

	website, err := suite.useCase.UpdateStatus(suite.ctx, userID, websiteID, enabled)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), enabled, website.Enabled)
}

func (suite *WebsiteUseCaseTestSuite) TestUpdateStatus_UserNotAuthorized() {
	userID := uuid.New()
	anotherUserID := uuid.New()
	websiteID := uuid.New()
	enabled := true

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByID", suite.ctx, websiteID).Return(domain.Website{UserID: anotherUserID}, nil).Once()

	_, err := suite.useCase.UpdateStatus(suite.ctx, userID, websiteID, enabled)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "not authorized", err.Error())
}

func (suite *WebsiteUseCaseTestSuite) TestUpdateLastCheck() {
	websiteID := uuid.New()

	suite.mockWebsiteService.On("UpdateLastCheck", suite.ctx, websiteID).Return(nil).Once()

	err := suite.useCase.UpdateLastCheck(suite.ctx, websiteID)

	assert.NoError(suite.T(), err)
}

func (suite *WebsiteUseCaseTestSuite) TestDelete_Successful() {
	userID := uuid.New()
	websiteID := uuid.New()

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByID", suite.ctx, websiteID).Return(domain.Website{UserID: userID}, nil).Once()
	suite.mockWebsiteService.On("Delete", suite.ctx, websiteID).Return(nil).Once()

	err := suite.useCase.Delete(suite.ctx, userID, websiteID)

	assert.NoError(suite.T(), err)
}

func (suite *WebsiteUseCaseTestSuite) TestDelete_UserNotAuthorized() {
	userID := uuid.New()
	anotherUserID := uuid.New()
	websiteID := uuid.New()

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByID", suite.ctx, websiteID).Return(domain.Website{UserID: anotherUserID}, nil).Once()

	err := suite.useCase.Delete(suite.ctx, userID, websiteID)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "not authorized", err.Error())
}

func (suite *WebsiteUseCaseTestSuite) TestUpdateStatus_NotAuthorized() {
	userID := uuid.New()
	anotherUserID := uuid.New()
	websiteID := uuid.New()
	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByID", suite.ctx, websiteID).Return(domain.Website{UserID: anotherUserID}, nil).Once()

	_, err := suite.useCase.UpdateStatus(suite.ctx, userID, websiteID, true)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "not authorized", err.Error())
}

func (suite *WebsiteUseCaseTestSuite) TestDelete_Successful2() {
	userID := uuid.New()
	websiteID := uuid.New()

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByID", suite.ctx, websiteID).Return(domain.Website{UserID: userID}, nil).Once()
	suite.mockWebsiteService.On("Delete", suite.ctx, websiteID).Return(nil).Once()

	err := suite.useCase.Delete(suite.ctx, userID, websiteID)
	assert.NoError(suite.T(), err)
}

func (suite *WebsiteUseCaseTestSuite) TestDelete_NotAuthorized() {
	userID := uuid.New()
	anotherUserID := uuid.New()
	websiteID := uuid.New()

	suite.mockUserService.On("GetByID", suite.ctx, userID).Return(domain.User{}, nil).Once()
	suite.mockWebsiteService.On("GetByID", suite.ctx, websiteID).Return(domain.Website{UserID: anotherUserID}, nil).Once()

	err := suite.useCase.Delete(suite.ctx, userID, websiteID)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "not authorized", err.Error())
}

func TestWebsiteUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(WebsiteUseCaseTestSuite))
}
