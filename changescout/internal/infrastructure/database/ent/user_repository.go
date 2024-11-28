package ent

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/user"
	"github.com/google/uuid"
	"time"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

func roleDomainToEnt(role domain.Role) user.Role {
	switch role {
	case domain.Admin:
		return user.RoleAdmin
	case domain.Regular:
		return user.RoleUser
	default:
		panic("invalid role")
	}
}

func entToRoleDomain(role user.Role) domain.Role {
	switch role {
	case user.RoleAdmin:
		return domain.Admin
	case user.RoleUser:
		return domain.Regular
	default:
		panic("invalid role")
	}
}

func entToUser(user *ent.User) domain.User {
	return domain.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
		Role:      entToRoleDomain(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u UserRepository) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	usr, err := u.client.
		User.
		Create().
		SetFirstName(user.FirstName).
		SetLastName(user.LastName).
		SetEmail(user.Email).
		SetPassword(user.Password).
		SetRole(roleDomainToEnt(user.Role)).
		SetIsActive(user.IsActive).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)

	if err != nil {
		return domain.User{}, err
	}

	return entToUser(usr), err
}

func (u UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	entity, err := u.client.
		User.
		Query().
		Where(user.ID(id)).
		Only(ctx)

	if err != nil {
		return domain.User{}, err
	}

	return entToUser(entity), nil
}

func (u UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	entity, err := u.client.
		User.
		Query().
		Where(user.Email(email)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return domain.User{}, database.ErrEntityNotFound
		}
		return domain.User{}, err
	}

	return entToUser(entity), nil
}

func (u UserRepository) UpdateUser(ctx context.Context, user domain.User) error {
	_user := u.client.User.UpdateOneID(user.ID)
	_mutation := _user.Mutation()

	if user.FirstName != "" {
		_mutation.SetFirstName(user.FirstName)
	}
	if user.LastName != "" {
		_mutation.SetLastName(user.LastName)
	}
	if user.Email != "" {
		_mutation.SetEmail(user.Email)
	}
	if user.Password != "" {
		_mutation.SetPassword(user.Password)
	}
	if user.Role != "" {
		_mutation.SetRole(roleDomainToEnt(user.Role))
	}
	if user.IsActive != false {
		_mutation.SetIsActive(user.IsActive)
	}

	_, err := _user.Save(ctx)

	return err
}
