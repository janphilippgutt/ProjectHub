package middleware

import (
	"log/slog"
	"net/http"
	"strings"
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

		noisyPath :=
			r.URL.Path == "/favicon.ico" ||
				strings.HasPrefix(r.URL.Path, "/uploads/")

		level := slog.LevelInfo

		if rec.status >= 500 {
			level = slog.LevelError
		} else if rec.status >= 400 {
			level = slog.LevelWarn
		} else if noisyPath {
			// Successful noisy requests → don’t log
			return
		}

		duration := time.Since(start)

		slog.LogAttrs(
			r.Context(),
			level,
			"request completed",
			slog.String("http.method", r.Method),
			slog.String("url.path", r.URL.Path),
			slog.Int("http.status_code", rec.status),
			slog.Int64("event.duration_ms", duration.Milliseconds()),
		)

	})
}
