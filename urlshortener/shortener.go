package urlshortener

import (
	"errors"
)

var ErrNotFound = errors.New("URL not found")
var ErrMissingHostname = errors.New("missing hostname")
var ErrMissingScheme = errors.New("missing scheme")

type Shortener interface {
	Shorten(rawURL string) (string, error)
}

type Unshortener interface {
	Unshorten(rawURL string) (string, error)
}

type ShortenUnshortener interface {
	Shortener
	Unshortener
}

type Usecase struct {
	store Storer
}

func NewUsecase(store Storer) *Usecase {
	return &Usecase{store: store}
}

func (c *Usecase) Shorten(rawURL string) (string, error) {
	u, err := NewURL(rawURL)
	if err != nil {
		return "", err
	}
	s, err := u.Shorten()
	if err != nil {
		return "", err
	}
	err = c.store.Save(rawURL, s.String())
	return s.String(), err
}

func (c *Usecase) Unshorten(rawURL string) (string, error) {
	u, err := NewURL(rawURL)
	if err != nil {
		return "", err
	}
	if err := u.Validate(); err != nil {
		return "", err
	}
	got, err := c.store.Get(rawURL)
	if err != nil {
		return "", ErrNotFound
	}
	return got, nil
}
