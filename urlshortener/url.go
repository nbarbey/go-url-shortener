package urlshortener

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"time"
)

type URL struct {
	*url.URL
	expiration *time.Time
}

var ErrInvalidURL = errors.New("invalid URL")

func MustNewURL(rawURL string, expiration *time.Time) URL {
	u, err := NewURL(rawURL, expiration)
	if err != nil {
		panic(fmt.Sprintf("unexpected error: `%s`", err))
	}
	return u
}

func NewURL(rawURL string, expiration *time.Time) (URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return URL{}, ErrInvalidURL
	}
	return URL{URL: u, expiration: expiration}, nil
}

func (u URL) Validate() error {
	if u.Scheme == "" {
		return ErrMissingScheme
	}
	if u.Hostname() == "" {
		return ErrMissingHostname
	}
	return nil
}

func (u URL) encode() string {
	payload := fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, u.Path)
	m := md5.Sum([]byte(payload))

	// base62 encoding from https://ucarion.com/go-base62
	var i big.Int
	i.SetBytes(m[:])
	return i.Text(62)
}

func (u URL) Shorten() (URL, error) {
	if err := u.Validate(); err != nil {
		return URL{}, err
	}
	shortenedPath := fmt.Sprintf("u/%s", u.encode())
	return URL{URL: &url.URL{Scheme: "https", Host: "localhost:8080", Path: shortenedPath}, expiration: u.expiration}, nil
}

func (u URL) Expiring() bool {
	return u.expiration != nil
}

func (u URL) ExpiredAt(t time.Time) bool {
	return u.Expiring() && u.expiration.Before(t)
}
