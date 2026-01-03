package middleware

import (
	"net/http"
	"net/url"

	"github.com/alexedwards/scs/v2"
)

// AuthRequired returns a chi-compatible middleware that redirects to /login?next=...
func AuthRequired(sess *scs.SessionManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// if not authenticated, redirect to login with next param
			if !sess.Exists(r.Context(), "authenticated") || !sess.GetBool(r.Context(), "authenticated") {
				nextURL := r.RequestURI // includes path + query
				loginURL := "/login?next=" + url.QueryEscape(nextURL)
				http.Redirect(w, r, loginURL, http.StatusSeeOther)
				return
			}
			// authenticated — continue
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAdmin(sess *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := sess.GetString(r.Context(), "role")
			if role != "admin" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
