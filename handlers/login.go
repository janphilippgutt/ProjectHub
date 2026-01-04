package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/janphilippgutt/casproject/internal/db"
)

type LoginData struct {
	Email string
	Next  string
	Error string
}

func Login(t *template.Template, sess *scs.SessionManager, pool *pgxpool.Pool) http.HandlerFunc {
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
				log.Println("user lookup error:", err)
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

			// user exists -> set session values from DB
			role := user.Role // comes from DB

			// Store session values
			sess.Put(r.Context(), "authenticated", true)
			sess.Put(r.Context(), "email", email)
			sess.Put(r.Context(), "role", role)

			log.Printf("login: user=%s role=%s\n", user.Email, role)

			// Redirect back to originally requested page if available
			if next != "" {
				http.Redirect(w, r, next, http.StatusSeeOther)
				return
			}

			// Default fallback
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Security note: for production, validate next so it only redirects to safe internal paths (avoid open redirect attacks).
