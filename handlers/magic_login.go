package handlers

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/janphilippgutt/casproject/internal/auth"
	"github.com/janphilippgutt/casproject/internal/db"
)

func MagicLogin(
	sess *scs.SessionManager,
	pool *pgxpool.Pool,
	tokenStore *auth.TokenStore,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "Missing token", http.StatusBadRequest)
			return
		}

		email, ok := tokenStore.Use(token)
		if !ok {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		user, err := db.GetUserByEmail(r.Context(), pool, email)
		if err != nil {
			log.Println("magic login db error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		sess.Put(r.Context(), "authenticated", true)
		sess.Put(r.Context(), "user_email", user.Email)
		sess.Put(r.Context(), "role", user.Role)

		http.Redirect(w, r, "/projects/new", http.StatusSeeOther)
	}
}
