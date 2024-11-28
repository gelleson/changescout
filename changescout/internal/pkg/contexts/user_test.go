package contexts

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserContext(t *testing.T) {
	mockUser := &domain.AuthClaims{
		ID:        uuid.New(),
		Email:     "john.doe@example.com",
		ExpiresAt: time.Now().Add(time.Hour),
	}

	tests := []struct {
		name     string
		ctx      context.Context
		wantUser *domain.AuthClaims
		wantIsOK bool
	}{
		{
			name:     "context has user",
			ctx:      WithUserContext(context.Background(), mockUser),
			wantUser: mockUser,
			wantIsOK: true,
		},
		{
			name:     "context has no user",
			ctx:      context.Background(),
			wantUser: nil,
			wantIsOK: false,
		},
		{
			name:     "context has incorrect type",
			ctx:      context.WithValue(context.Background(), UserContextKey{}, "incorrect_type"),
			wantUser: nil,
			wantIsOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, gotIsOK := UserContext(tt.ctx)
			if gotIsOK != tt.wantIsOK {
				t.Errorf("UserContext() gotIsOK = %v, want %v", gotIsOK, tt.wantIsOK)
			}
			if gotUser != tt.wantUser {
				t.Errorf("UserContext() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestWithUserContext(t *testing.T) {
	user := &domain.AuthClaims{
		ID:    uuid.New(),
		Email: "john.doe@example.com",
	}

	tests := []struct {
		name string
		ctx  context.Context
	}{
		{"With existing context", context.TODO()},
		{"With nil context", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := WithUserContext(tt.ctx, user)
			assert.NotNil(t, ctx, "Context should not be nil")

			extractedUser, ok := ctx.Value(UserContextKey{}).(*domain.AuthClaims)
			assert.True(t, ok, "User should be in the context")
			assert.Equal(t, user, extractedUser, "Extracted user should be equal to the original user")
		})
	}
}
