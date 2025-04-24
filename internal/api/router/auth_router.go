// internal/api/router/auth_router.go
package router

// import (
//     "net/http"
//     "time"

//     "github.com/Neimess/shortener/internal/api/handler"
//     "github.com/Neimess/shortener/internal/api/middleware"
//     "github.com/Neimess/shortener/internal/infrastructure/cache"
// )

// func RegisterAuthRoutes(
//     mux *http.ServeMux,
//     h handler.AuthHandler,
//     cacheClient cache.ExpirableCache,
//     spamLimit int,
//     spamWindow time.Duration,
// ) {

// 	chain := NewChain(ChainOptions{
//         Cache:      cacheClient,
//         SpamLimit:  spamLimit,
//         SpamWindow: spamWindow,
//     })
//     mux.Handle("/register", chain(http.HandlerFunc(h.Register)))
//     mux.Handle("/login",    chain(http.HandlerFunc(h.Login)))
// }
