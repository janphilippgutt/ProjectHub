package middleware

import (
	"context"
	"log/slog"
	"strings"
	"time"
)

type ContextHandler struct {
	handler slog.Handler
}

func NewContextHandler(h slog.Handler) slog.Handler {
	return &ContextHandler{handler: h}
}

func (h *ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ContextHandler) Handle(ctx context.Context, r slog.Record) error {

	reqID := RequestIDFromContext(ctx)
	if reqID != "" {
		r.AddAttrs(slog.String("request_id", reqID))
	}

	r.AddAttrs(slog.Time("@timestamp", time.Now().UTC()))

	r.AddAttrs(slog.String("log.level", strings.ToLower(r.Level.String())))

	return h.handler.Handle(ctx, r)
}

func (h ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ContextHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h ContextHandler) WithGroup(name string) slog.Handler {
	return &ContextHandler{handler: h.handler.WithGroup(name)}
}
