package mongo

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/census/internal/db"
	"github.com/GeneralKenobi/census/internal/db/mongo/repository"
	"go.mongodb.org/mongo-driver/mongo"
)

func (mongoCtx *Context) TransactionalRepository(ctx context.Context) (db.Repository, db.Transaction, error) {
	session, err := mongoCtx.client.StartSession()
	if err != nil {
		return nil, nil, fmt.Errorf("error starting session: %w", err)
	}
	err = session.StartTransaction()
	if err != nil {
		return nil, nil, fmt.Errorf("error starting transaction: %w", err)
	}

	queryCtx := sessionQueryContext{
		session:  session,
		database: mongoCtx.database(),
	}
	transactionalRepository := repository.New(&queryCtx)

	sessionCtx := mongo.NewSessionContext(ctx, session)
	transactionCtx := &transactionContext{sessionCtx: sessionCtx}

	return transactionalRepository, transactionCtx, nil
}

// transactionContext implements the db.Transaction interface.
type transactionContext struct {
	sessionCtx mongo.SessionContext
}

var _ db.Transaction = (*transactionContext)(nil) // Interface guard

func (transactionCtx *transactionContext) Commit() error {
	return transactionCtx.sessionCtx.CommitTransaction(transactionCtx.sessionCtx)
}

func (transactionCtx *transactionContext) Rollback() error {
	return transactionCtx.sessionCtx.AbortTransaction(transactionCtx.sessionCtx)
}

// sessionQueryContext implements a stateful repository.QueryContext, i.e. with an associated session.
type sessionQueryContext struct {
	session  mongo.Session
	database *mongo.Database
}

var _ repository.QueryContext = (*sessionQueryContext)(nil) // Interface guard

func (queryCtx *sessionQueryContext) Database() *mongo.Database {
	return queryCtx.database
}

func (queryCtx *sessionQueryContext) QueryContext(ctx context.Context) context.Context {
	// For the query to be associated with this session and its transaction it needs to be given a session context built with this session.
	return mongo.NewSessionContext(ctx, queryCtx.session)
}
