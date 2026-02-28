package middleware

import (
	"log"
	"net/http"
	"time"
)

type MetricsMiddleware struct{}

func NewMetricsMiddleware() *MetricsMiddleware {
	return &MetricsMiddleware{}
}

func (m *MetricsMiddleware) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Custom response writer to capture status code
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		log.Printf("[Metrics] %s %s | Status: %d | Duration: %v", r.Method, r.URL.Path, rw.status, duration)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
