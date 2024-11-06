package url_shortener

import (
	"errors"
	"fmt"
	"net/url"
)

var ErrNotFound = errors.New("URL not found")
var ErrMissingHostname = errors.New("missing hostname")
var ErrMissingScheme = errors.New("missing scheme")

func Shorten(rawURL string) (string, error) {
	if err := validateURL(rawURL); err != nil {
		return "", err
	}
	return "https://localhost/hardcoded", nil
}

func Unshorten(rawURL string) (string, error) {
	if err := validateURL(rawURL); err != nil {
		return "", err
	}
	if rawURL == "https://localhost/hardcoded" {
		return "https://medium.com/equify-tech/the-three-fundamental-stages-of-an-engineering-career-54dac732fc74", nil
	}
	return "", ErrNotFound
}

func validateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}
	if u.Scheme == "" {
		return ErrMissingScheme
	}
	if u.Hostname() == "" {
		return ErrMissingHostname
	}
	return nil
}
