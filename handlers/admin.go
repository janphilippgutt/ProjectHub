package handlers

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/janphilippgutt/casproject/internal/models"
	"github.com/janphilippgutt/casproject/internal/repository"
)

type AdminData struct {
	BasePageData
	Email string
}

type AdminProjectsData struct {
	BasePageData
	Pending  []models.Project
	Approved []models.Project
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

		data := AdminData{
			BasePageData: NewBaseData(r.Context(), sess),
			Email:        email,
		}
		if err := t.ExecuteTemplate(w, "admin", data); err != nil {
			log.Println("admin template error:", err)
			// Do not attempt to write another header after partial write;
			// ExecuteTemplate typically writes the whole body, but if it fails
			// before writing anything the following error response is safe.
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func ListUnapprovedProjects(t *template.Template, repo *repository.ProjectRepository, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pending, err := repo.ListUnapproved(r.Context())
		if err != nil {
			log.Println("list unapproved projects error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		approved, err := repo.ListApproved(r.Context())
		if err != nil {
			log.Println("list approved projects error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := AdminProjectsData{
			BasePageData: NewBaseData(r.Context(), sess),
			Pending:      pending,
			Approved:     approved,
		}

		if err := t.ExecuteTemplate(w, "admin_projects", data); err != nil {
			log.Println("template execute error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func ApproveProject(repo *repository.ProjectRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "id")

		projectID, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid project ID", http.StatusBadRequest)
			return
		}

		if err := repo.Approve(r.Context(), projectID); err != nil {
			log.Println("approve project error:", err)
			http.Error(w, "Failed to approve project", http.StatusInternalServerError)
			return
		}

		// PRG pattern
		http.Redirect(w, r, "/admin/projects", http.StatusSeeOther)
	}
}

func AdminArchivedProjects(
	t *template.Template,
	repo *repository.ProjectRepository,
	sess *scs.SessionManager,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projects, err := repo.ListArchived(r.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := struct {
			BasePageData
			Archived []models.Project
		}{
			BasePageData: NewBaseData(r.Context(), sess),
			Archived:     projects,
		}

		t.ExecuteTemplate(w, "admin_archived_projects", data)
	}
}

func RestoreProject(repo *repository.ProjectRepository, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid project ID", http.StatusBadRequest)
			return
		}

		if err := repo.Restore(ctx, id); err != nil {
			slog.Error(
				"failed to restore project",
				"event.category", "admin",
				"event.type", "restore",
				"project.id", id,
				"error", err,
			)
			http.Error(w, "Could not restore project", http.StatusInternalServerError)
			return
		}

		userEmail := sess.GetString(ctx, "user_email")
		slog.Info(
			"project restored",
			"event.category", "admin",
			"event.type", "restore",
			"user.id", userEmail,
			"project.id", id,
		)

		http.Redirect(w, r, "/admin/projects", http.StatusSeeOther)
	}
}

func DeleteProjectForever(
	repo *repository.ProjectRepository,
	sess *scs.SessionManager,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		idStr := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid project ID", http.StatusBadRequest)
			return
		}

		if err := repo.DeleteForever(ctx, id); err != nil {
			slog.Error(
				"failed to permanently delete project",
				"event.category", "admin",
				"event.type", "delete_forever",
				"project.id", id,
				"error", err,
			)
			http.Error(w, "Could not delete project", http.StatusInternalServerError)
			return
		}

		slog.Info(
			"project permanently deleted",
			"event.category", "admin",
			"event.type", "delete_forever",
			"user.id", sess.GetString(ctx, "user_email"),
			"project.id", id,
		)

		http.Redirect(w, r, "/admin/projects/archived", http.StatusSeeOther)
	}
}
