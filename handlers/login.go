package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type LoginData struct {
	Email string
	Next  string
}

func Login(t *template.Template, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			// Read ?next=... from URL (if any)
			next := r.URL.Query().Get("next")

			data := LoginData{
				Email: "",
				Next:  next,
			}

			if err := t.ExecuteTemplate(w, "login", data); err != nil {
				log.Println("template execute error:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}

		case http.MethodPost:
			email := r.FormValue("email")
			next := r.FormValue("next") // comes from hidden form field

			if email == "" {
				http.Error(w, "Email is required", http.StatusBadRequest)
				return
			}

			log.Println("Login attempt for:", email)

			// Store session values
			sess.Put(r.Context(), "authenticated", true)
			sess.Put(r.Context(), "email", email)

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
