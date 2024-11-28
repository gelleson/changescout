package ent

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/check"
	"github.com/google/uuid"
	"time"
)

type CheckRepository struct {
	client *ent.Client
}

func NewCheckRepository(client *ent.Client) *CheckRepository {
	return &CheckRepository{
		client: client,
	}
}

func (r *CheckRepository) CreateCheck(ctx context.Context, check domain.Check) (domain.Check, error) {
	created, err := r.client.Check.Create().
		SetWebsiteID(check.WebsiteID).
		SetResult(check.Result).
		SetHasError(check.HasError).
		SetErrorMessage(check.ErrorMessage).
		SetHasDiff(check.HasChanges).
		SetDiffChange(check.DiffResult).
		SetCreatedAt(time.Now()).
		SetDiffChange(check.DiffResult).
		Save(ctx)
	if err != nil {
		return domain.Check{}, err
	}

	return mapCheckToEntity(created), nil
}

func (r *CheckRepository) GetCheckByID(ctx context.Context, id uuid.UUID) (domain.Check, error) {
	found, err := r.client.Check.Query().
		Where(check.ID(id)).
		Only(ctx)
	if err != nil {
		return domain.Check{}, err
	}

	return mapCheckToEntity(found), nil
}

func (r *CheckRepository) ClearChecksByWebsite(ctx context.Context, websiteID uuid.UUID) error {
	_, err := r.client.Check.Delete().
		Where(check.WebsiteID(websiteID)).
		Exec(ctx)
	return err
}

func (r *CheckRepository) ListChecks(ctx context.Context, filters database.CheckFilters, pagination domain.Pagination) ([]domain.Check, int, error) {
	query := r.client.Check.Query()

	if filters.WebsiteID != nil {
		query = query.Where(check.WebsiteID(*filters.WebsiteID))
	}
	if filters.FromDate != nil && !filters.FromDate.IsZero() {
		query = query.Where(check.CreatedAtGTE(*filters.FromDate))
	}
	if filters.ToDate != nil && !filters.ToDate.IsZero() {
		query = query.Where(check.CreatedAtLTE(*filters.ToDate))
	}

	// Get total count
	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination
	query = query.
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Order(ent.Desc(check.FieldCreatedAt))

	checks, err := query.All(ctx)
	if err != nil {
		return nil, 0, err
	}

	var result []domain.Check
	for _, c := range checks {
		result = append(result, mapCheckToEntity(c))
	}

	return result, total, nil
}

func (r *CheckRepository) UpdateCheck(ctx context.Context, check domain.Check) (domain.Check, error) {
	updated, err := r.client.Check.UpdateOneID(check.ID).
		SetWebsiteID(check.WebsiteID).
		SetResult(check.Result).
		Save(ctx)
	if err != nil {
		return domain.Check{}, err
	}

	return mapCheckToEntity(updated), nil
}

func (r *CheckRepository) GetLatestCheckByWebsite(ctx context.Context, websiteID uuid.UUID) (domain.Check, error) {
	found, err := r.client.Check.Query().
		Where(check.WebsiteID(websiteID)).
		Order(ent.Desc(check.FieldCreatedAt)).
		First(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return domain.Check{}, err
	}

	if err != nil && ent.IsNotFound(err) {
		return domain.Check{}, domain.ErrCheckNotFound
	}

	return mapCheckToEntity(found), nil
}

// Helper function to map ent.Check to domain.Check
func mapCheckToEntity(check *ent.Check) domain.Check {
	return domain.Check{
		ID:           check.ID,
		WebsiteID:    check.WebsiteID,
		DiffResult:   check.DiffChange,
		HasError:     check.HasError,
		ErrorMessage: check.ErrorMessage,
		HasChanges:   check.HasDiff,
		Result:       check.Result,
		CreatedAt:    check.CreatedAt,
	}
}
