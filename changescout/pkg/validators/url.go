package validators

import (
	"net/url"
)

func IsValidURL(u string) bool {
	_url, err := url.Parse(u)
	if err != nil {
		return false
	}
	return _url.Scheme != "" && _url.Host != ""
}
