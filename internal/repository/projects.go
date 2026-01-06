// Repository for projects data concerns

package repository

import (
	"context"

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

func (r *ProjectRepository) Create(
	ctx context.Context,
	title string,
	description string,
	authorEmail string,
) error {
	_, err := r.DB.Exec(ctx, `
		INSERT INTO projects (title, project_description, author_email, approved)
		VALUES ($1, $2, $3, false)
	`,
		title,
		description,
		authorEmail,
	)

	return err
}
