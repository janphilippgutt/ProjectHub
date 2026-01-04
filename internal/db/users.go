package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/janphilippgutt/casproject/internal/models"
)

func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (*models.User, error) {
	row := pool.QueryRow(
		ctx,
		`SELECT id, email, role, created_at
		 FROM users
		 WHERE email = $1`,
		email,
	)

	var u models.User
	err := row.Scan(&u.ID, &u.Email, &u.Role, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
