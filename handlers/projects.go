package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/janphilippgutt/casproject/internal/models"
	"github.com/janphilippgutt/casproject/internal/repository"
)

// Establish a struct, even if empty to have predictable templates that can later be filled with properties like
// Error string
// Title string
type ProjectNewData struct{}

type ProjectsPageData struct {
	Projects []models.Project
}

func NewProject(t *template.Template, repo *repository.ProjectRepository, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			t.ExecuteTemplate(w, "project_new", ProjectNewData{})
			return

		case http.MethodPost:
			title := r.FormValue("title")
			description := r.FormValue("description")
			authorEmail := sess.GetString(r.Context(), "email")

			err := repo.Create(r.Context(), title, description, authorEmail)
			if err != nil {
				log.Println("create project error:", err)
				http.Error(w, "Failed to create project", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}
}

func ListProjects(t *template.Template, repo *repository.ProjectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := repo.ListApproved(r.Context())
		if err != nil {
			log.Println("list projects error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := ProjectsPageData{Projects: projects}
		if err := t.ExecuteTemplate(w, "projects", data); err != nil {
			log.Println("template execute error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
