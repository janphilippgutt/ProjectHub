package handlers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/janphilippgutt/casproject/internal/models"
	"github.com/janphilippgutt/casproject/internal/repository"
)

// Establish a struct, even if empty to have predictable templates that can later be filled with properties like
// Error string
// Title string
type ProjectNewData struct {
	Title       string
	Description string
	Error       string
}

type ProjectsPageData struct {
	Projects []models.Project
}

func NewProject(t *template.Template, repo *repository.ProjectRepository, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodGet:
			data := ProjectNewData{
				// PopString reads and removes the message; empty string = no message = safe
				Error: sess.PopString(r.Context(), "flash_error"),
			}
			if err := t.ExecuteTemplate(w, "base", data); err != nil {
				log.Println("template execute error:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return

		case http.MethodPost:

			r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				http.Error(w, "Could not parse form", http.StatusBadRequest)
				return
			}

			title := strings.TrimSpace(r.FormValue("title"))
			description := strings.TrimSpace(r.FormValue("description"))
			authorEmail := sess.GetString(r.Context(), "email")

			if title == "" || description == "" {
				sess.Put(r.Context(), "flash_error", "Title and description are required")
				http.Redirect(w, r, "/projects/new", http.StatusSeeOther)
				return
			}

			var imagePath string
			file, header, err := r.FormFile("image")
			if err != nil {
				if err == http.ErrMissingFile {
					imagePath = "" // optional file not provided
				} else {
					log.Println("upload error:", err)
					http.Error(w, "Invalid file upload", http.StatusBadRequest)
					return
				}
			} else {
				defer file.Close()
				if header.Filename != "" {
					buf := make([]byte, 512)
					if _, err := file.Read(buf); err != nil {
						http.Error(w, "Invalid file", http.StatusBadRequest)
						return
					}

					contentType := http.DetectContentType(buf)
					if contentType != "image/jpeg" && contentType != "image/png" {
						http.Error(w, "Only JPEG and PNG allowed", http.StatusBadRequest)
						return
					}

					if _, err := file.Seek(0, 0); err != nil {
						http.Error(w, "Invalid file stream", http.StatusBadRequest)
						return
					}

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
			}

			// Insert project
			if err := repo.Create(r.Context(), title, description, imagePath, authorEmail); err != nil {
				log.Println("create project error:", err)
				http.Error(w, "Failed to create project", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/", http.StatusSeeOther)

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
