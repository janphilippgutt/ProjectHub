package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type AboutData struct {
	BasePageData
	Header  string
	Message string
}

func About(t *template.Template, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		base := NewBaseData(r.Context(), sess)

		data := AboutData{
			BasePageData: base,
			Message:      "Read more about the project here",
			Header:       "Welcome to the About Page!",
		}

		if err := t.ExecuteTemplate(w, "about", data); err != nil {
			log.Println("template execute error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
