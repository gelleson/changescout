package services

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/google/uuid"
)

//go:generate mockery --name CheckRepository
type CheckRepository interface {
	database.CheckRepository
}

type CheckService struct {
	repository CheckRepository
}

func NewCheckService(repository CheckRepository) *CheckService {
	return &CheckService{
		repository: repository,
	}
}

func (s CheckService) GetLatestCheckByWebsite(ctx context.Context, websiteID uuid.UUID) (domain.Check, error) {
	return s.repository.GetLatestCheckByWebsite(ctx, websiteID)
}

func (s CheckService) CreateCheck(ctx context.Context, check domain.Check) (domain.Check, error) {
	return s.repository.CreateCheck(ctx, check)
}

func (s CheckService) ClearChecksByWebsite(ctx context.Context, websiteID uuid.UUID) error {
	return s.repository.ClearChecksByWebsite(ctx, websiteID)
}

func (s CheckService) ListChecks(ctx context.Context, filters database.CheckFilters, pagination domain.Pagination) ([]domain.Check, int, error) {
	return s.repository.ListChecks(ctx, filters, pagination)
}

func (s CheckService) UpdateCheck(ctx context.Context, check domain.Check) (domain.Check, error) {
	return s.repository.UpdateCheck(ctx, check)
}

func (s CheckService) GetCheckByID(ctx context.Context, id uuid.UUID) (domain.Check, error) {
	return s.repository.GetCheckByID(ctx, id)
}

var _ database.CheckRepository = (*CheckService)(nil)
