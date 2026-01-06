package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
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

			// Parse multipart form (max 10 MB)
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				http.Error(w, "Could not parse form", http.StatusBadRequest)
				return
			}

			var imagePath string

			file, header, err := r.FormFile("image")
			if err == nil {
				defer file.Close()

				filename := generateImageFilename(header.Filename)
				dstPath := filepath.Join("uploads", "projects", filename)

				dst, err := os.Create(dstPath)
				if err != nil {
					http.Error(w, "Could not save image", http.StatusInternalServerError)
					return
				}
				defer dst.Close()

				if _, err := io.Copy(dst, file); err != nil {
					http.Error(w, "Could not write image", http.StatusInternalServerError)
					return
				}

				imagePath = "/uploads/projects/" + filename
			}

			title := r.FormValue("title")
			description := r.FormValue("description")
			authorEmail := sess.GetString(r.Context(), "email")

			repoErr := repo.Create(r.Context(), title, description, imagePath, authorEmail)
			if repoErr != nil {
				log.Println("create project error:", repoErr)
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

func generateImageFilename(original string) string {
	ext := filepath.Ext(original)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}
