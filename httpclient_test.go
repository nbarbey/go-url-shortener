package url_shortener

import (
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPShortenUnshortener(t *testing.T) {
	ShortenUnshortenerFromBuilder(t, func() ShortenUnshortener {
		app := NewApplication()
		testServer := httptest.NewServer(app.server.mux)
		return NewHTTPClientFromResty(resty.NewWithClient(testServer.Client()).
			SetBaseURL(testServer.URL))
	})

	t.Run("URL not found", func(t *testing.T) {
		app := NewApplication()
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
		app := NewApplication()
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
