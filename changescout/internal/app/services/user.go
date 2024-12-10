package services

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockery --name UserRepository
type UserRepository interface {
	database.UserRepository
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository database.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (u UserService) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return u.userRepository.GetUserByID(ctx, id)
}

func (u UserService) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return u.userRepository.GetUserByEmail(ctx, email)
}

func (u UserService) Create(ctx context.Context, user domain.User) (domain.User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}
	user.Password = string(password)
	return u.userRepository.CreateUser(ctx, user)
}

// GetRole determines the role based on the number of users in the system.
// If there are no users, it assigns the RoleAdmin to the first user.
// Otherwise, it assigns the RoleUser.
func (u UserService) GetRole(ctx context.Context) (domain.Role, error) {
	total, err := u.userRepository.GetTotalUsers(ctx)
	if err != nil {
		return domain.RoleUser, err
	}
	if total == 0 {
		return domain.RoleAdmin, nil
	}
	return domain.RoleUser, nil
}
