package url_shortener

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnshorten_not_found(t *testing.T) {
	_, err := Unshorten("https://localhost/abcd1234")

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestUnshorten_invalid_url_invalid_character(t *testing.T) {
	_, err := Unshorten("https:// ")

	assert.ErrorContains(t, err, "invalid character")
}

func TestUnshorten_invalid_url_missing_hostname(t *testing.T) {
	_, err := Unshorten("https://")

	assert.ErrorIs(t, err, ErrMissingHostname)
}

func TestUnshorten_invalid_url_missing_scheme(t *testing.T) {
	_, err := Unshorten("toto.com")

	assert.ErrorIs(t, err, ErrMissingScheme)
}

func TestShorten_invalid_url_invalid_character(t *testing.T) {
	_, err := Shorten("https:// ")

	assert.ErrorContains(t, err, "invalid character")
}

func TestShorten_invalid_url_missing_hostname(t *testing.T) {
	_, err := Shorten("https://")

	assert.ErrorIs(t, err, ErrMissingHostname)
}

func TestShorten_invalid_url_missing_scheme(t *testing.T) {
	_, err := Shorten("toto.com")

	assert.ErrorIs(t, err, ErrMissingScheme)
}
