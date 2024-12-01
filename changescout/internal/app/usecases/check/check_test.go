package check

import (
	"context"
	"fmt"
	"github.com/gelleson/changescout/changescout/internal/app/services/diff"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/check/mocks"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type CheckTestSuite struct {
	suite.Suite
	useCase        *UseCase
	websiteService *mocks.WebsiteService
	httpService    *mocks.HttpService
	checkService   *mocks.DBService
	diffService    *mocks.DiffService
	ctx            context.Context
}

func (s *CheckTestSuite) SetupTest() {
	s.websiteService = mocks.NewWebsiteService(s.T())
	s.httpService = mocks.NewHttpService(s.T())
	s.checkService = mocks.NewDBService(s.T())
	s.diffService = mocks.NewDiffService(s.T())
	s.ctx = context.Background()

	s.useCase = NewUseCase(
		s.websiteService,
		s.httpService,
		s.checkService,
		s.diffService,
	)
}

func TestCheckSuite(t *testing.T) {
	suite.Run(t, new(CheckTestSuite))
}

func (s *CheckTestSuite) TestCheckNoChanges() {
	// Arrange
	websiteID := uuid.New()
	website := domain.Website{
		ID:  websiteID,
		URL: "https://example.com",
	}
	previousCheck := domain.Check{
		ID:        uuid.New(),
		WebsiteID: websiteID,
		Result:    []byte("previous content"),
	}
	currentContent := []byte("current content")
	diffResult := diff.Result{HasChanges: false}

	s.websiteService.On("GetByID", s.ctx, websiteID).Return(website, nil)
	s.httpService.On("Request", website).Return(currentContent, nil)
	s.checkService.On("GetLatestCheckByWebsite", s.ctx, websiteID).Return(previousCheck, nil)
	s.diffService.On("Compare", previousCheck.Result, currentContent).Return(diffResult, nil)

	// Act
	result, err := s.useCase.Check(s.ctx, websiteID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), result)
	s.websiteService.AssertExpectations(s.T())
	s.httpService.AssertExpectations(s.T())
	s.checkService.AssertExpectations(s.T())
	s.diffService.AssertExpectations(s.T())
}

func (s *CheckTestSuite) TestCheckWithChanges() {
	// Arrange
	websiteID := uuid.New()
	website := domain.Website{
		ID:  websiteID,
		URL: "https://example.com",
	}
	previousContent := []byte("previous content")
	previousCheck := domain.Check{
		ID:        uuid.New(),
		WebsiteID: websiteID,
		Result:    previousContent,
	}
	currentContent := []byte("updated content")
	diffResult := diff.Result{
		HasChanges: true,
		Changes: []diff.Change{
			{
				Type:    diff.Modified,
				Content: "updated content",
			},
		},
	}

	s.websiteService.On("GetByID", s.ctx, websiteID).Return(website, nil)
	s.httpService.On("Request", website).Return(currentContent, nil)
	s.checkService.On("GetLatestCheckByWebsite", s.ctx, websiteID).Return(previousCheck, nil)
	s.diffService.On("Compare", previousContent, currentContent).Return(diffResult, nil)
	s.checkService.On("CreateCheck", s.ctx, mock.MatchedBy(func(check domain.Check) bool {
		return check.WebsiteID == websiteID &&
			string(check.Result) == string(currentContent) &&
			check.HasChanges == true &&
			check.HasError == false
	})).Return(domain.Check{}, nil)

	// Act
	result, err := s.useCase.Check(s.ctx, websiteID)

	// Assert
	assert.NoError(s.T(), err)
	expectedResult := domain.CheckResult{
		OldValue:   previousContent,
		NewValue:   currentContent,
		HasChanges: true,
		Check:      diffResult,
	}
	assert.Equal(s.T(), expectedResult, result)
	s.websiteService.AssertExpectations(s.T())
	s.httpService.AssertExpectations(s.T())
	s.checkService.AssertExpectations(s.T())
	s.diffService.AssertExpectations(s.T())
}

func (s *CheckTestSuite) TestWebsiteNotFound() {
	// Arrange
	websiteID := uuid.New()
	expectedErr := domain.ErrWebsiteNotFound

	s.websiteService.On("GetByID", s.ctx, websiteID).Return(domain.Website{}, expectedErr)

	// Act
	result, err := s.useCase.Check(s.ctx, websiteID)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to get website")
	assert.ErrorIs(s.T(), err, expectedErr)
	assert.Empty(s.T(), result)
	s.websiteService.AssertExpectations(s.T())
}

func (s *CheckTestSuite) TestRequestError() {
	// Arrange
	websiteID := uuid.New()
	website := domain.Website{
		ID:  websiteID,
		URL: "https://example.com",
	}
	requestErr := domain.ErrRequestFailed

	s.websiteService.On("GetByID", s.ctx, websiteID).Return(website, nil)
	s.httpService.On("Request", website).Return(nil, requestErr)
	s.checkService.On("CreateCheck", s.ctx, mock.MatchedBy(func(check domain.Check) bool {
		return check.WebsiteID == websiteID &&
			check.HasError == true &&
			check.ErrorMessage == requestErr.Error()
	})).Return(domain.Check{}, nil)

	// Act
	result, err := s.useCase.Check(s.ctx, websiteID)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to make HTTP request")
	assert.ErrorIs(s.T(), err, requestErr)
	assert.Empty(s.T(), result)
	s.websiteService.AssertExpectations(s.T())
	s.httpService.AssertExpectations(s.T())
	s.checkService.AssertExpectations(s.T())
}

func (s *CheckTestSuite) TestCompareError() {
	// Arrange
	websiteID := uuid.New()
	website := domain.Website{
		ID:  websiteID,
		URL: "https://example.com",
	}
	currentContent := []byte("content")
	previousCheck := domain.Check{
		ID:        uuid.New(),
		WebsiteID: websiteID,
		Result:    []byte("previous content"),
	}
	compareErr := fmt.Errorf("comparison failed")

	s.websiteService.On("GetByID", s.ctx, websiteID).Return(website, nil)
	s.httpService.On("Request", website).Return(currentContent, nil)
	s.checkService.On("GetLatestCheckByWebsite", s.ctx, websiteID).Return(previousCheck, nil)
	s.diffService.On("Compare", previousCheck.Result, currentContent).Return(diff.Result{}, compareErr)

	// Act
	result, err := s.useCase.Check(s.ctx, websiteID)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to compare results")
	assert.ErrorIs(s.T(), err, compareErr)
	assert.Empty(s.T(), result)
	s.websiteService.AssertExpectations(s.T())
	s.httpService.AssertExpectations(s.T())
	s.checkService.AssertExpectations(s.T())
	s.diffService.AssertExpectations(s.T())
}

func (s *CheckTestSuite) TestCreateCheckError() {
	// Arrange
	websiteID := uuid.New()
	website := domain.Website{
		ID:  websiteID,
		URL: "https://example.com",
	}
	currentContent := []byte("content")
	diffResult := diff.Result{HasChanges: true}
	createErr := fmt.Errorf("failed to create check")

	s.websiteService.On("GetByID", s.ctx, websiteID).Return(website, nil)
	s.httpService.On("Request", website).Return(currentContent, nil)
	s.checkService.On("GetLatestCheckByWebsite", s.ctx, websiteID).Return(domain.Check{}, domain.ErrCheckNotFound)
	s.diffService.On("Compare", []byte(nil), currentContent).Return(diffResult, nil)
	s.checkService.On("CreateCheck", s.ctx, mock.Anything).Return(domain.Check{}, createErr)

	// Act
	result, err := s.useCase.Check(s.ctx, websiteID)

	// Assert
	assert.Error(s.T(), err)
	assert.Empty(s.T(), result)
	s.websiteService.AssertExpectations(s.T())
	s.httpService.AssertExpectations(s.T())
	s.checkService.AssertExpectations(s.T())
	s.diffService.AssertExpectations(s.T())
}

func (s *CheckTestSuite) TestNoPreviousCheck() {
	// Arrange
	websiteID := uuid.New()
	website := domain.Website{
		ID:  websiteID,
		URL: "https://example.com",
	}
	currentContent := []byte("new content")
	diffResult := diff.Result{HasChanges: true}
	s.useCase.diffService = diff.NewDiffService()

	s.websiteService.On("GetByID", s.ctx, websiteID).Return(website, nil)
	s.httpService.On("Request", website).Return(currentContent, nil)
	s.checkService.On("GetLatestCheckByWebsite", s.ctx, websiteID).Return(domain.Check{}, domain.ErrCheckNotFound)
	s.checkService.On("CreateCheck", s.ctx, mock.MatchedBy(func(check domain.Check) bool {
		return check.WebsiteID == websiteID &&
			string(check.Result) == string(currentContent)
	})).Return(domain.Check{}, nil)

	// Act
	result, err := s.useCase.Check(s.ctx, websiteID)

	// Assert
	assert.NoError(s.T(), err)
	expectedResult := domain.CheckResult{
		OldValue:   nil,
		NewValue:   currentContent,
		HasChanges: true,
		Check:      diffResult,
	}
	assert.Equal(s.T(), expectedResult.HasChanges, result.HasChanges)
	s.websiteService.AssertExpectations(s.T())
	s.httpService.AssertExpectations(s.T())
	s.checkService.AssertExpectations(s.T())
	s.diffService.AssertExpectations(s.T())
}

func (s *CheckTestSuite) TestViewSuccess() {
	// Arrange
	websiteID := uuid.New()
	website := domain.Website{
		ID:  websiteID,
		URL: "https://example.com",
	}
	expectedContent := []byte("processed content")

	s.websiteService.On("GetByID", s.ctx, websiteID).Return(website, nil)
	s.httpService.On("Request", website).Return(expectedContent, nil)

	// Act
	result, err := s.useCase.View(s.ctx, websiteID)

	// Assert
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expectedContent, result)
	s.websiteService.AssertExpectations(s.T())
	s.httpService.AssertExpectations(s.T())
}

func (s *CheckTestSuite) TestViewWebsiteNotFound() {
	// Arrange
	websiteID := uuid.New()
	expectedErr := domain.ErrWebsiteNotFound

	s.websiteService.On("GetByID", s.ctx, websiteID).Return(domain.Website{}, expectedErr)

	// Act
	result, err := s.useCase.View(s.ctx, websiteID)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to get website")
	assert.ErrorIs(s.T(), err, expectedErr)
	assert.Nil(s.T(), result)
	s.websiteService.AssertExpectations(s.T())
}

func (s *CheckTestSuite) TestViewRequestError() {
	// Arrange
	websiteID := uuid.New()
	website := domain.Website{
		ID:  websiteID,
		URL: "https://example.com",
	}
	requestErr := domain.ErrRequestFailed

	s.websiteService.On("GetByID", s.ctx, websiteID).Return(website, nil)
	s.httpService.On("Request", website).Return(nil, requestErr)
	s.checkService.On("CreateCheck", s.ctx, mock.MatchedBy(func(check domain.Check) bool {
		return check.WebsiteID == websiteID &&
			check.HasError == true &&
			check.ErrorMessage == requestErr.Error()
	})).Return(domain.Check{}, nil)

	// Act
	result, err := s.useCase.View(s.ctx, websiteID)

	// Assert
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "failed to make HTTP request")
	assert.ErrorIs(s.T(), err, requestErr)
	assert.Nil(s.T(), result)
	s.websiteService.AssertExpectations(s.T())
	s.httpService.AssertExpectations(s.T())
	s.checkService.AssertExpectations(s.T())
}
