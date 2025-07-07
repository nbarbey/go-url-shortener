package urlshortener

import (
	"github.com/goccha/logging/restylog"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPShortenUnshortener(t *testing.T) {
	ShortenUnshortenerFromBuilder(t, func() ShortenUnshortener {
		app := NewInMemoryApplication()
		testServer := httptest.NewServer(app.server.mux)
		return NewHTTPClientFromResty(resty.NewWithClient(testServer.Client()).
			SetBaseURL(testServer.URL))
	})

	t.Run("URL not found", func(t *testing.T) {
		app := NewInMemoryApplication()
		request := httptest.NewRequest(http.MethodGet, "/unshorten?url=https%3A%2F%2Flocalhost%2Fabcd1234", nil)
		recorder := httptest.NewRecorder()

		app.server.mux.ServeHTTP(recorder, request)

		result := recorder.Result()
		data, err := io.ReadAll(result.Body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
		assert.Equal(t, `{"error": "URL not found"}`, string(data))
	})
	t.Run("missing hostname", func(t *testing.T) {
		app := NewInMemoryApplication()
		request := httptest.NewRequest(http.MethodGet, "/unshorten?url=https%3A//", nil)
		recorder := httptest.NewRecorder()

		app.server.mux.ServeHTTP(recorder, request)

		result := recorder.Result()
		data, err := io.ReadAll(result.Body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, result.StatusCode)
		assert.Equal(t, `{"error": "missing hostname"}`, string(data))
	})
}

func TestHTTPShortenWithExpiration(t *testing.T) {
	app := NewInMemoryApplication()
	testServer := httptest.NewServer(app.server.mux)
	restyClient := resty.NewWithClient(testServer.Client())
	restyClient.SetLogger(&restylog.Logger{})
	client := NewHTTPClientFromResty(restyClient.SetBaseURL(testServer.URL))
	clock := clockwork.NewFakeClock()
	app.WithClock(clock)

	location, err := time.LoadLocation("Local")
	require.NoError(t, err)
	expirationTime := clock.Now().In(location).Add(1 * time.Hour)
	short, err := client.Shorten("https://developer.hashicorp.com/vault/tutorials/get-started/understand-static-dynamic-secrets", &expirationTime)
	require.NoError(t, err)

	_, err = client.Unshorten(short)
	require.NoError(t, err)

	clock.Advance(2 * time.Hour)

	_, err = client.Unshorten(short)
	assert.ErrorIs(t, err, ErrExpired)

	count, err := app.CountingUsecase.countStore.Get(short)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}
