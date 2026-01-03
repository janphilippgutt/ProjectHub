package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/alexedwards/scs/v2"

	"github.com/janphilippgutt/casproject/handlers"
)

func mustParse(name string, files ...string) *template.Template {
	t := template.Must(template.ParseFiles(files...))
	log.Printf("parsed templates for %s: %q\n", name, t.DefinedTemplates())
	return t
}

func main() {

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false // for local dev; set true in production

	// parse per-page template sets (base + specific page)
	tpls := map[string]*template.Template{
		"home":  mustParse("home", "templates/base.html", "templates/home.html"),
		"login": mustParse("login", "templates/base.html", "templates/login.html"),
		"about": mustParse("about", "templates/base.html", "templates/about.html"),
	}

	r := chi.NewRouter()

	// Wrap router with session manager middleware
	r.Use(sessionManager.LoadAndSave)

	// inject the correct template set into each handler
	r.Get("/", handlers.Home(tpls["home"], sessionManager))
	r.Get("/login", handlers.Login(tpls["login"], sessionManager))
	r.Post("/login", handlers.Login(tpls["login"], sessionManager))
	r.Get("/about", handlers.About(tpls["about"]))

	log.Println("Server running on :8080") // log -> timestamps included, consistent logging style, logs can easily be redirected later
	log.Fatal(http.ListenAndServe(":8080", r))
}
