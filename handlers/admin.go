package handlers

import (
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/alexedwards/scs/v2"
	"github.com/janphilippgutt/casproject/internal/models"
	"github.com/janphilippgutt/casproject/internal/repository"
)

type AdminData struct {
	Email string
}

type AdminProjectsData struct {
	Projects []models.Project
}

func Admin(t *template.Template, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Defensive auth check: if not authenticated, redirect to login with next.
		// This mirrors AuthRequired middleware behavior and ensures safety
		// if the handler is accidentally wired without middleware.
		if !sess.Exists(r.Context(), "authenticated") || !sess.GetBool(r.Context(), "authenticated") {
			next := url.QueryEscape(r.RequestURI)
			http.Redirect(w, r, "/login?next="+next, http.StatusSeeOther)
			return
		}

		// Authorization check: only allow users with role == "admin"
		role := sess.GetString(r.Context(), "role")
		if role != "admin" {
			// Authenticated but not authorized --> 403 Forbidden
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// At this point, the user is authenticated and authorized.
		email := sess.GetString(r.Context(), "email")

		data := AdminData{Email: email}
		if err := t.ExecuteTemplate(w, "admin", data); err != nil {
			log.Println("admin template error:", err)
			// Do not attempt to write another header after partial write;
			// ExecuteTemplate typically writes the whole body, but if it fails
			// before writing anything the following error response is safe.
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func ListUnapprovedProjects(t *template.Template, repo *repository.ProjectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := repo.ListUnapproved(r.Context())
		if err != nil {
			log.Println("list unapproved projects error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := AdminProjectsData{Projects: projects}
		if err := t.ExecuteTemplate(w, "admin_projects", data); err != nil {
			log.Println("template execute error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
