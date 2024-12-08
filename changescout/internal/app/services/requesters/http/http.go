package http

import (
	"errors"
	"github.com/gelleson/changescout/changescout/internal/domain"
	"io"
	"net/http"
)

//go:generate mockery --name Doer
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type HttpService struct {
	doer Doer
}

func New(doer Doer) *HttpService {
	return &HttpService{
		doer: doer,
	}
}

func (h HttpService) Request(site domain.Website) ([]byte, error) {
	req, err := http.NewRequest(site.Setting.Method, site.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = site.Setting.Headers
	req.Header.Set("User-Agent", site.Setting.UserAgent)
	req.Header.Set("Referer", site.Setting.Referer)

	resp, err := h.doer.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New("bad status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
