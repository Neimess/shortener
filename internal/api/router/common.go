package router

import (
	"net/http"
	"time"

	"github.com/Neimess/shortener/internal/api/middleware"
	"github.com/Neimess/shortener/internal/infrastructure/cache"
)

type ChainOptions struct {
	Cache      cache.FullCache
	SpamLimit  int
	SpamWindow time.Duration
}

func NewChain(opts ChainOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := middleware.LoggerMiddleware(next)
		h = middleware.PrometheusMiddleware(h)
		h = middleware.CORSMiddleware(h)
		h = middleware.SpamFilterMiddleware(opts.Cache, opts.SpamLimit, opts.SpamWindow)(h)
		return h
	}
}
