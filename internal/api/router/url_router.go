package router

import (
	"net/http"
	"time"

	handler "github.com/Neimess/shortener/internal/api/handler/url"
	"github.com/Neimess/shortener/internal/infrastructure/cache"
)

func RegisterURLRoutes(
	mux *http.ServeMux,
	h handler.URLHandler,
	cacheClient cache.FullCache,
	spamLimit int,
	spamWindow time.Duration,
) {
	chain := NewChain(ChainOptions{
		Cache:      cacheClient,
		SpamLimit:  spamLimit,
		SpamWindow: spamWindow,
	})

	mux.Handle("/shorten", chain(http.HandlerFunc(h.Shorten)))
	mux.Handle("/healthz", chain(http.HandlerFunc(h.Health)))
	mux.Handle("/", chain(http.HandlerFunc(h.Redirect)))

}
