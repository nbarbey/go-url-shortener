package url_shortener

import (
	"errors"
)

var ErrNotFound = errors.New("URL not found")
var ErrMissingHostname = errors.New("missing hostname")
var ErrMissingScheme = errors.New("missing scheme")

func (a *Application) Shorten(rawURL string) (string, error) {
	u, err := NewURL(rawURL)
	if err != nil {
		return "", err
	}
	s, err := u.Shorten()
	if err != nil {
		return "", err
	}
	err = a.store.Save(rawURL, s.String())
	return s.String(), err
}

func (a *Application) Unshorten(rawURL string) (string, error) {
	u, err := NewURL(rawURL)
	if err != nil {
		return "", err
	}
	if err := u.Validate(); err != nil {
		return "", err
	}
	got, err := a.store.Get(rawURL)
	if err != nil {
		return "", ErrNotFound
	}
	return got, nil
}
