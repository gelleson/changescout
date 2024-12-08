//go:build integration

package browser

import (
	"github.com/gelleson/changescout/changescout/internal/domain"
	"github.com/gelleson/changescout/changescout/internal/utils/transform"
	"github.com/go-rod/rod"
	"github.com/stretchr/testify/suite"
	"testing"
)

type BrowserServiceTestSuite struct {
	suite.Suite
	browserService *BrowserService
}

func (suite *BrowserServiceTestSuite) SetupTest() {
	suite.browserService = &BrowserService{
		browser: rod.New().MustConnect(),
	}
}

func (suite *BrowserServiceTestSuite) TestRequestSuccess() {
	site := domain.Website{
		URL: "https://go-rod.github.io/#/network/README?id=throttling",
		Setting: domain.Setting{
			RenderedOption: domain.RenderedOption{
				WaitForTimeout: transform.ToPtr(1),
			},
		},
	}

	html, err := suite.browserService.Request(site)
	suite.NoError(err)
	suite.NotEmpty(html)
}

func TestBrowserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BrowserServiceTestSuite))
}
