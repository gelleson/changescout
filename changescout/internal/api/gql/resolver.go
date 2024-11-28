package gql

import (
	"github.com/gelleson/changescout/changescout/internal/app/services"
	"github.com/gelleson/changescout/changescout/internal/app/usecases"
	"github.com/gelleson/changescout/changescout/internal/app/usecases/auth"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	WebsiteUseCase      *usecases.WebsiteUseCase
	AuthUseCase         *auth.UseCase
	NotificationService *services.NotificationService
}
