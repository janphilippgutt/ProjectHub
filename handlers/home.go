// handlers package defines how to handle requests

package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type HomeData struct {
	Name string
}

func Home(t *template.Template, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var name string

		if sess.Exists(r.Context(), "authenticated") && sess.GetBool(r.Context(), "authenticated") {
			// logged-in user
			name = sess.GetString(r.Context(), "email")
		} else {
			// public visitor
			name = "Guest"
		}

		data := HomeData{Name: name}

		// Execute the page's entry template (named "home")
		if err := t.ExecuteTemplate(w, "home", data); err != nil {
			// log server-side, but don't attempt to overwrite an already-written response
			log.Println("template execute error:", err)
			// safe: send an error only if nothing was written yet — but keeping it simple:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
