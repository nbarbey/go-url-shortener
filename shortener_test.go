package url_shortener

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShortener_not_found(t *testing.T) {
	_, err := Unshorten("https://localhost/abcd1234")

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestShortener_invalid_url_invalid_character(t *testing.T) {
	_, err := Unshorten("https:// ")

	assert.ErrorContains(t, err, "invalid character")
}

func TestShortener_invalid_url_missing_hostname(t *testing.T) {
	_, err := Unshorten("https://")

	assert.ErrorIs(t, err, ErrMissingHostname)
}
