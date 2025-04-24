package middleware

import (
	"net/http"
	"time"

	"github.com/Neimess/shortener/internal/infrastructure/cache"
	"github.com/Neimess/shortener/internal/util"
)

func SpamFilterMiddleware(store cache.FullCache, limit int, ttl time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := util.ClientIP(r)
			if ip == "" {
				http.Error(w, "cannot determine IP", http.StatusBadRequest)
				return
			}

			key := "spam:" + ip
			count, err := store.Incr(r.Context(), key)
			if err != nil {
				http.Error(w, "internal redis error", http.StatusInternalServerError)
				return
			}

			if count == 1 {
				_ = store.Expire(r.Context(), key, ttl)
			}

			if count > int64(limit) {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
