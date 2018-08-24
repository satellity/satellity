package session

import (
	"context"

	"github.com/go-pg/pg"
	"github.com/godiscourse/godiscourse/durable"
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
)

func Database(ctx context.Context) *pg.DB {
	v, _ := ctx.Value(keyDatabase).(*pg.DB)
	return v
}

func WithDatabase(ctx context.Context, database *pg.DB) context.Context {
	return context.WithValue(ctx, keyDatabase, database)
}

func Logger(ctx context.Context) *durable.Logger {
	v, _ := ctx.Value(keyLogger).(*durable.Logger)
	return v
}

func WithLogger(ctx context.Context, logger *durable.Logger) context.Context {
	return context.WithValue(ctx, keyLogger, logger)
}

func Render(ctx context.Context) *render.Render {
	v, _ := ctx.Value(keyRender).(*render.Render)
	return v
}

func WithRender(ctx context.Context, r *render.Render) context.Context {
	return context.WithValue(ctx, keyRender, r)
}
