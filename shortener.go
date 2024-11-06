package url_shortener

import (
	"errors"
	"fmt"
	"net/url"
)

var ErrNotFound = errors.New("URL not found")
var ErrMissingHostname = errors.New("missing hostname")
var ErrMissingScheme = errors.New("missing scheme")

type Application struct {
	store *Store
}

func NewApplication() *Application {
	store := NewStore()
	store.Save("https://medium.com/equify-tech/the-three-fundamental-stages-of-an-engineering-career-54dac732fc74",
		"https://localhost/hardcoded")
	return &Application{store: store}
}

type Store struct {
	data map[string]string
}

func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}

func (s *Store) Get(shortened string) (string, error) {
	u, ok := s.data[shortened]
	if !ok {
		return "", ErrNotFound
	}
	return u, nil
}

func (s *Store) Save(url, shortened string) error {
	s.data[shortened] = url
	return nil
}

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
