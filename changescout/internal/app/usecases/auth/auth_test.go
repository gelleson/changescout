package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gelleson/changescout/changescout/internal/app/usecases/auth/mocks"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AuthUseCaseTestSuite struct {
	suite.Suite
	useCase      *UseCase
	mockUserRepo *mocks.UserService
	jwtSecret    []byte
	ctx          context.Context
	tokenExpiry  time.Duration
}

// SetupTest sets up the environment for each test.
func (suite *AuthUseCaseTestSuite) SetupTest() {
	suite.mockUserRepo = new(mocks.UserService)
	suite.jwtSecret = []byte("secret")
	suite.tokenExpiry = time.Hour
	suite.ctx = context.TODO()
	suite.useCase = NewUseCase(suite.mockUserRepo, suite.jwtSecret, suite.tokenExpiry)
}

// Test case for invalid credentials
func (suite *AuthUseCaseTestSuite) TestAuthenticateByPassword_InvalidCredentials() {
	email := "user@example.com"
	wrongPassword := "wrongpassword"
	hashedPassword := "correctHash"

	suite.mockUserRepo.On("GetByEmail", suite.ctx, email).Return(domain.User{Email: email, Password: hashedPassword}, nil)
	authpkg := new(mocks.UserRepository)
	authpkg.On("CheckPassword", hashedPassword, wrongPassword).Return(errors.New("invalid credentials"))

	_, err := suite.useCase.AuthenticateByPassword(suite.ctx, email, wrongPassword)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), domain.ErrInvalidCredentials, err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// Test case for user registration when user already exists
func (suite *AuthUseCaseTestSuite) TestRegistrationByPassword_UserAlreadyExists() {
	email := "user@example.com"
	existingUser := domain.User{ID: uuid.New(), Email: email}

	suite.mockUserRepo.On("GetByEmail", suite.ctx, email).Return(existingUser, nil)

	_, err := suite.useCase.RegistrationByPassword(suite.ctx, "John", "Doe", email, "password")

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), domain.ErrUserAlreadyExists, err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// Test case for successful user registration with role fetching
func (suite *AuthUseCaseTestSuite) TestRegistrationByPassword_SuccessWithRole() {
	email := "newuser@example.com"
	password := "password"
	role := domain.RoleUser
	user := domain.User{FirstName: "John", LastName: "Doe", Email: email, Password: password, Role: role}

	suite.mockUserRepo.On("GetByEmail", suite.ctx, email).Return(domain.User{}, database.ErrEntityNotFound)
	suite.mockUserRepo.On("Create", suite.ctx, mock.AnythingOfType("domain.User")).Return(user, nil)
	suite.mockUserRepo.On("GetRole", suite.ctx).Return(role, nil)

	token, err := suite.useCase.RegistrationByPassword(suite.ctx, "John", "Doe", email, password)

	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// TestSuite runs the suite to execute the tests
func TestAuthUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(AuthUseCaseTestSuite))
}
