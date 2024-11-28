package sender

import "github.com/gelleson/changescout/changescout/internal/domain"

type Sender interface {
	Send(notification string, conf domain.Notification) error
}

type Senders map[domain.NotificationType]Sender

type SenderService struct {
	senders Senders
}

func NewSenderService(senders Senders) *SenderService {
	return &SenderService{
		senders: senders,
	}
}

func (s *SenderService) Send(notification string, conf domain.Notification) error {
	sender, ok := s.senders[conf.Type]
	if !ok {
		return nil
	}

	return sender.Send(notification, conf)
}
