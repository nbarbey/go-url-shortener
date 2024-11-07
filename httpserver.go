package url_shortener

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type HTTPServer struct {
	mux *http.ServeMux
}

func NewHTTPServer(s ShortenUnshortener) *HTTPServer {
	mux := http.NewServeMux()
	mux = withShortenerHandler(s)(mux)
	mux = withUnhortenerHandler(s)(mux)
	return &HTTPServer{mux: mux}
}

func (s *HTTPServer) Start() error {
	go func() { _ = http.ListenAndServe("localhost:8080", s.mux) }()
	return nil
}

type muxModifier func(mux *http.ServeMux) *http.ServeMux

func withShortenerHandler(s Shortener) muxModifier {
	return func(mux *http.ServeMux) *http.ServeMux {
		mux.HandleFunc("/shorten", func(writer http.ResponseWriter, request *http.Request) {
			escapedURL := request.URL.Query().Get("url")
			rawURL, err := url.QueryUnescape(escapedURL)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			shortened, err := s.Shorten(rawURL)
			switch {
			case err == nil:
				_, _ = writer.Write([]byte(fmt.Sprintf(`{"shortened": "%s"}`, shortened)))
			case errors.Is(err, ErrNotFound):
				fallthrough
			case errors.Is(err, ErrMissingScheme):
				fallthrough
			case errors.Is(err, ErrMissingHostname):
				fallthrough
			case errors.Is(err, ErrInvalidURL):
				writer.WriteHeader(http.StatusBadRequest)
				_, _ = writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			default:
				writer.WriteHeader(http.StatusInternalServerError)
			}
		})
		return mux
	}
}

func withUnhortenerHandler(u Unshortener) muxModifier {
	return func(mux *http.ServeMux) *http.ServeMux {
		mux.HandleFunc("/unshorten", func(writer http.ResponseWriter, request *http.Request) {
			escapedURL := request.URL.Query().Get("url")
			rawURL, err := url.QueryUnescape(escapedURL)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			shortened, err := u.Unshorten(rawURL)
			switch {
			case errors.Is(err, ErrNotFound):
				fallthrough
			case errors.Is(err, ErrMissingScheme):
				fallthrough
			case errors.Is(err, ErrMissingHostname):
				fallthrough
			case errors.Is(err, ErrInvalidURL):
				writer.WriteHeader(http.StatusBadRequest)
				_, _ = writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			case err == nil:
				_, _ = writer.Write([]byte(fmt.Sprintf(`{"unshortened": "%s"}`, shortened)))
			default:
				writer.WriteHeader(http.StatusInternalServerError)
			}
		})
		return mux
	}
}
