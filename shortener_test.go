package url_shortener

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnshorten_not_found(t *testing.T) {
	app := NewApplication()
	_, err := app.Unshorten("https://localhost/abcd1234")

	assert.ErrorIs(t, err, ErrNotFound)
}

func TestUnshorten_invalid_url_invalid_character(t *testing.T) {
	app := NewApplication()
	_, err := app.Unshorten("https:// ")

	assert.ErrorContains(t, err, "invalid character")
}

func TestUnshorten_invalid_url_missing_hostname(t *testing.T) {
	app := NewApplication()
	_, err := app.Unshorten("https://")

	assert.ErrorIs(t, err, ErrMissingHostname)
}

func TestUnshorten_invalid_url_missing_scheme(t *testing.T) {
	app := NewApplication()
	_, err := app.Unshorten("toto.com")

	assert.ErrorIs(t, err, ErrMissingScheme)
}

func TestShorten_invalid_url_invalid_character(t *testing.T) {
	app := NewApplication()
	_, err := app.Shorten("https:// ")

	assert.ErrorContains(t, err, "invalid character")
}

func TestShorten_invalid_url_missing_hostname(t *testing.T) {
	app := NewApplication()
	_, err := app.Shorten("https://")

	assert.ErrorIs(t, err, ErrMissingHostname)
}

func TestShorten_invalid_url_missing_scheme(t *testing.T) {
	app := NewApplication()
	_, err := app.Shorten("toto.com")

	assert.ErrorIs(t, err, ErrMissingScheme)
}

func TestShortenAndUnshorten_ok_first_example(t *testing.T) {
	app := NewApplication()
	url := "https://medium.com/equify-tech/the-three-fundamental-stages-of-an-engineering-career-54dac732fc74"
	shortenedURL, err := app.Shorten(url)
	require.NoError(t, err)

	gotURL, err := app.Unshorten(shortenedURL)
	require.NoError(t, err)

	assert.Equal(t, url, gotURL)
}
