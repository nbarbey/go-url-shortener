package urlshortener

import (
	"errors"
	"time"

	"github.com/jonboulle/clockwork"
)

var ErrNotFound = errors.New("URL not found")
var ErrMissingHostname = errors.New("missing hostname")
var ErrMissingScheme = errors.New("missing scheme")
var ErrExpired = errors.New("URL expired")

type Shortener interface {
	Shorten(rawURL string, expiration *time.Time) (string, error)
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
	clock clockwork.Clock
}

func NewUsecase(store Storer) *Usecase {
	return &Usecase{
		store: store,
		clock: clockwork.NewRealClock(),
	}
}

func (u *Usecase) WithClock(clock clockwork.Clock) {
	u.clock = clock
}

func (c *Usecase) Shorten(rawURL string, expiration *time.Time) (string, error) {
	u, err := NewURL(rawURL, expiration)
	if err != nil {
		return "", err
	}
	s, err := u.Shorten()
	if err != nil {
		return "", err
	}
	err = c.store.Save(rawURL, s.String(), expiration)
	return s.String(), err
}

func (c *Usecase) Unshorten(rawURL string) (string, error) {
	u, err := NewURL(rawURL, nil)
	if err != nil {
		return "", err
	}
	if err := u.Validate(); err != nil {
		return "", err
	}
	storedURL, err := c.store.Get(rawURL)
	if err != nil {
		return "", ErrNotFound
	}
	location, err := time.LoadLocation("Local")
	if err != nil {
		return "", err
	}
	if storedURL.ExpiredAt(c.clock.Now().In(location)) {
		return "", ErrExpired
	}
	return storedURL.String(), nil
}

type CountingUsecase struct {
	*Usecase
	countStore CountStorer
}

func NewCountingUsecase(store Storer, countStore CountStorer) *CountingUsecase {
	return &CountingUsecase{Usecase: NewUsecase(store), countStore: countStore}
}

func (c *CountingUsecase) Unshorten(rawURL string) (string, error) {
	got, err := c.Usecase.Unshorten(rawURL)
	if err == nil {
		_ = c.countStore.Increment(rawURL)
	}

	return got, err
}
