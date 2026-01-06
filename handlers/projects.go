package handlers

import (
	"html/template"
	"net/http"
)

// Establish a struct, even if empty to have predictable templates that can later be filled with properties like
// Error string
// Title string
type ProjectNewData struct{}

func NewProject(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "project_new", ProjectNewData{}); err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
		}
	}
}
