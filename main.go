package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/janphilippgutt/casproject/handlers"
	"github.com/janphilippgutt/casproject/internal/db"
	"github.com/janphilippgutt/casproject/middleware"
)

func mustParse(name string, files ...string) *template.Template {
	t := template.Must(template.ParseFiles(files...))
	log.Printf("parsed templates for %s: %q\n", name, t.DefinedTemplates())
	return t
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found")
	}

	dbPool, err := db.Connect()
	if err != nil {
		log.Fatal("database connection failed:", err)
	}
	defer dbPool.Close()

	log.Println("database connected")

	// TEMPORARY TEST
	user, err := db.GetUserByEmail(context.Background(), dbPool, "admin@example.com")
	if err != nil {
		log.Fatal("query failed:", err)
	}

	log.Printf("loaded user: id=%d email=%s role=%s\n",
		user.ID, user.Email, user.Role)

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
		"admin": mustParse("admin", "templates/base.html", "templates/admin.html"),
	}

	r := chi.NewRouter()

	// Wrap router with session manager middleware
	r.Use(sessionManager.LoadAndSave)

	// create the middleware
	authMW := middleware.AuthRequired(sessionManager)
	requireAdmin := middleware.RequireAdmin(sessionManager)

	// use it for a route (we will add /admin next)
	r.With(authMW, requireAdmin).Get("/admin", handlers.Admin(tpls["admin"], sessionManager))

	// inject the correct template set into each handler
	r.Get("/", handlers.Home(tpls["home"], sessionManager))
	r.Get("/login", handlers.Login(tpls["login"], sessionManager))
	r.Post("/login", handlers.Login(tpls["login"], sessionManager))
	r.Get("/about", handlers.About(tpls["about"]))

	log.Println("Server running on :8080") // log -> timestamps included, consistent logging style, logs can easily be redirected later
	log.Fatal(http.ListenAndServe(":8080", r))
}
