package models

import (
	"context"

	"github.com/godiscourse/godiscourse/api/durable"
)

type Context struct {
	context  context.Context
	database *durable.Database
}

func WrapContext(ctx context.Context, db *durable.Database) *Context {
	return &Context{context: ctx, database: db}
}
