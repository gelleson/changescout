package directive

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/pkg/contexts"
	"github.com/google/uuid"
)

func isAuth(ctx context.Context) (*domain.AuthClaims, error) {
	user, isAuthenticated := contexts.UserContext(ctx)
	if !isAuthenticated {
		return nil, errors.New("not authenticated")
	}
	if user.ID == uuid.Nil {
		return nil, errors.New("not authenticated")
	}

	return user, nil
}

type Handler func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error)

func IsAuth() Handler {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
		_, err = isAuth(ctx)
		if err != nil {
			return nil, err
		}
		return next(ctx)
	}
}

func HasRole() func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []domain.Role) (res interface{}, err error) {
	return func(ctx context.Context, obj interface{}, next graphql.Resolver, roles []domain.Role) (res interface{}, err error) {
		user, err := isAuth(ctx)
		if err != nil {
			return nil, err
		}
		for _, role := range roles {
			if user.Role == role {
				return next(ctx)
			}
		}
		return nil, errors.New("not authorized")
	}
}
