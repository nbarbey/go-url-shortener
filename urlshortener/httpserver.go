package urlshortener

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/MadAppGang/httplog"
)

type HTTPServer struct {
	mux *http.ServeMux
}

func NewHTTPServer(s ShortenUnshortener, c CountStorer) *HTTPServer {
	mux := http.NewServeMux()
	mws := []middleware{newRateLimiterMiddleware(), middlewareFunc(httplog.Logger)}
	mux = withShortenerHandler(s, mws...)(mux)
	mux = withUnhortenerHandler(s, mws...)(mux)
	mux = withCount(c)(mux)
	mux = withURedirectHandler(s, mws...)(mux)
	return &HTTPServer{mux: mux}
}

func withCount(c CountStorer, mws ...middleware) func(mux *http.ServeMux) *http.ServeMux {
	return func(mux *http.ServeMux) *http.ServeMux {
		handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			escapedURL := request.URL.Query().Get("url")
			rawURL, err := url.QueryUnescape(escapedURL)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			count, _ := c.Get(rawURL)
			writer.Write([]byte(fmt.Sprintf(`{"count": %d}`, count)))
		})
		mux.Handle("/count", middlewares(mws).Handler(handler))
		return mux
	}
}

func (s *HTTPServer) Start() error {
	return http.ListenAndServe(":8080", s.mux)
}

type muxModifier func(mux *http.ServeMux) *http.ServeMux

func withShortenerHandler(s Shortener, mws ...middleware) muxModifier {
	return func(mux *http.ServeMux) *http.ServeMux {
		var handler http.Handler
		handler = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			escapedURL := request.URL.Query().Get("url")
			rawURL, err := url.QueryUnescape(escapedURL)
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			var expiration *time.Time
			escapedExpiration := request.URL.Query().Get("expiration")
			if escapedExpiration != "" {
				expirationString, err := url.QueryUnescape(escapedExpiration)
				if err != nil {
					writer.WriteHeader(http.StatusBadRequest)
					return
				}
				location, err := time.LoadLocation("Local")
				if err != nil {
					writer.WriteHeader(http.StatusBadRequest)
					return
				}
				e, err := time.ParseInLocation("2006-01-02_15:04:05", expirationString, location)
				if err != nil {
					writer.WriteHeader(http.StatusBadRequest)
					return
				}

				expiration = &e
			}
			shortened, err := s.Shorten(rawURL, expiration)
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

		mux.Handle("/shorten", middlewares(mws).Handler(handler))
		return mux
	}
}

func withUnhortenerHandler(u Unshortener, mws ...middleware) muxModifier {
	return func(mux *http.ServeMux) *http.ServeMux {
		handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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
				fallthrough
			case errors.Is(err, ErrExpired):
				writer.WriteHeader(http.StatusBadRequest)
				_, _ = writer.Write([]byte(fmt.Sprintf(`{"error": "%s"}`, err.Error())))
			case err == nil:
				_, _ = writer.Write([]byte(fmt.Sprintf(`{"unshortened": "%s"}`, shortened)))
			default:
				writer.WriteHeader(http.StatusInternalServerError)
			}
		})

		mux.Handle("/unshorten", middlewares(mws).Handler(handler))
		return mux
	}
}

func withURedirectHandler(u Unshortener, mws ...middleware) muxModifier {
	return func(mux *http.ServeMux) *http.ServeMux {
		handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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
		})
		mux.Handle("/u/{path}", middlewares(mws).Handler(handler))
		return mux
	}
}
