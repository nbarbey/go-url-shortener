package url_shortener

import (
	"errors"
	"fmt"
	"net/url"
)

var ErrNotFound = errors.New("URL not found")
var ErrMissingHostname = errors.New("missing hostname")
var ErrMissingScheme = errors.New("missing scheme")

func (a *Application) Shorten(rawURL string) (string, error) {
	if err := validateURL(rawURL); err != nil {
		return "", err
	}
	s := "https://localhost/hardcoded"
	err := a.store.Save(rawURL, s)
	return s, err
}

func (a *Application) Unshorten(rawURL string) (string, error) {
	if err := validateURL(rawURL); err != nil {
		return "", err
	}
	u, err := a.store.Get(rawURL)
	if err != nil {
		return "", ErrNotFound
	}
	return u, nil
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
