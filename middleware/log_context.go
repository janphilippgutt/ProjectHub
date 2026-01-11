package middleware

import "context"

type logCtxKeyType struct{}

var logCtxKey = logCtxKeyType{}

func WithLogAttrs(ctx context.Context, attrs ...any) context.Context {
	return context.WithValue(ctx, logCtxKey, attrs)
}

func LogAttrsFromContext(ctx context.Context) []any {
	attrs, _ := ctx.Value(logCtxKey).([]any)
	return attrs
}
