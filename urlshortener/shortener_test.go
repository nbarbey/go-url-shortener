package urlshortener

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplicationShortenUnshortener(t *testing.T) {
	ShortenUnshortenerFromBuilder(t, func() ShortenUnshortener {
		return NewInMemoryApplication()
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
		_, err := app.Shorten("https:// ", nil)

		assert.ErrorContains(t, err, "invalid URL")
	})

	t.Run("invalid_url_missing_hostname", func(t *testing.T) {
		app := builder()
		_, err := app.Shorten("https://", nil)

		assert.ErrorIs(t, err, ErrMissingHostname)
	})

	t.Run("invalid_url_missing_scheme", func(t *testing.T) {
		app := builder()
		_, err := app.Shorten("toto.com", nil)

		assert.ErrorIs(t, err, ErrMissingScheme)
	})

	t.Run("ok_first_example", func(t *testing.T) {
		app := builder()
		u := "https://medium.com/leboncoin-tech-blog/seriously-you-should-be-having-fun-writing-software-at-work-fa92c7cd008c"
		shortenedURL, err := app.Shorten(u, nil)
		require.NoError(t, err)

		assert.Equal(t, "https://localhost:8080/u/6lHWylUzE7YYSRslbslMap", shortenedURL)

		gotURL, err := app.Unshorten(shortenedURL)
		require.NoError(t, err)

		assert.Equal(t, u, gotURL)
	})

	t.Run("ok_random_path", func(t *testing.T) {
		app := builder()
		url := "https://localhost/bla/bla/bla"
		shortenedURL, err := app.Shorten(url, nil)
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
	shortenedURL, err := app.Shorten(url, nil)
	require.NoError(t, err)

	gotURL, err := app.Unshorten(shortenedURL)
	require.NoError(t, err)
	return gotURL
}

func TestHTTPUnshorten_with_redirect(t *testing.T) {
	app := NewInMemoryApplication()
	testServer := httptest.NewServer(app.server.mux)
	client := NewHTTPClientFromResty(resty.NewWithClient(testServer.Client()).
		SetBaseURL(testServer.URL))

	rawURL := "https://developer.hashicorp.com/vault/tutorials/get-started/understand-static-dynamic-secrets"
	short, err := client.Shorten(rawURL, nil)
	require.NoError(t, err)

	assert.Equal(t, "https://localhost:8080/u/6Hgh0HxUDE0TQs8NYZDHtP", short)

	request := httptest.NewRequest("GET", short, nil)
	recorder := httptest.NewRecorder()
	app.server.mux.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusTemporaryRedirect, recorder.Code)
	assert.Equal(t, rawURL, recorder.Header().Get("Location"))
}

func TestHTTPUnshorten_counter(t *testing.T) {
	app := NewInMemoryApplication()
	testServer := httptest.NewServer(app.server.mux)
	client := NewHTTPClientFromResty(resty.NewWithClient(testServer.Client()).
		SetBaseURL(testServer.URL))

	short, err := client.Shorten("https://developer.hashicorp.com/vault/tutorials/get-started/understand-static-dynamic-secrets", nil)
	require.NoError(t, err)

	for _ = range 10 {
		handle(app, httptest.NewRequest("GET", short, nil))
	}

	request := httptest.NewRequest("GET", fmt.Sprintf("/count?url=%s", url.QueryEscape(short)), nil)
	recorder := handle(app, request)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, `{"count": 10}`, recorder.Body.String())
}

func handle(app *Application, request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	app.server.mux.ServeHTTP(recorder, request)
	return recorder
}

func TestShortenWithExpiration(t *testing.T) {
	app := NewInMemoryApplication()
	clock := clockwork.NewFakeClock()
	app.WithClock(clock)

	expirationTime := clock.Now().Add(1 * time.Hour)
	short, err := app.Shorten("https://developer.hashicorp.com/vault/tutorials/get-started/understand-static-dynamic-secrets", &expirationTime)
	require.NoError(t, err)

	_, err = app.Unshorten(short)
	require.NoError(t, err)

	clock.Advance(2 * time.Hour)
	_, err = app.Unshorten(short)
	require.ErrorIs(t, err, ErrExpired)
}
