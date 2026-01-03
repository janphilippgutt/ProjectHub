package handlers

import (
	"html/template"
	"log"
	"net/http"
)

type LoginData struct {
	Email string
}

func Login(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := LoginData{Email: ""}

		if err := t.ExecuteTemplate(w, "login", data); err != nil {
			log.Println("template execute error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
