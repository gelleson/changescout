package usecases

import (
	"context"
	"errors"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/gelleson/changescout/changescout/pkg/crons"
	"github.com/google/uuid"
	"time"
)

type WebsiteService interface {
	Create(ctx context.Context, website domain.Website) (domain.Website, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Website, error)
	GetByURL(ctx context.Context, url string) (domain.Website, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, pagination domain.Pagination) ([]domain.Website, error)
	Update(ctx context.Context, website domain.Website) (domain.Website, error)
	GetDueForCheck(ctx context.Context, pagination domain.Pagination) ([]domain.Website, error)
	List(ctx context.Context, filters database.WebsiteFilters, pagination domain.Pagination) ([]domain.Website, int, error)
	GetByStatus(ctx context.Context, enabled bool, pagination domain.Pagination) ([]domain.Website, error)
	UpdateLastCheck(ctx context.Context, websiteID uuid.UUID) error
	UpdateStatus(ctx context.Context, websiteID uuid.UUID, enabled bool) (domain.Website, error)

	Delete(ctx context.Context, id uuid.UUID) error
}

type UserService interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type WebsiteUseCase struct {
	websiteService WebsiteService
	userService    UserService
}

func NewWebsiteUseCase(websiteService WebsiteService, userService UserService) *WebsiteUseCase {
	return &WebsiteUseCase{
		websiteService: websiteService,
		userService:    userService,
	}
}

func (u WebsiteUseCase) Create(ctx context.Context, userId uuid.UUID, website domain.Website) (domain.Website, error) {
	_, err := u.userService.GetByID(ctx, userId)
	if err != nil {
		return domain.Website{}, domain.ErrInvalidToken
	}
	website.UserID = userId

	return u.websiteService.Create(ctx, website)
}

func (u WebsiteUseCase) GetByID(ctx context.Context, userId uuid.UUID, id uuid.UUID) (domain.Website, error) {
	_, err := u.userService.GetByID(ctx, userId)
	if err != nil {
		return domain.Website{}, err
	}

	return u.websiteService.GetByID(ctx, id)
}

func (u WebsiteUseCase) GetByURL(ctx context.Context, userId uuid.UUID, url string) (domain.Website, error) {
	_, err := u.userService.GetByID(ctx, userId)
	if err != nil {
		return domain.Website{}, err
	}

	return u.websiteService.GetByURL(ctx, url)
}

func (u WebsiteUseCase) GetByUserID(ctx context.Context, userId uuid.UUID, pagination domain.Pagination) ([]domain.Website, error) {
	_, err := u.userService.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return u.websiteService.GetByUserID(ctx, userId, pagination)
}

func (u WebsiteUseCase) Update(ctx context.Context, userId uuid.UUID, id uuid.UUID, website domain.Website) (domain.Website, error) {
	_, err := u.userService.GetByID(ctx, userId)
	if err != nil {
		return domain.Website{}, err
	}
	site, err := u.websiteService.GetByID(ctx, id)
	if err != nil {
		return domain.Website{}, err
	}
	if site.UserID != userId {
		return domain.Website{}, errors.New("not authorized")
	}

	if website.Cron != "" {
		c := crons.NewScheduler()
		if err := c.Validate(website.Cron); err != nil {
			return domain.Website{}, err
		}
		nextCheckAt, err := c.NextRun(website.Cron, time.Now())
		if err != nil {
			return domain.Website{}, err
		}
		website.NextCheckAt = &nextCheckAt
	}

	return u.websiteService.Update(ctx, website)
}

func (u WebsiteUseCase) UpdateLastCheck(ctx context.Context, id uuid.UUID) error {
	return u.websiteService.UpdateLastCheck(ctx, id)
}

func (u WebsiteUseCase) UpdateStatus(ctx context.Context, userId uuid.UUID, id uuid.UUID, enabled bool) (domain.Website, error) {
	_, err := u.userService.GetByID(ctx, userId)
	if err != nil {
		return domain.Website{}, err
	}
	site, err := u.websiteService.GetByID(ctx, id)
	if err != nil {
		return domain.Website{}, err
	}
	if site.UserID != userId {
		return domain.Website{}, errors.New("not authorized")
	}

	return u.websiteService.UpdateStatus(ctx, id, enabled)
}

func (u WebsiteUseCase) Delete(ctx context.Context, userId uuid.UUID, id uuid.UUID) error {
	_, err := u.userService.GetByID(ctx, userId)
	if err != nil {
		return err
	}
	site, err := u.websiteService.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if site.UserID != userId {
		return errors.New("not authorized")
	}

	return u.websiteService.Delete(ctx, id)
}
