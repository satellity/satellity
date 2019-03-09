package session

import (
	"context"

	"github.com/godiscourse/godiscourse/api/durable"
	"github.com/unrolled/render"
)

type contextValueKey int

const (
	keyRequest           contextValueKey = 0
	keyDatabase          contextValueKey = 1
	keyLogger            contextValueKey = 2
	keyRender            contextValueKey = 3
	keyRemoteAddress     contextValueKey = 11
	keyAuthorizationInfo contextValueKey = 12
	keyRequestBody       contextValueKey = 13
)

// Database read database from context
func Database(ctx context.Context) *durable.Database {
	v, _ := ctx.Value(keyDatabase).(*durable.Database)
	return v
}

// WithDatabase put database into context
func WithDatabase(ctx context.Context, database *durable.Database) context.Context {
	return context.WithValue(ctx, keyDatabase, database)
}

// Logger read logger from context
func Logger(ctx context.Context) *durable.Logger {
	v, _ := ctx.Value(keyLogger).(*durable.Logger)
	return v
}

// WithLogger put logger into context
func WithLogger(ctx context.Context, logger *durable.Logger) context.Context {
	return context.WithValue(ctx, keyLogger, logger)
}

// Render read render from context
func Render(ctx context.Context) *render.Render {
	v, _ := ctx.Value(keyRender).(*render.Render)
	return v
}

// WithRender put render to context
func WithRender(ctx context.Context, r *render.Render) context.Context {
	return context.WithValue(ctx, keyRender, r)
}

// RequestBody read request body from context
func RequestBody(ctx context.Context) string {
	v, _ := ctx.Value(keyRequestBody).(string)
	return v
}

// WithRequestBody put request body to context
func WithRequestBody(ctx context.Context, body string) context.Context {
	return context.WithValue(ctx, keyRequestBody, body)
}
