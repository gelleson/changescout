package gql

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gelleson/changescout/changescout/internal/api/gql/generated"
	"github.com/gelleson/changescout/changescout/internal/app/services"
	"github.com/gelleson/changescout/changescout/internal/app/services/diff"
	http2 "github.com/gelleson/changescout/changescout/internal/app/services/requesters/http"
	"github.com/gelleson/changescout/changescout/internal/app/usecases"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/auth"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/check"
	entrepo "github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent"
	"github.com/gelleson/changescout/changescout/internal/pkg/contexts"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type HandlerConfig struct {
	Secret string
	Client *ent.Client
}

type Handler struct {
	schema            *handler.Server
	playgroundHandler http.Handler
}

func (h *Handler) Schema() *handler.Server {
	return h.schema
}

func (h *Handler) PlaygroundHandler() http.Handler {
	return h.playgroundHandler
}

func BuildHandler(conf *HandlerConfig) *Handler {
	return &Handler{
		schema: handler.NewDefaultServer(
			generated.NewExecutableSchema(
				generated.Config{
					Resolvers: &Resolver{
						WebsiteUseCase: usecases.NewWebsiteUseCase(
							services.NewWebsiteService(
								entrepo.NewWebsiteRepository(conf.Client),
							),
							services.NewUserService(
								entrepo.NewUserRepository(conf.Client),
							),
						),
						AuthUseCase: auth.NewUseCase(
							services.NewUserService(
								entrepo.NewUserRepository(conf.Client),
							),
							[]byte(conf.Secret),
							time.Hour*24,
						),
						NotificationService: services.NewNotificationService(
							entrepo.NewNotificationRepository(conf.Client),
						),
						CheckUseCase: check.NewUseCase(
							services.NewWebsiteService(
								entrepo.NewWebsiteRepository(conf.Client),
							),
							http2.New(http.DefaultClient),
							services.NewCheckService(
								entrepo.NewCheckRepository(conf.Client),
							),
							diff.NewDiffService(),
						),
					},
					Directives: generated.DirectiveRoot{
						IsAuthenticated: func(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
							user, isAuthenticated := contexts.UserContext(ctx)
							if !isAuthenticated {
								return nil, errors.New("not authenticated")
							}
							if user.ID == uuid.Nil {
								return nil, errors.New("not authenticated")
							}
							return next(ctx)
						},
					},
				},
			),
		),
		playgroundHandler: playground.Handler("GraphQL", "/query"),
	}
}
