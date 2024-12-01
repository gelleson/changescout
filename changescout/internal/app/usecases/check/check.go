package check

import (
	"context"
	"fmt"
	"github.com/gelleson/changescout/changescout/internal/app/services/diff"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/pkg/processors"
	"github.com/google/uuid"
	"net/http"
)

//go:generate mockery --name Doer
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

//go:generate mockery --name WebsiteService
type WebsiteService interface {
	GetByID(ctx context.Context, id uuid.UUID) (domain.Website, error)
}

//go:generate mockery --name HttpService
type HttpService interface {
	Request(site domain.Website) ([]byte, error)
}

//go:generate mockery --name DiffService
type DiffService interface {
	Compare(previous, current []byte) (diff.Result, error)
}

//go:generate mockery --name DBService
type DBService interface {
	GetLatestCheckByWebsite(ctx context.Context, websiteID uuid.UUID) (domain.Check, error)
	CreateCheck(ctx context.Context, check domain.Check) (domain.Check, error)
}

type UseCase struct {
	websiteService WebsiteService
	httpService    HttpService
	checkService   DBService
	diffService    DiffService
}

func NewUseCase(
	websiteService WebsiteService,
	httpService HttpService,
	checkService DBService,
	diffService DiffService,
) *UseCase {
	return &UseCase{
		websiteService: websiteService,
		httpService:    httpService,
		checkService:   checkService,
		diffService:    diffService,
	}
}

func (u UseCase) View(ctx context.Context, websiteID uuid.UUID) ([]byte, error) {
	site, err := u.websiteService.GetByID(ctx, websiteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get website: %w", err)
	}

	// Make HTTP request
	body, err := u.makeRequestAndHandleError(ctx, site)
	if err != nil {
		return nil, err
	}

	processor := processors.New(
		processors.NewHTMLProcessor(site.Setting),
		processors.NewJSONPathProcessor(site.Setting),
		processors.NewDeduplicationProcessor(site.Setting),
		processors.NewTrimProcessor(site.Setting),
		processors.NewSortProcessor(site.Setting),
	)

	body = processor.Run(body)

	return body, nil
}

func (u UseCase) Check(ctx context.Context, websiteID uuid.UUID) (domain.CheckResult, error) {
	// Get website details
	site, err := u.websiteService.GetByID(ctx, websiteID)
	if err != nil {
		return domain.CheckResult{}, fmt.Errorf("failed to get website: %w", err)
	}

	// Make HTTP request
	body, err := u.View(ctx, site.ID)
	if err != nil {
		return domain.CheckResult{}, err
	}

	// Compare with previous check
	diffResult, prevResult, err := u.compareWithPreviousCheck(ctx, site.ID, body)
	if err != nil {
		return domain.CheckResult{}, err
	}

	// If no changes, return early
	if !diffResult.HasChanges {
		return domain.CheckResult{}, nil
	}

	// Create new check record
	if err := u.createSuccessfulCheck(ctx, site.ID, body, diffResult); err != nil {
		return domain.CheckResult{}, err
	}

	return domain.CheckResult{
		OldValue:   prevResult,
		NewValue:   body,
		HasChanges: diffResult.HasChanges,
		Check:      diffResult,
	}, nil
}

// makeRequestAndHandleError handles the HTTP request and records any errors
func (u UseCase) makeRequestAndHandleError(ctx context.Context, site domain.Website) ([]byte, error) {
	body, err := u.httpService.Request(site)
	if err != nil {
		if createErr := u.createFailedCheck(ctx, site.ID, err); createErr != nil {
			return nil, fmt.Errorf("failed to create error check: %w (original error: %v)", createErr, err)
		}
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	return body, nil
}

// compareWithPreviousCheck gets the latest check and compares results
func (u UseCase) compareWithPreviousCheck(ctx context.Context, websiteID uuid.UUID, currentBody []byte) (diff.Result, []byte, error) {
	latestCheck, err := u.checkService.GetLatestCheckByWebsite(ctx, websiteID)
	if err != nil && !domain.IsErrCheckNotFound(err) {
		return diff.Result{}, nil, fmt.Errorf("failed to get latest check: %w", err)
	}

	var prevResult []byte
	if latestCheck.Result != nil {
		prevResult = latestCheck.Result
	}

	diffResult, err := u.diffService.Compare(prevResult, currentBody)
	if err != nil {
		return diff.Result{}, nil, fmt.Errorf("failed to compare results: %w", err)
	}

	return diffResult, prevResult, nil
}

// createFailedCheck creates a check record for a failed request
func (u UseCase) createFailedCheck(ctx context.Context, websiteID uuid.UUID, requestError error) error {
	_, err := u.checkService.CreateCheck(ctx, domain.Check{
		WebsiteID:    websiteID,
		Result:       nil,
		DiffResult:   &diff.Result{},
		HasChanges:   true,
		HasError:     true,
		ErrorMessage: requestError.Error(),
	})
	return err
}

// createSuccessfulCheck creates a check record for a successful comparison
func (u UseCase) createSuccessfulCheck(ctx context.Context, websiteID uuid.UUID, body []byte, diffResult diff.Result) error {
	_, err := u.checkService.CreateCheck(ctx, domain.Check{
		WebsiteID:  websiteID,
		Result:     body,
		DiffResult: &diffResult,
		HasChanges: true,
		HasError:   false,
	})
	return err
}
