package middleware

import (
	"net/http"
	"strconv"
	"time"

	metrics "github.com/Neimess/shortener/internal/maintenance/metric"
)

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		duration := time.Since(start).Seconds()
		metrics.RequestCount.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(rw.statusCode)).Inc()
		metrics.RequestDuration.WithLabelValues(r.URL.Path).Observe(duration)
	})
}
