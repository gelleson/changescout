package browser

import (
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/utils/transform"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"time"
)

type BrowserService struct {
	browser *rod.Browser
}

func New(opts ...func(*options) *options) *BrowserService {
	opt := transform.Pipe[*options](
		&options{},
		opts...,
	)
	var b *rod.Browser
	if opt.managedInstanceURL != nil {
		l := launcher.MustNewManaged(*opt.managedInstanceURL)

		b = rod.New().
			ControlURL(l.MustLaunch()).
			MustConnect()
	}
	if opt.managedInstanceURL == nil {
		b = rod.New().MustConnect()
	}

	return &BrowserService{
		browser: b,
	}
}

func (b BrowserService) Request(site domain.Website) ([]byte, error) {
	p, err := b.browser.Page(proto.TargetCreateTarget{URL: site.URL})
	if err != nil {
		return nil, err
	}

	var defaultTimeout = time.Second * 10
	if site.Setting.RenderedOption.WaitForTimeout != nil {
		defaultTimeout = time.Second * time.Duration(*site.Setting.RenderedOption.WaitForTimeout)
	}

	if err := p.WaitIdle(defaultTimeout); err != nil {
		return nil, err
	}

	if site.Setting.RenderedOption.WaitForTimeout != nil {
		time.Sleep(time.Second * time.Duration(*site.Setting.RenderedOption.WaitForTimeout))
	}

	html, err := p.HTML()
	if err != nil {
		return nil, err
	}

	return []byte(html), nil
}
