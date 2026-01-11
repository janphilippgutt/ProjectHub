package middleware

import (
	"context"
	"log/slog"
)

type ContextHandler struct {
	handler slog.Handler
}

func NewContextHandler(h slog.Handler) slog.Handler {
	return ContextHandler{handler: h}
}

func (h ContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, attr := range LogAttrsFromContext(ctx) {
		r.AddAttrs(slog.Any(attr.(string), attr))
	}
	return h.handler.Handle(ctx, r)
}

func (h ContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return ContextHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h ContextHandler) WithGroup(name string) slog.Handler {
	return ContextHandler{handler: h.handler.WithGroup(name)}
}
