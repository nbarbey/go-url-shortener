package urlshortener

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApplicationShortenUnshortener(t *testing.T) {
	ShortenUnshortenerFromBuilder(t, func() ShortenUnshortener {
		return NewApplication()
	})
}

func ShortenUnshortenerFromBuilder(t *testing.T, builder func() ShortenUnshortener) {
	t.Run("not_found", func(t *testing.T) {
		app := builder()
		_, err := app.Unshorten("https://localhost/abcd1234")

		assert.ErrorIs(t, err, ErrNotFound)
	})

	t.Run("invalid_url_invalid_character", func(t *testing.T) {
		app := builder()
		_, err := app.Unshorten("https:// ")

		assert.ErrorContains(t, err, "invalid URL")
	})

	t.Run("invalid_url_missing_hostname", func(t *testing.T) {
		app := builder()
		_, err := app.Unshorten("https://")

		assert.ErrorIs(t, err, ErrMissingHostname)
	})

	t.Run("invalid_url_missing_scheme", func(t *testing.T) {
		app := builder()
		_, err := app.Unshorten("toto.com")

		assert.ErrorIs(t, err, ErrMissingScheme)
	})

	t.Run("invalid_url_invalid_character", func(t *testing.T) {
		app := builder()
		_, err := app.Shorten("https:// ")

		assert.ErrorContains(t, err, "invalid URL")
	})

	t.Run("invalid_url_missing_hostname", func(t *testing.T) {
		app := builder()
		_, err := app.Shorten("https://")

		assert.ErrorIs(t, err, ErrMissingHostname)
	})

	t.Run("invalid_url_missing_scheme", func(t *testing.T) {
		app := builder()
		_, err := app.Shorten("toto.com")

		assert.ErrorIs(t, err, ErrMissingScheme)
	})

	t.Run("ok_first_example", func(t *testing.T) {
		app := builder()
		url := "https://medium.com/equify-tech/the-three-fundamental-stages-of-an-engineering-career-54dac732fc74"
		shortenedURL, err := app.Shorten(url)
		require.NoError(t, err)

		gotURL, err := app.Unshorten(shortenedURL)
		require.NoError(t, err)

		assert.Equal(t, url, gotURL)
	})

	t.Run("ok_random_path", func(t *testing.T) {
		app := builder()
		url := "https://localhost/bla/bla/bla"
		shortenedURL, err := app.Shorten(url)
		require.NoError(t, err)

		gotURL, err := app.Unshorten(shortenedURL)
		require.NoError(t, err)

		assert.Equal(t, url, gotURL)
	})

	t.Run("ok_two_paths", func(t *testing.T) {
		app := builder()
		url1 := "https://foobar/first"
		gotURL1 := shortenUnshorten(t, app, url1)
		assert.Equal(t, url1, gotURL1)

		url2 := "https://foobar/second"
		gotURL2 := shortenUnshorten(t, app, url2)
		assert.Equal(t, url2, gotURL2)
	})
}

func shortenUnshorten(t *testing.T, app ShortenUnshortener, url string) string {
	shortenedURL, err := app.Shorten(url)
	require.NoError(t, err)

	gotURL, err := app.Unshorten(shortenedURL)
	require.NoError(t, err)
	return gotURL
}
