package url_shortener

import (
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

// base62 encoding from https://ucarion.com/go-base62
func (u URL) encode() string {
	var i big.Int
	payload := fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, u.Path)
	i.SetBytes([]byte(payload[:]))
	return i.Text(62)
}

func (u URL) Shorten() (URL, error) {
	if err := u.Validate(); err != nil {
		return URL{}, err
	}
	shortenedPath := u.encode()
	return URL{URL: &url.URL{Scheme: "https", Host: "localhost", Path: shortenedPath}}, nil
}
