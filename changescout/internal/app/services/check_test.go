package services_test

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/app/services"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/gelleson/changescout/changescout/internal/app/services/mocks"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CheckServiceSuite struct {
	suite.Suite
	checkService   *services.CheckService
	mockRepository *mocks.CheckRepository
}

func (suite *CheckServiceSuite) SetupTest() {
	suite.mockRepository = new(mocks.CheckRepository)
	suite.checkService = services.NewCheckService(suite.mockRepository)
}

func (suite *CheckServiceSuite) TestGetLatestCheckByWebsite() {
	websiteID := uuid.New()
	expectedCheck := domain.Check{ID: uuid.New()}

	suite.mockRepository.On("GetLatestCheckByWebsite", mock.Anything, websiteID).Return(expectedCheck, nil)

	check, err := suite.checkService.GetLatestCheckByWebsite(context.Background(), websiteID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedCheck, check)

	suite.mockRepository.AssertExpectations(suite.T())
}

func (suite *CheckServiceSuite) TestCreateCheck() {
	check := domain.Check{ID: uuid.New()}
	expectedCheck := domain.Check{ID: uuid.New()}

	suite.mockRepository.On("CreateCheck", mock.Anything, check).Return(expectedCheck, nil)

	resultCheck, err := suite.checkService.CreateCheck(context.Background(), check)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedCheck, resultCheck)

	suite.mockRepository.AssertExpectations(suite.T())
}

func (suite *CheckServiceSuite) TestClearChecksByWebsite() {
	websiteID := uuid.New()

	suite.mockRepository.On("ClearChecksByWebsite", mock.Anything, websiteID).Return(nil)

	err := suite.checkService.ClearChecksByWebsite(context.Background(), websiteID)
	assert.NoError(suite.T(), err)

	suite.mockRepository.AssertExpectations(suite.T())
}

func (suite *CheckServiceSuite) TestListChecks() {
	filters := database.CheckFilters{}
	pagination := domain.Pagination{}
	expectedChecks := []domain.Check{{ID: uuid.New()}}
	expectedCount := 10

	suite.mockRepository.On("ListChecks", mock.Anything, filters, pagination).Return(expectedChecks, expectedCount, nil)

	checks, count, err := suite.checkService.ListChecks(context.Background(), filters, pagination)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedChecks, checks)
	assert.Equal(suite.T(), expectedCount, count)

	suite.mockRepository.AssertExpectations(suite.T())
}

func (suite *CheckServiceSuite) TestUpdateCheck() {
	check := domain.Check{ID: uuid.New()}
	expectedCheck := domain.Check{ID: uuid.New()}

	suite.mockRepository.On("UpdateCheck", mock.Anything, check).Return(expectedCheck, nil)

	resultCheck, err := suite.checkService.UpdateCheck(context.Background(), check)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedCheck, resultCheck)

	suite.mockRepository.AssertExpectations(suite.T())
}

func (suite *CheckServiceSuite) TestGetCheckByID() {
	checkID := uuid.New()
	expectedCheck := domain.Check{ID: checkID}

	suite.mockRepository.On("GetCheckByID", mock.Anything, checkID).Return(expectedCheck, nil)

	check, err := suite.checkService.GetCheckByID(context.Background(), checkID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedCheck, check)

	suite.mockRepository.AssertExpectations(suite.T())
}

func TestCheckServiceSuite(t *testing.T) {
	suite.Run(t, new(CheckServiceSuite))
}
