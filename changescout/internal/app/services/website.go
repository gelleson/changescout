package services

import (
	"context"
	"errors"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/gelleson/changescout/changescout/internal/utils/transform"
	"github.com/gelleson/changescout/changescout/pkg/crons"
	"github.com/gelleson/changescout/changescout/pkg/validators"
	"github.com/google/uuid"
	"time"
)

const MaxInterval = time.Hour * 24

//go:generate mockery --name WebsiteRepository
type WebsiteRepository interface {
	database.WebsiteRepository
}

type WebsiteService struct {
	websiteRepository WebsiteRepository
	scheduler         *crons.Scheduler
}

func NewWebsiteService(websiteRepository database.WebsiteRepository) *WebsiteService {
	return &WebsiteService{
		websiteRepository: websiteRepository,
		scheduler:         crons.NewScheduler(),
	}
}

func (w WebsiteService) Create(ctx context.Context, website domain.Website) (domain.Website, error) {
	if err := checkMode(website.Mode); err != nil {
		return domain.Website{}, err
	}

	if !validators.IsValidURL(website.URL) {
		return domain.Website{}, errors.New("invalid url")
	}

	c := crons.NewScheduler()
	if err := c.Validate(website.Cron); err != nil {
		return domain.Website{}, err
	}

	nextCheckAt, err := c.NextRun(website.Cron, time.Now())
	if err != nil {
		return domain.Website{}, err
	}

	website.NextCheckAt = &nextCheckAt

	return w.websiteRepository.CreateWebsite(ctx, website)
}

func (w WebsiteService) GetByID(ctx context.Context, id uuid.UUID) (domain.Website, error) {
	return w.websiteRepository.GetWebsiteByID(ctx, id)
}

func (w WebsiteService) GetByURL(ctx context.Context, url string) (domain.Website, error) {
	return w.websiteRepository.GetWebsiteByURL(ctx, url)
}

func (w WebsiteService) GetByUserID(ctx context.Context, userID uuid.UUID, pagination domain.Pagination) ([]domain.Website, error) {
	return w.websiteRepository.GetWebsiteByUserID(ctx, userID, pagination)
}

func (w WebsiteService) Update(ctx context.Context, website domain.Website) (domain.Website, error) {
	return w.websiteRepository.UpdateWebsite(ctx, website)
}

func (w WebsiteService) GetDueForCheck(ctx context.Context, pagination domain.Pagination) ([]domain.Website, error) {
	return w.websiteRepository.GetWebsitesDueForCheck(ctx, pagination)
}

func (w WebsiteService) List(ctx context.Context, filters database.WebsiteFilters, pagination domain.Pagination) ([]domain.Website, int, error) {
	return w.websiteRepository.ListWebsites(ctx, filters, pagination)
}

func (w WebsiteService) GetByStatus(ctx context.Context, enabled bool, pagination domain.Pagination) ([]domain.Website, error) {
	return w.websiteRepository.GetWebsitesByStatus(ctx, enabled, pagination)
}

func (w WebsiteService) UpdateLastCheck(ctx context.Context, websiteID uuid.UUID) error {
	site, err := w.websiteRepository.GetWebsiteByID(ctx, websiteID)
	if err != nil {
		return err
	}
	site.LastCheckAt = transform.ToPtr(time.Now())

	nextCheckAt, err := w.scheduler.NextRun(site.Cron, time.Now())
	if err != nil {
		return err
	}
	site.NextCheckAt = &nextCheckAt

	_, err = w.websiteRepository.UpdateWebsite(ctx, site)
	return err
}

func (w WebsiteService) UpdateStatus(ctx context.Context, id uuid.UUID, enabled bool) (domain.Website, error) {
	return w.websiteRepository.UpdateStatusWebsite(ctx, id, enabled)
}

func (w WebsiteService) Delete(ctx context.Context, id uuid.UUID) error {
	return w.websiteRepository.DeleteWebsite(ctx, id)
}

func checkMode(mode domain.Mode) error {
	switch mode {
	case domain.ModePlain:
		return nil
	default:
		return database.ErrModeNotCorrect
	}
}
