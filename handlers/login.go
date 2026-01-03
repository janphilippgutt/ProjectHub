package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type LoginData struct {
	Email string
}

func Login(t *template.Template, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			// Render the login form
			data := LoginData{Email: ""}
			if err := t.ExecuteTemplate(w, "login", data); err != nil {
				log.Println("template execute error:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		case http.MethodPost:
			// Read submitted email
			email := r.FormValue("email")
			if email == "" {
				http.Error(w, "Email is required", http.StatusBadRequest)
				return
			}

			log.Println("Login attempt for:", email)

			// Store in session
			sess.Put(r.Context(), "authenticated", true)
			sess.Put(r.Context(), "email", email)

			// Redirect to home after login (PRG pattern)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
