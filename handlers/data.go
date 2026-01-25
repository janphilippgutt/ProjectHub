package handlers

import (
	"context"
	"log/slog"

	"github.com/alexedwards/scs/v2"
)

// shared view-model that carries user-related UI state
type BasePageData struct {
	IsAuthenticated bool
	IsAdmin         bool
	UserEmail       string
}

// derive user state from the request context via the session manager.
func NewBaseData(ctx context.Context, sess *scs.SessionManager) BasePageData {
	email := sess.GetString(ctx, "user_email")
	role := sess.GetString(ctx, "role")

	slog.Info(
		"base data",
		"user_email", email,
		"role", role,
	)

	return BasePageData{
		IsAuthenticated: email != "",
		IsAdmin:         role == "admin",
		UserEmail:       email,
	}
}
