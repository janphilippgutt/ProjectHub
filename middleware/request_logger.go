package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK, // default
		}

		// Call the next handler
		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		reqID := RequestIDFromContext(r.Context())

		// Decide log level
		level := slog.LevelInfo
		if rec.status >= 500 {
			level = slog.LevelError
		} else if rec.status >= 400 && rec.status != http.StatusNotFound {
			level = slog.LevelWarn
		}

		slog.LogAttrs(
			r.Context(),
			level,
			"request completed",
			slog.String("request_id", reqID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.status),
			slog.Int64("duration_ms", duration.Milliseconds()),
		)

	})
}
