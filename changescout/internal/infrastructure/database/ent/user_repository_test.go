package ent

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/user"
	"github.com/gelleson/changescout/changescout/internal/utils/testdb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	repo *UserRepository
	ctx  context.Context
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

func (s *UserRepositoryTestSuite) SetupTest() {
	client := testdb.NewEntClient()
	s.repo = NewUserRepository(client)
	s.ctx = context.Background()
}

func (s *UserRepositoryTestSuite) TestCreateUser() {
	tests := []struct {
		name    string
		user    domain.User
		wantErr bool
	}{
		{
			name: "create regular user",
			user: domain.User{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "hashedpassword",
				Role:      domain.RoleUser,
				IsActive:  true,
			},
			wantErr: false,
		},
		{
			name: "create admin user",
			user: domain.User{
				FirstName: "RoleAdmin",
				LastName:  "User",
				Email:     "admin@example.com",
				Password:  "hashedpassword",
				Role:      domain.RoleAdmin,
				IsActive:  true,
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			user: domain.User{
				FirstName: "Duplicate",
				LastName:  "User",
				Email:     "john@example.com", // Same as first test case
				Password:  "hashedpassword",
				Role:      domain.RoleUser,
				IsActive:  true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.CreateUser(s.ctx, tt.user)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.NotEqual(s.T(), uuid.Nil, got.ID)
			assert.Equal(s.T(), tt.user.FirstName, got.FirstName)
			assert.Equal(s.T(), tt.user.LastName, got.LastName)
			assert.Equal(s.T(), tt.user.Email, got.Email)
			assert.Equal(s.T(), tt.user.Password, got.Password)
			assert.Equal(s.T(), tt.user.Role, got.Role)
			assert.Equal(s.T(), tt.user.IsActive, got.IsActive)
			assert.NotZero(s.T(), got.CreatedAt)
			assert.NotZero(s.T(), got.UpdatedAt)
		})
	}
}

func (s *UserRepositoryTestSuite) TestGetUserByID() {
	// Create test user
	user := domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "hashedpassword",
		Role:      domain.RoleUser,
		IsActive:  true,
	}

	created, err := s.repo.CreateUser(s.ctx, user)
	assert.NoError(s.T(), err)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "existing user",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.GetUserByID(s.ctx, tt.id)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.Equal(s.T(), created.ID, got.ID)
			assert.Equal(s.T(), created.Email, got.Email)
		})
	}
}

func (s *UserRepositoryTestSuite) TestGetUserByEmail() {
	// Create test user
	user := domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "hashedpassword",
		Role:      domain.RoleUser,
		IsActive:  true,
	}

	created, err := s.repo.CreateUser(s.ctx, user)
	assert.NoError(s.T(), err)

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "existing email",
			email:   created.Email,
			wantErr: false,
		},
		{
			name:    "non-existing email",
			email:   "nonexistent@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, err := s.repo.GetUserByEmail(s.ctx, tt.email)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)
			assert.Equal(s.T(), created.Email, got.Email)
			assert.Equal(s.T(), created.ID, got.ID)
		})
	}
}

func (s *UserRepositoryTestSuite) TestUpdateUser() {
	// Create initial user
	user := domain.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "hashedpassword",
		Role:      domain.RoleUser,
		IsActive:  true,
	}

	created, err := s.repo.CreateUser(s.ctx, user)
	assert.NoError(s.T(), err)

	tests := []struct {
		name    string
		update  domain.User
		wantErr bool
		verify  func(assert.TestingT, domain.User)
	}{
		{
			name: "update name",
			update: domain.User{
				ID:        created.ID,
				FirstName: "Jane",
				LastName:  "Smith",
			},
			wantErr: false,
			verify: func(t assert.TestingT, u domain.User) {
				assert.Equal(t, "Jane", u.FirstName)
				assert.Equal(t, "Smith", u.LastName)
			},
		},
		{
			name: "update email",
			update: domain.User{
				ID:    created.ID,
				Email: "jane@example.com",
			},
			wantErr: false,
			verify: func(t assert.TestingT, u domain.User) {
				assert.Equal(t, "jane@example.com", u.Email)
			},
		},
		{
			name: "update role",
			update: domain.User{
				ID:   created.ID,
				Role: domain.RoleAdmin,
			},
			wantErr: false,
			verify: func(t assert.TestingT, u domain.User) {
				assert.Equal(t, domain.RoleAdmin, u.Role)
			},
		},
		{
			name: "update non-existent user",
			update: domain.User{
				ID:        uuid.New(),
				FirstName: "Should Fail",
			},
			wantErr: true,
			verify:  func(t assert.TestingT, u domain.User) {},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			err := s.repo.UpdateUser(s.ctx, tt.update)
			if tt.wantErr {
				assert.Error(s.T(), err)
				return
			}

			assert.NoError(s.T(), err)

			// Verify the update
			updated, err := s.repo.GetUserByID(s.ctx, tt.update.ID)
			assert.NoError(s.T(), err)
			tt.verify(s.T(), updated)
		})
	}
}

func (s *UserRepositoryTestSuite) TestRoleConversion() {
	tests := []struct {
		name     string
		domain   domain.Role
		expected user.Role
	}{
		{
			name:     "admin role",
			domain:   domain.RoleAdmin,
			expected: user.RoleAdmin,
		},
		{
			name:     "regular role",
			domain:   domain.RoleUser,
			expected: user.RoleUser,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			entRole := roleDomainToEnt(tt.domain)
			assert.Equal(s.T(), tt.expected, entRole)

			// Test conversion back to domain
			domainRole := entToRoleDomain(entRole)
			assert.Equal(s.T(), tt.domain, domainRole)
		})
	}
}

func (s *CheckRepositoryTestSuite) TestGetTotalUsers() {
	// Ensure that the test user is created in SetupTest

	// Get the total number of users
	count, err := s.userRepo.GetTotalUsers(s.ctx)
	assert.NoError(s.T(), err)

	// Check that the number of users is at least 1 (since SetupTest adds one user)
	assert.GreaterOrEqual(s.T(), count, 1)

	// Additional checks can be added by creating more users
	// Create an additional test user
	newUser := domain.User{
		Email:     "newuser@example.com",
		Role:      domain.RoleUser,
		IsActive:  true,
		Password:  "password",
		FirstName: "New",
		LastName:  "User",
	}

	_, err = s.userRepo.CreateUser(s.ctx, newUser)
	assert.NoError(s.T(), err)

	// Check the user count again
	count, err = s.userRepo.GetTotalUsers(s.ctx)
	assert.NoError(s.T(), err)
	assert.GreaterOrEqual(s.T(), count, 2)
}
