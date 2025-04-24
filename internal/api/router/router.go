package router

import (
	"net/http"
	"time"

	"github.com/Neimess/shortener/internal/infrastructure/cache"
	"github.com/Neimess/shortener/internal/util/jwt"
)

type Router struct {
	Mux        *http.ServeMux
	Cache      cache.FullCache
	JWT        jwtutil.JWTManager
	SpamLimit  int
	SpamWindow time.Duration
}

func New(mux *http.ServeMux, cache cache.FullCache, jwtMgr jwtutil.JWTManager, spamLimit int, spamWindow time.Duration) *Router {
	return &Router{Mux: mux, Cache: cache, JWT: jwtMgr, SpamLimit: spamLimit, SpamWindow: spamWindow}
}
