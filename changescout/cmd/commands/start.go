package commands

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/garsue/watermillzap"
	"github.com/gelleson/changescout/changescout/internal/api/broker"
	"github.com/gelleson/changescout/changescout/internal/api/gql"
	httpplatform "github.com/gelleson/changescout/changescout/internal/api/http"
	"github.com/gelleson/changescout/changescout/internal/api/http/middlewares"
	"github.com/gelleson/changescout/changescout/internal/app/services"
	"github.com/gelleson/changescout/changescout/internal/app/services/diff"
	"github.com/gelleson/changescout/changescout/internal/app/services/requesters"
	"github.com/gelleson/changescout/changescout/internal/app/services/sender"
	"github.com/gelleson/changescout/changescout/internal/app/services/sender/providers/telegram"
	"github.com/gelleson/changescout/changescout/internal/app/usecases"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/check"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/notification"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/scheduler"
	"github.com/gelleson/changescout/changescout/internal/domain"
	entrepo "github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent"
	"github.com/gelleson/changescout/changescout/internal/pkg/clis"
	"github.com/gelleson/changescout/changescout/internal/platform/logger"
	"github.com/gelleson/changescout/changescout/internal/utils/transform"
	"github.com/gelleson/changescout/changescout/pkg/flags"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"log"
)

var StartServer = &cli.Command{
	Name:  "start",
	Usage: "Start the service",
	Flags: flags.Build(
		clis.FlagsLogLevel,
		clis.FlagsPort,
		clis.FlagsSecret,
		clis.FlagsSecretExpiration,
		clis.FlagsDBUrl,
		clis.FlagsDBEngine,
		clis.FlagsBrokerEnabled,
		clis.FlagsSchedulerEnabled,
		clis.FlagsSchedulerInterval,
		clis.FlagsBrowserManagedInstanceURL,
		clis.FlagsBrowserDisable,
	),
	Action: func(c *cli.Context) error {
		logger.SetLevel(clis.FlagsLogLevel.Get(c))
		server := httpplatform.New(
			logger.L("http"),
			httpplatform.WithPort(clis.FlagsPort.Get(c)),
		)
		client, err := entrepo.Build(c.Context, &entrepo.BuildConfig{
			DBEngine: clis.FlagsDBEngine.Get(c),
			DBURL:    clis.FlagsDBUrl.Get(c),
		})
		if err != nil {
			log.Fatal("failed to open database connection", zap.Error(err))
		}
		server.WithMiddlewares(
			logger.WithLogger(logger.L("http")),
			middlewares.JWTAuth(middlewares.JWTAuthConfig{
				Secret: []byte(clis.FlagsSecret.Get(c)),
			}),
		)

		pubSub := gochannel.NewGoChannel(
			gochannel.Config{
				OutputChannelBuffer: 100,
			},
			watermillzap.NewLogger(logger.L("pubsub")),
		)

		b := broker.New(
			logger.L("broker"),
			&broker.UseCases{
				CheckUseCase: check.NewUseCase(
					services.NewWebsiteService(
						entrepo.NewWebsiteRepository(client),
					),
					requesters.New(requesters.Options{
						Browser: requesters.BrowserOption{
							Enable: !clis.FlagsBrowserDisable.Get(c),
							ManagedInstanceURL: transform.ToPtr(
								clis.FlagsBrowserManagedInstanceURL.Get(c),
							),
						},
					}),
					services.NewCheckService(
						entrepo.NewCheckRepository(client),
					),
					diff.NewDiffService(),
				),
				WebsiteUseCase: usecases.NewWebsiteUseCase(
					services.NewWebsiteService(
						entrepo.NewWebsiteRepository(client),
					),
					services.NewUserService(
						entrepo.NewUserRepository(client),
					),
				),
				NotificationUseCase: notification.NewUseCase(
					sender.NewSenderService(sender.Senders{
						domain.TelegramNotificationType: telegram.New(),
					}),
					services.NewWebsiteService(
						entrepo.NewWebsiteRepository(client),
					),
					services.NewNotificationService(
						entrepo.NewNotificationRepository(client),
					),
				),
			},
			broker.PubSubLib{
				broker.MainProviderName: pubSub,
			},
		)

		s := scheduler.NewUseCase(
			pubSub,
			services.NewWebsiteService(
				entrepo.NewWebsiteRepository(client),
			),
			clis.FlagsSchedulerInterval.Get(c),
		)

		if clis.FlagsBrokerEnabled.Get(c) {
			b.AddHandler(
				"websites.check",
				b.HandleWebsiteCheck(),
			)
			go b.Run(context.Background())
		}

		if clis.FlagsSchedulerEnabled.Get(c) {
			go s.Run(context.Background())
		}

		ghandler := gql.BuildHandler(&gql.HandlerConfig{
			Secret:           clis.FlagsSecret.Get(c),
			SecretExpiration: clis.FlagsSecretExpiration.Get(c),
			Client:           client,
		})
		server.Register("POST", "/query", func(c echo.Context) error {
			ghandler.Schema().ServeHTTP(c.Response(), c.Request())
			return nil
		})

		server.Register("GET", "/playground", func(c echo.Context) error {
			ghandler.PlaygroundHandler().ServeHTTP(c.Response(), c.Request())
			return nil
		})

		log.Fatal(server.Start())

		return nil
	},
}
