package requesters

import (
	"github.com/gelleson/changescout/changescout/internal/app/services/requesters/browser"
	httprequesters "github.com/gelleson/changescout/changescout/internal/app/services/requesters/http"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"net/http"
)

type Provider interface {
	Request(domain.Website) ([]byte, error)
}
type Providers map[domain.Mode]Provider

type Requester struct {
	providers Providers
}

type BrowserOption struct {
	Enable             bool
	ManagedInstanceURL *string
}

type Options struct {
	Browser BrowserOption
}

func New(opt Options) *Requester {
	return &Requester{
		providers: Providers{
			domain.ModePlain: httprequesters.New(http.DefaultClient),
			domain.ModeRenderer: optional[Provider](
				opt.Browser.Enable,
				func() Provider {
					return browser.New(
						browser.WithManagedInstanceURL(opt.Browser.ManagedInstanceURL),
					)
				},
				httprequesters.New(http.DefaultClient),
			),
		},
	}
}

func (r *Requester) Request(site domain.Website) ([]byte, error) {
	return r.providers[site.Mode].Request(site)
}

func optional[T any](enabled bool, builder func() T, fallback T) T {
	if enabled {
		return builder()
	}
	return fallback
}
