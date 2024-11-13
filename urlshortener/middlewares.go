package urlshortener

import (
	"net/http"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type middleware interface {
	Handler(h http.Handler) http.Handler
}

type middlewareFunc func(h http.Handler) http.Handler

func (m middlewareFunc) Handler(h http.Handler) http.Handler { return m(h) }

type middlewares []middleware

func (mws middlewares) Handler(h http.Handler) http.Handler {
	for _, mw := range mws {
		h = mw.Handler(h)
	}
	return h
}

func newRateLimiterMiddleware() *stdlib.Middleware {
	return stdlib.NewMiddleware(limiter.New(memory.NewStore(), limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  1000,
	}))
}
