package auth

import (
	"context"
	"fmt"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	authpkg "github.com/gelleson/changescout/changescout/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"

	"github.com/gelleson/changescout/changescout/internal/domain"
)

// UserRepository interface (should already exist from your previous code)
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
}

// UseCase handles authentication logic, including user validation and JWT token generation.
type UseCase struct {
	userRepository UserRepository
	jwtSecret      []byte
	tokenExpiry    time.Duration
}

// NewUseCase initializes a new AuthUseCase with the given UserRepository, JWT secret, and token expiry duration.
func NewUseCase(userRepository UserRepository, jwtSecret []byte, tokenExpiry time.Duration) *UseCase {
	return &UseCase{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
		tokenExpiry:    tokenExpiry,
	}
}

// AuthenticateByPassword validates the user against the provided email and password, and returns a JWT token if successful.
// Returns a token string and an error if authentication fails.
func (uc *UseCase) AuthenticateByPassword(ctx context.Context, email string, password string) (string, error) {
	user, err := uc.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := authpkg.CheckPassword(user.Password, password); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := uc.generateJWT(user)
	if err != nil {
		return "", fmt.Errorf("generating JWT: %w", err)
	}
	return token, nil
}

// RegistrationByPassword registers a new user with the provided credentials.
func (uc *UseCase) RegistrationByPassword(ctx context.Context, firstName, lastName, email, password string) (string, error) {
	existingUser, err := uc.userRepository.GetByEmail(ctx, email)
	if err != nil && !database.IsNotFound(err) {
		return "", fmt.Errorf("getting user by email: %w", err)
	}
	if existingUser.ID != uuid.Nil {
		return "", domain.ErrUserAlreadyExists
	}

	user := domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		Role:      domain.Regular,
	}
	newUser, err := uc.userRepository.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("creating user: %w", err)
	}

	token, err := uc.generateJWT(newUser)
	if err != nil {
		return "", fmt.Errorf("generating JWT: %w", err)
	}

	return token, nil
}

// generateJWT generates a JWT token for the provided user.
// Returns the token string and an error if the JWT generation fails.
func (uc *UseCase) generateJWT(user domain.User) (string, error) {
	claims := domain.AuthClaims{
		ID:        user.ID,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(uc.tokenExpiry),
	}
	claims.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:   "changescout",
		Subject:  "changescout",
		Audience: []string{"changescout"},
	}

	return authpkg.GenerateJWTWithClaims(claims, uc.jwtSecret)
}
