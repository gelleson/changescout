package broker

import (
	"github.com/gelleson/changescout/changescout/internal/app/usecases"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/check"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/notification"
)

type UseCases struct {
	CheckUseCase        *check.UseCase
	WebsiteUseCase      *usecases.WebsiteUseCase
	NotificationUseCase *notification.UseCase
}
