package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		duration := time.Since(start)
		log.Printf(
			"[%s] %s %d %s (%s)",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}
