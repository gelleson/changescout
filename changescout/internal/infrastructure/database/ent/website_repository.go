package ent

import (
	"context"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/predicate"
	"github.com/gelleson/changescout/changescout/internal/infrastructure/database/ent/ent/website"
	"github.com/gelleson/changescout/changescout/internal/utils/transform"
	"github.com/google/uuid"
	"time"
)

type WebsiteRepository struct {
	client *ent.Client
}

func NewWebsiteRepository(client *ent.Client) *WebsiteRepository {
	return &WebsiteRepository{
		client: client,
	}
}

func (w WebsiteRepository) CreateWebsite(ctx context.Context, website domain.Website) (domain.Website, error) {
	site, err := w.client.
		Website.
		Create().
		SetName(website.Name).
		SetURL(website.URL).
		SetEnabled(website.Enabled).
		SetMode(string(website.Mode)).
		SetCron(website.Cron).
		SetSetting(&website.Setting).
		SetUserID(website.UserID).
		SetNillableLastCheckAt(website.LastCheckAt).
		SetNextCheckAt(transform.ToValueOrDefault(website.NextCheckAt, time.Time{})).
		Save(ctx)

	if err != nil {
		return domain.Website{}, err
	}

	return entToWebsite(site), nil
}

func (w WebsiteRepository) GetWebsiteByID(ctx context.Context, id uuid.UUID) (domain.Website, error) {
	entity, err := w.client.
		Website.
		Query().
		Where(website.ID(id)).
		Only(ctx)

	if err != nil {
		return domain.Website{}, err
	}

	return entToWebsite(entity), nil
}

func (w WebsiteRepository) GetWebsiteByURL(ctx context.Context, url string) (domain.Website, error) {
	entity, err := w.client.
		Website.
		Query().
		Where(website.URL(url)).
		Only(ctx)

	if err != nil {
		return domain.Website{}, err
	}

	return entToWebsite(entity), nil
}

func (w WebsiteRepository) GetWebsiteByUserID(ctx context.Context, userID uuid.UUID, pagination domain.Pagination) ([]domain.Website, error) {
	entities, err := w.client.
		Website.
		Query().
		Where(website.UserID(userID)).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		All(ctx)

	if err != nil {
		return nil, err
	}

	return entToWebsites(entities), nil
}

func (w WebsiteRepository) UpdateWebsite(ctx context.Context, website domain.Website) (domain.Website, error) {
	_website := w.client.Website.UpdateOneID(website.ID)
	_mutation := _website.Mutation()

	if website.Name != "" {
		_mutation.SetName(website.Name)
	}
	if website.URL != "" {
		_mutation.SetURL(website.URL)
	}
	if website.Enabled != false {
		_mutation.SetEnabled(website.Enabled)
	}
	if website.Mode != "" {
		mode, err := modeDomainToEnt(website.Mode)
		if err != nil {
			return domain.Website{}, err
		}
		_mutation.SetMode(mode)
	}
	if website.LastCheckAt != nil {
		_mutation.SetLastCheckAt(*website.LastCheckAt)
	}
	if website.Setting.Method != "" {
		_mutation.SetSetting(&website.Setting)
	}
	if website.UserID != uuid.Nil {
		_mutation.SetUserID(website.UserID)
	}
	if website.NextCheckAt != nil {
		_mutation.SetNextCheckAt(*website.NextCheckAt)
	}
	if website.Cron != "" {
		_mutation.SetCron(website.Cron)
	}

	site, err := _website.Save(ctx)
	if err != nil {
		return domain.Website{}, err
	}

	return entToWebsite(site), err
}

func (w WebsiteRepository) UpdateStatusWebsite(ctx context.Context, siteID uuid.UUID, enabled bool) (domain.Website, error) {
	_website := w.client.Website.UpdateOneID(siteID)
	_mutation := _website.Mutation()

	_mutation.SetEnabled(enabled)

	site, err := _website.Save(ctx)
	if err != nil {
		return domain.Website{}, err
	}

	return entToWebsite(site), err
}

func modeDomainToEnt(mode domain.Mode) (string, error) {
	switch mode {
	case domain.ModePlain:
		return string(domain.ModePlain), nil
	default:
		return "", database.ErrModeNotCorrect
	}
}

func entToWebsite(website *ent.Website) domain.Website {
	return domain.Website{
		ID:          website.ID,
		Name:        website.Name,
		URL:         website.URL,
		Enabled:     website.Enabled,
		Cron:        website.Cron,
		Mode:        domain.Mode(website.Mode),
		Setting:     transform.ToValueOrDefault(website.Setting, domain.Setting{}),
		UserID:      website.UserID,
		CreatedAt:   website.CreatedAt,
		UpdatedAt:   website.UpdatedAt,
		NextCheckAt: transform.ToPtr(website.NextCheckAt),
		LastCheckAt: transform.ToPtr(website.LastCheckAt),
	}
}

func entToWebsites(websites []*ent.Website) []domain.Website {
	domains := make([]domain.Website, len(websites))
	for i, site := range websites {
		domains[i] = entToWebsite(site)
	}
	return domains
}

func websitesToEnt(websites []domain.Website) []*ent.Website {
	domains := make([]*ent.Website, len(websites))
	for i, site := range websites {
		domains[i] = websiteToEnt(site)
	}
	return domains
}

func websiteToEnt(website domain.Website) *ent.Website {

	return &ent.Website{
		Name:        website.Name,
		URL:         website.URL,
		Enabled:     website.Enabled,
		Mode:        string(website.Mode),
		Cron:        website.Cron,
		Setting:     &website.Setting,
		UserID:      website.UserID,
		LastCheckAt: transform.ToValueOrDefault(website.LastCheckAt, time.Time{}),
		NextCheckAt: transform.ToValueOrDefault(website.NextCheckAt, time.Time{}),
		CreatedAt:   website.CreatedAt,
		UpdatedAt:   website.UpdatedAt,
	}
}

func (w WebsiteRepository) ListWebsites(ctx context.Context, filters database.WebsiteFilters, pagination domain.Pagination) ([]domain.Website, int, error) {
	predicates := make([]predicate.Website, 0)
	if filters.UserID != nil {
		predicates = append(predicates, website.UserID(*filters.UserID))
	}
	if filters.Enabled != nil {
		predicates = append(predicates, website.Enabled(*filters.Enabled))
	}
	if filters.Mode != nil {
		predicates = append(predicates, website.Mode(string(*filters.Mode)))
	}
	if filters.URLQuery != nil {
		predicates = append(predicates, website.URLContains(*filters.URLQuery))
	}
	if filters.FromDate != nil {
		predicates = append(predicates, website.LastCheckAtGTE(*filters.FromDate))
	}
	if filters.ToDate != nil {
		predicates = append(predicates, website.LastCheckAtLTE(*filters.ToDate))
	}

	entities, err := w.client.
		Website.
		Query().
		Where(predicates...).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		All(ctx)

	if err != nil {
		return nil, 0, err
	}

	total, err := w.client.
		Website.
		Query().
		Where(predicates...).
		Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return entToWebsites(entities), int(total), nil
}

func (w WebsiteRepository) GetWebsitesByStatus(ctx context.Context, enabled bool, pagination domain.Pagination) ([]domain.Website, error) {
	predicates := make([]predicate.Website, 0)
	if enabled {
		predicates = append(predicates, website.EnabledEQ(true))
	} else {
		predicates = append(predicates, website.EnabledEQ(false))
	}

	entities, err := w.client.
		Website.
		Query().
		Where(predicates...).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		All(ctx)

	if err != nil {
		return nil, err
	}

	return entToWebsites(entities), nil
}

func (w WebsiteRepository) SearchWebsites(ctx context.Context, query string, pagination domain.Pagination) ([]domain.Website, error) {
	predicates := make([]predicate.Website, 0)
	if query != "" {
		predicates = append(predicates, website.Or(
			website.NameContains(query),
			website.URLContains(query),
		))
	}

	entities, err := w.client.
		Website.
		Query().
		Where(predicates...).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		All(ctx)

	if err != nil {
		return nil, err
	}

	return entToWebsites(entities), nil
}

func (w WebsiteRepository) GetWebsitesDueForCheck(ctx context.Context, pagination domain.Pagination) ([]domain.Website, error) {
	now := time.Now()

	entities, err := w.client.
		Website.
		Query().
		Where(
			website.And(
				website.NextCheckAtLTE(now),
				website.EnabledEQ(true),
			),
		).
		Order(ent.Asc(website.FieldNextCheckAt)). // Add ordering to make results consistent
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		All(ctx)

	if err != nil {
		return nil, err
	}

	return entToWebsites(entities), nil
}

func (w WebsiteRepository) DeleteWebsite(ctx context.Context, id uuid.UUID) error {
	return w.client.
		Website.
		DeleteOneID(id).
		Exec(ctx)
}

var _ database.WebsiteRepository = (*WebsiteRepository)(nil)
