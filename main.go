package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/janphilippgutt/casproject/handlers"
	"github.com/janphilippgutt/casproject/internal/auth"
	"github.com/janphilippgutt/casproject/internal/db"
	"github.com/janphilippgutt/casproject/internal/repository"
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

	tokenStore := auth.NewTokenStore()

	tokenStore.StartCleanup(1 * time.Minute)

	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Cookie.Secure = false // for local dev; set true in production

	// parse per-page template sets (base + specific page)
	tpls := map[string]*template.Template{
		"home":           mustParse("home", "templates/base.html", "templates/home.html"),
		"login":          mustParse("login", "templates/base.html", "templates/login.html"),
		"about":          mustParse("about", "templates/base.html", "templates/about.html"),
		"admin":          mustParse("admin", "templates/base.html", "templates/admin.html"),
		"new_project":    mustParse("new_project", "templates/base.html", "templates/project_new.html"),
		"projects":       mustParse("projects", "templates/base.html", "templates/projects.html"),
		"admin_projects": mustParse("admin_projects", "templates/base.html", "templates/admin_projects.html"),
	}

	r := chi.NewRouter()

	// Wrap router with session manager middleware
	r.Use(sessionManager.LoadAndSave)

	// create the middleware
	authMW := middleware.AuthRequired(sessionManager)
	requireAdmin := middleware.RequireAdmin(sessionManager)

	// Create repository once
	projectRepo := &repository.ProjectRepository{DB: dbPool}

	// Create file server
	fileServer := http.FileServer(http.Dir("./uploads"))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fileServer))

	// use it for a route
	r.With(authMW, requireAdmin).Get("/admin", handlers.Admin(tpls["admin"], sessionManager))
	r.With(authMW).Get("/projects/new", handlers.NewProject(tpls["new_project"], projectRepo, sessionManager))
	r.With(authMW).Post("/projects/new", handlers.NewProject(tpls["new_project"], projectRepo, sessionManager))
	r.With(authMW, requireAdmin).Get("/admin/projects", handlers.ListUnapprovedProjects(tpls["admin_projects"], projectRepo))
	r.With(authMW, requireAdmin).Post("/admin/projects/{id}/approve", handlers.ApproveProject(projectRepo))

	// inject the correct template set into each handler
	r.Get("/", handlers.Home(tpls["home"], sessionManager))
	r.Get("/login", handlers.Login(tpls["login"], sessionManager, dbPool, tokenStore))
	r.Post("/login", handlers.Login(tpls["login"], sessionManager, dbPool, tokenStore))
	r.Get("/magic-login", handlers.MagicLogin(sessionManager, dbPool, tokenStore))
	r.Get("/about", handlers.About(tpls["about"]))
	r.Get("/projects", handlers.ListProjects(tpls["projects"], projectRepo))

	log.Println("Server running on :8080") // log -> timestamps included, consistent logging style, logs can easily be redirected later
	log.Fatal(http.ListenAndServe(":8080", r))
}

// later add MIME validation, size limits, image type checks
