package url_shortener

import (
	"errors"
	"fmt"
	"net/url"
)

var ErrNotFound = errors.New("URL not found")
var ErrMissingHostname = errors.New("missing hostname")

func Unshorten(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}
	if u.Hostname() == "" {
		return "", ErrMissingHostname
	}
	return "", ErrNotFound
}
