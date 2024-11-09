package urlshortener

import (
	"errors"
	"fmt"
	"github.com/MadAppGang/httplog"
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
	mux = withURedirectHandler(s)(mux)

	return &HTTPServer{mux: mux}
}

func (s *HTTPServer) Start() error {
	return http.ListenAndServe(":8080", s.mux)
}

type muxModifier func(mux *http.ServeMux) *http.ServeMux

func withShortenerHandler(s Shortener) muxModifier {
	return func(mux *http.ServeMux) *http.ServeMux {
		mux.Handle("/shorten", httplog.Logger(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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
		})))
		return mux
	}
}

func withUnhortenerHandler(u Unshortener) muxModifier {
	return func(mux *http.ServeMux) *http.ServeMux {
		mux.Handle("/unshorten", httplog.Logger(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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
		})))
		return mux
	}
}

func withURedirectHandler(u Unshortener) muxModifier {
	return func(mux *http.ServeMux) *http.ServeMux {
		mux.Handle("/u/{path}", httplog.Logger(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			path := request.PathValue("path")

			rawURL := fmt.Sprintf("https://localhost:8080/u/%s", path)
			unshortened, err := u.Unshorten(rawURL)
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
				writer.Header().Set("Location", unshortened)
				writer.WriteHeader(http.StatusTemporaryRedirect)
			default:
				writer.WriteHeader(http.StatusInternalServerError)
			}
		})))
		return mux
	}
}
