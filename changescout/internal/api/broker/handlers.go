package broker

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"go.uber.org/zap"
)

type HandlerFunc func(msg *message.Message) error

func (b *Broker) HandleWebsiteCheck() HandlerFunc {
	return func(msg *message.Message) error {
		defer msg.Ack()
		var site domain.Website
		err := json.Unmarshal(msg.Payload, &site)
		if err != nil {
			return err
		}

		b.logger.Debug("Checking website", zap.String("url", site.URL))

		res, err := b.usecases.CheckUseCase.Check(msg.Context(), site.ID)
		if err != nil {
			return err
		}

		if !res.HasChanges {
			return nil
		}

		if err := b.usecases.WebsiteUseCase.UpdateLastCheck(msg.Context(), site.ID); err != nil {
			b.logger.Error("Failed to update last check", zap.Error(err))
			return err
		}

		if err := b.usecases.NotificationUseCase.NotifyChanges(msg.Context(), site.ID, res); err != nil {
			return err
		}

		return nil
	}
}
