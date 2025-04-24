package router

import (
	urlH "github.com/Neimess/shortener/internal/api/handler/url"
	"net/http"
)

func (r *Router) URLRoutes(h urlH.URLHandler) {
	chain := NewChain(ChainOptions{Cache: r.Cache, SpamLimit: r.SpamLimit, SpamWindow: r.SpamWindow})
	r.Mux.Handle("/shorten", chain(http.HandlerFunc(h.Shorten)))
	r.Mux.Handle("/healthz", chain(http.HandlerFunc(h.Health)))
	r.Mux.Handle("/", chain(http.HandlerFunc(h.Redirect)))
}
