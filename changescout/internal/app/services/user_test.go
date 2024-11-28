package services

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"github.com/gelleson/changescout/changescout/internal/app/services/mocks"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.UserRepository
	service  *UserService
	ctx      context.Context
}

func (s *UserServiceTestSuite) SetupTest() {
	s.mockRepo = mocks.NewUserRepository(s.T())
	s.service = NewUserService(s.mockRepo)
	s.ctx = context.Background()
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (s *UserServiceTestSuite) TestGetByID() {
	tests := []struct {
		name    string
		id      uuid.UUID
		mock    func()
		want    domain.User
		wantErr bool
	}{
		{
			name: "successful get",
			id:   uuid.New(),
			mock: func() {
				expectedUser := domain.User{ID: uuid.New(), Email: "test@example.com"}
				s.mockRepo.On("GetUserByID", s.ctx, mock.AnythingOfType("uuid.UUID")).Once().
					Return(expectedUser, nil)
			},
			want:    domain.User{},
			wantErr: false,
		},
		{
			name: "user not found",
			id:   uuid.New(),
			mock: func() {
				s.mockRepo.On("GetUserByID", s.ctx, mock.AnythingOfType("uuid.UUID")).Once().
					Return(domain.User{}, errors.New("user not found"))
			},
			want:    domain.User{},
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

func (s *UserServiceTestSuite) TestGetByEmail() {
	tests := []struct {
		name    string
		email   string
		mock    func()
		want    domain.User
		wantErr bool
	}{
		{
			name:  "successful get",
			email: "test@example.com",
			mock: func() {
				expectedUser := domain.User{ID: uuid.New(), Email: "test@example.com"}
				s.mockRepo.On("GetUserByEmail", s.ctx, "test@example.com").
					Return(expectedUser, nil)
			},
			want:    domain.User{},
			wantErr: false,
		},
		{
			name:  "user not found",
			email: "notfound@example.com",
			mock: func() {
				s.mockRepo.On("GetUserByEmail", s.ctx, "notfound@example.com").
					Return(domain.User{}, errors.New("user not found"))
			},
			want:    domain.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			got, err := s.service.GetByEmail(s.ctx, tt.email)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.NotEmpty(got)
		})
	}
}

func (s *UserServiceTestSuite) TestCreate() {
	tests := []struct {
		name    string
		user    domain.User
		mock    func()
		want    domain.User
		wantErr bool
	}{
		{
			name: "successful creation",
			user: domain.User{
				Email:    "test@example.com",
				Password: "password123",
			},
			mock: func() {
				expectedID := uuid.New()
				expectedUser := domain.User{
					ID:       expectedID,
					Email:    "test@example.com",
					Password: mock.Anything,
				}
				s.mockRepo.On("CreateUser", s.ctx, mock.MatchedBy(func(u domain.User) bool {
					return u.Email == "test@example.com" && len(u.Password) > 0 && bcrypt.CompareHashAndPassword([]byte(expectedUser.Password), []byte(u.Password)) != nil

				})).Once().Return(expectedUser, nil)
			},
			wantErr: false,
		},
		{
			name: "creation failed",
			user: domain.User{
				Email:    "test@example.com",
				Password: "password123",
			},
			mock: func() {
				s.mockRepo.On("CreateUser", s.ctx, mock.AnythingOfType("domain.User")).Once().
					Return(domain.User{}, errors.New("creation failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.mock()
			got, err := s.service.Create(s.ctx, tt.user)
			if tt.wantErr {
				s.Error(err)
				return
			}
			s.NoError(err)
			s.NotEmpty(got)
			s.NotEqual(tt.user.Password, got.Password) // Password should be hashed
		})
	}
}
