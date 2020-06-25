package session

import (
	"context"
	"satellity/internal/durable"

	"github.com/unrolled/render"
)

type contextValueKey int

const (
	keyLogger      contextValueKey = 1
	keyRender      contextValueKey = 2
	keyDatabase    contextValueKey = 3
	keyRequestBody contextValueKey = 13
)

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

// Database read the database from context
func Database(ctx context.Context) *durable.Database {
	v, _ := ctx.Value(keyDatabase).(*durable.Database)
	return v
}

// WithDatabase put database in context
func WithDatabase(ctx context.Context, database *durable.Database) context.Context {
	return context.WithValue(ctx, keyDatabase, database)
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
