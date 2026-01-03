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

			// For now, just render the form again with the email pre-filled
			data := LoginData{Email: email}
			if err := t.ExecuteTemplate(w, "login", data); err != nil {
				log.Println("template execute error:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
