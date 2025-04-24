package router

import (
	"net/http"
	"time"

	handler "github.com/Neimess/shortener/internal/api/handler/auth"
	"github.com/Neimess/shortener/internal/infrastructure/cache"
)

func RegisterAuthRoutes(
	mux *http.ServeMux,
	h handler.AuthHandler,
	cacheClient cache.FullCache,
	spamLimit int,
	spamWindow time.Duration,
) {

	chain := NewChain(ChainOptions{
		Cache:      cacheClient,
		SpamLimit:  spamLimit,
		SpamWindow: spamWindow,
	})
	mux.Handle("/auth/register", chain(http.HandlerFunc(h.Register)))
	mux.Handle("/auth/login", chain(http.HandlerFunc(h.Login)))
	mux.Handle("/auth/refresh", chain(http.HandlerFunc(h.Refresh)))
}
