package gql

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gelleson/changescout/changescout/internal/api/gql/directive"
	"github.com/gelleson/changescout/changescout/internal/api/gql/generated"
	"github.com/gelleson/changescout/changescout/internal/app/services"
	"github.com/gelleson/changescout/changescout/internal/app/services/diff"
	http2 "github.com/gelleson/changescout/changescout/internal/app/services/requesters/http"
	"github.com/gelleson/changescout/changescout/internal/app/usecases"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/auth"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/check"
	entrepo "github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent"
	"net/http"
	"time"
)

type HandlerConfig struct {
	Secret           string
	SecretExpiration time.Duration
	Client           *ent.Client
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
							conf.SecretExpiration,
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
						IsAuthenticated: directive.IsAuth(),
						HasRole:         directive.HasRole(),
					},
				},
			),
		),
		playgroundHandler: playground.Handler("GraphQL", "/query"),
	}
}
