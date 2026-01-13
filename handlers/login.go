package handlers

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/janphilippgutt/casproject/internal/auth"
	"github.com/janphilippgutt/casproject/internal/db"
)

type LoginData struct {
	Email string
	Next  string
	Error string
}

func Login(t *template.Template, sess *scs.SessionManager, pool *pgxpool.Pool, tokenStore *auth.TokenStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			// Read ?next=... from URL (if any)
			next := r.URL.Query().Get("next")

			data := LoginData{
				Email: "",
				Next:  next,
				Error: "",
			}

			if err := t.ExecuteTemplate(w, "login", data); err != nil {
				log.Println("template execute error:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}

		case http.MethodPost:

			// Use request context for DB calls and cancellation

			ctx := r.Context()

			email := r.FormValue("email")
			next := r.FormValue("next") // comes from hidden form field

			if email == "" {
				data := LoginData{
					Email: "",
					Next:  next,
					Error: "Email is required",
				}
				if err := t.ExecuteTemplate(w, "login", data); err != nil {
					log.Println("template execute error:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			// Look up user in DB
			user, err := db.GetUserByEmail(ctx, pool, email)
			if err != nil {
				// If no rows, show friendly error. For any other DB error, log and show generic error.
				// pgx returns pgx.ErrNoRows for not found (use errors.Is check if needed).

				slog.LogAttrs(
					r.Context(),
					slog.LevelWarn,
					"login failed: user not found",
					slog.String("event.category", "auth"),
					slog.String("event.type", "fail"),
					slog.String("user.id", email),
				)

				data := LoginData{
					Email: email,
					Next:  next,
					Error: "No account found for that email.",
				}
				if err := t.ExecuteTemplate(w, "login", data); err != nil {
					log.Println("template execute error:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			token, err := auth.GenerateToken(32)
			if err != nil {

				slog.LogAttrs(
					r.Context(),
					slog.LevelError,
					"login failed: token generation error",
					slog.String("event.category", "auth"),
					slog.String("event.type", "error"),
					slog.String("user.id", user.Email),
					slog.Any("error", err),
				)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			tokenStore.Add(token, user.Email, 15*time.Minute)

			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}

			// dev: print login link to terminal. In production send per email instead.
			log.Println("Magic login link:")
			log.Println("http://localhost:" + port + "/magic-login?token=" + token)

			slog.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"magic login link issued",
				slog.String("event.category", "auth"),
				slog.String("event.type", "login"),
				slog.String("user.id", user.Email),
			)

			// Show confirmation instead of logging user in
			w.Write([]byte("Check your email for the login link (see server log)."))
			return

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Security note: for production, validate next so it only redirects to safe internal paths (avoid open redirect attacks).
