package urlshortener

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/big"
	"net/url"
)

type URL struct {
	*url.URL
}

var ErrInvalidURL = errors.New("invalid URL")

func NewURL(rawURL string) (URL, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return URL{}, ErrInvalidURL
	}
	return URL{URL: u}, nil
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
	shortenedPath := fmt.Sprintf("unshorten/%s", u.encode())
	return URL{URL: &url.URL{Scheme: "https", Host: "localhost", Path: shortenedPath}}, nil
}
