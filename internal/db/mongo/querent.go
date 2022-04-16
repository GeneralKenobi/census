package mongo

import (
	"context"
	"github.com/GeneralKenobi/census/internal/db"
	"github.com/GeneralKenobi/census/internal/db/mongo/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mongoCtx *Context) Repository(ctx context.Context) (db.Repository, error) {
	queryCtx := standaloneQueryContext{database: mongoCtx.database()}
	return repository.New(&queryCtx), nil
}

// standaloneQueryContext implements a stateless repository.QueryContext, i.e. with no session state.
type standaloneQueryContext struct {
	database *mongo.Database
}

var _ repository.QueryContext = (*standaloneQueryContext)(nil) // Interface guard

func (standaloneCtx *standaloneQueryContext) Database() *mongo.Database {
	return standaloneCtx.database
}

func (standaloneCtx *standaloneQueryContext) QueryContext(ctx context.Context) context.Context {
	// For session-less queries any context can be used - no need to enhance it with session data.
	return ctx
}
