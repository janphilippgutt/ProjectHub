package main

import (
	"html/template"
	"log"
	"net/http"
)

type PageData struct {
	Name string
}

var homeTemplate *template.Template

func homeHandler(w http.ResponseWriter, r *http.Request) {

	// Create data per Request
	data := PageData{
		Name: "Zag",
	}

	// Execute template per Request
	err := homeTemplate.ExecuteTemplate(w, "base.html", data) // base.html being the entry point, pulling in blocks from child templates
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func main() {

	// Parse template once at startup
	var err error

	homeTemplate, err = template.ParseFiles(
		"templates/base.html",
		"templates/home.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", homeHandler)

	log.Println("Server running on :8080") // log -> timestamps included, consistent logging style, logs can easily be redirected later
	serverErr := http.ListenAndServe(":8080", nil)
	if serverErr != nil {
		log.Fatal(serverErr)
	}
}
