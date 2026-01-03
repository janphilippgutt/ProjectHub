package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type AboutData struct {
	Header  string
	Message string
}

func About(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := AboutData{Message: "Read more about the project here", Header: "Welcome to the About Page!"}

		if err := t.ExecuteTemplate(w, "about", data); err != nil {
			log.Println("template execute error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
