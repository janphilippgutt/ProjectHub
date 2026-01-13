package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func RequestLogger(sess *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
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

			outcome := "success"

			if rec.status >= 400 {
				outcome = "failure"
			}

			duration := time.Since(start)

			userID := sess.GetString(r.Context(), "email") // optional

			slog.LogAttrs(
				r.Context(),
				level,
				"request completed",
				slog.String("event.category", "http"),
				slog.String("event.type", "request"),
				slog.String("http.method", r.Method),
				slog.String("url.path", r.URL.Path),
				slog.Int("http.status_code", rec.status),
				slog.String("event.outcome", outcome),
				slog.Int64("event.duration_ms", duration.Milliseconds()),
				slog.String("user.id", userID), // optional
			)
		})
	}
}
