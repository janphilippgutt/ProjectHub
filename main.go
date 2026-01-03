package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/janphilippgutt/casproject/handlers"
)

func mustParse(name string, files ...string) *template.Template {
	t := template.Must(template.ParseFiles(files...))
	log.Printf("parsed templates for %s: %q\n", name, t.DefinedTemplates())
	return t
}

func main() {
	// parse per-page template sets (base + specific page)
	tpls := map[string]*template.Template{
		"home":  mustParse("home", "templates/base.html", "templates/home.html"),
		"login": mustParse("login", "templates/base.html", "templates/login.html"),
	}

	r := chi.NewRouter()

	// inject the correct template set into each handler
	r.Get("/", handlers.Home(tpls["home"]))
	r.Get("/login", handlers.Login(tpls["login"]))

	log.Println("Server running on :8080") // log -> timestamps included, consistent logging style, logs can easily be redirected later
	log.Fatal(http.ListenAndServe(":8080", r))
}
