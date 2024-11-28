package contexts

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
)

type UserContextKey struct{}

func UserContext(ctx context.Context) (*domain.AuthClaims, bool) {
	user, isOK := ctx.Value(UserContextKey{}).(*domain.AuthClaims)
	return user, isOK
}

func WithUserContext(ctx context.Context, user *domain.AuthClaims) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, UserContextKey{}, user)
}
