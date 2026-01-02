package handlers

import (
	"html/template"
	"net/http"
)

type PageData struct {
	Name string
}

func Home(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create data per Request
		data := PageData{
			Name: "Zag",
		}

		// Execute template per Request
		err := t.ExecuteTemplate(w, "base.html", data) // base.html being the entry point, pulling in blocks from child templates
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
		}
	}

}
