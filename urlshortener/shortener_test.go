package urlshortener

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
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

		assert.Equal(t, "https://localhost:8080/u/1oPzkR9KEQU5LZniKkpIub", shortenedURL)

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

func TestHTTPUnshorten_with_redirect(t *testing.T) {
	app := NewApplication()
	testServer := httptest.NewServer(app.server.mux)
	client := NewHTTPClientFromResty(resty.NewWithClient(testServer.Client()).
		SetBaseURL(testServer.URL))

	rawURL := "https://developer.hashicorp.com/vault/tutorials/get-started/understand-static-dynamic-secrets"
	short, err := client.Shorten(rawURL)
	require.NoError(t, err)

	assert.Equal(t, "https://localhost:8080/u/6Hgh0HxUDE0TQs8NYZDHtP", short)

	request := httptest.NewRequest("GET", short, nil)
	recorder := httptest.NewRecorder()
	app.server.mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	assert.Equal(t, rawURL, recorder.Header().Get("Location"))
}
