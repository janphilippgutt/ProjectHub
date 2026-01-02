package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/janphilippgutt/casproject/handlers"
)

func main() {

	r := chi.NewRouter()

	// Parse template once at startup

	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/home.html"))
	/*if err != nil {
		log.Fatal(err)
	}*/

	r.Get("/", handlers.Home(tmpl))

	// http.HandleFunc("/", homeHandler)

	log.Println("Server running on :8080") // log -> timestamps included, consistent logging style, logs can easily be redirected later
	serverErr := http.ListenAndServe(":8080", r)
	if serverErr != nil {
		log.Fatal(serverErr)
	}
}
