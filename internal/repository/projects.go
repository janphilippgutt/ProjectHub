// Repository for projects data concerns

package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/janphilippgutt/casproject/internal/models"
)

type ProjectRepository struct {
	DB *pgxpool.Pool
}

func (r *ProjectRepository) ListApproved(ctx context.Context) ([]models.Project, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, title, project_description, image_path, author_email, approved, created_at
		FROM projects
		WHERE approved = true
		AND deleted_at IS NULL
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var p models.Project
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.ImagePath,
			&p.AuthorEmail,
			&p.Approved,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}

		projects = append(projects, p)
	}

	return projects, rows.Err()
}

func (r *ProjectRepository) ListUnapproved(ctx context.Context) ([]models.Project, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, title, project_description, image_path, author_email, approved, created_at
		FROM projects
		WHERE approved = false
		AND deleted_at IS NULL
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Description,
			&p.ImagePath,
			&p.AuthorEmail,
			&p.Approved,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *ProjectRepository) Create(
	ctx context.Context,
	title string,
	description string,
	imagePath string,
	authorEmail string,
) error {
	_, err := r.DB.Exec(ctx, `
		INSERT INTO projects (title, project_description, image_path, author_email, approved)
		VALUES ($1, $2, $3, $4, false)
	`,
		title,
		description,
		imagePath,
		authorEmail,
	)

	return err
}

func (r *ProjectRepository) Approve(
	ctx context.Context,
	projectID int,
) error {
	_, err := r.DB.Exec(ctx, `
		UPDATE projects
		SET approved = true
		WHERE id = $1
	`, projectID)

	return err
}

func (r *ProjectRepository) SoftDelete(ctx context.Context, projectID int) error {
	query := `
		UPDATE projects
		SET deleted_at = NOW()
		WHERE id = $1
		AND deleted_at IS NULL
	`
	cmd, err := r.DB.Exec(ctx, query, projectID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("project not found or already deleted")
	}

	return nil
}
