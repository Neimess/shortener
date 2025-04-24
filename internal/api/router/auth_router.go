package router

import (
	"net/http"

	authH "github.com/Neimess/shortener/internal/api/handler/auth"
)

func (r *Router) AuthRoutes(h authH.AuthHandler) {
	chain := NewChain(ChainOptions{Cache: r.Cache, SpamLimit: r.SpamLimit, SpamWindow: r.SpamWindow})

	r.Mux.Handle("/auth/register", chain(http.HandlerFunc(h.Register)))
	r.Mux.Handle("/auth/login", chain(http.HandlerFunc(h.Login)))
	r.Mux.Handle("/auth/token/refresh", chain(http.HandlerFunc(h.Refresh)))
}
