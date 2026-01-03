package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type AdminData struct {
	Email string
}

func Admin(t *template.Template, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read email from session
		email := ""
		if sess.Exists(r.Context(), "authenticated") && sess.GetBool(r.Context(), "authenticated") {
			email = sess.GetString(r.Context(), "email")
		}

		data := AdminData{Email: email}
		if err := t.ExecuteTemplate(w, "admin", data); err != nil {
			log.Println("admin template error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
