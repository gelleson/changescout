package broker

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/garsue/watermillzap"
	"go.uber.org/zap"
	"strings"
)

type Provider interface {
	message.Publisher
	message.Subscriber
}

type PubSubLib map[providerName]Provider

type Broker struct {
	router   *message.Router
	pubsubs  PubSubLib
	usecases *UseCases
	logger   *zap.Logger
}

type providerName string

const (
	MainProviderName providerName = "main"
)

func New(logger *zap.Logger, usecases *UseCases, pubsubs PubSubLib) *Broker {
	return &Broker{
		router: message.NewDefaultRouter(
			watermillzap.NewLogger(logger),
		),
		pubsubs:  pubsubs,
		logger:   zap.L(),
		usecases: usecases,
	}
}

func (b *Broker) AddHandler(topic string, handler func(msg *message.Message) error) {
	b.router.AddNoPublisherHandler(
		strings.ReplaceAll(topic, ".", "_"),
		topic,
		b.pubsubs[MainProviderName],
		handler,
	)
}

func (b *Broker) Run(ctx context.Context) error {
	return b.router.Run(ctx)
}
