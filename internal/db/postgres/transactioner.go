package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GeneralKenobi/census/internal/db"
	"github.com/GeneralKenobi/census/internal/db/postgres/repository"
)

func (postgresCtx *Context) TransactionalRepository(ctx context.Context) (db.Repository, db.Transaction, error) {
	transaction, err := postgresCtx.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	transactionalRepository := repository.New(transaction)
	transactionCtx := &transactionContext{transaction: transaction}
	return transactionalRepository, transactionCtx, nil
}

// transactionContext implements the db.Transaction interface.
type transactionContext struct {
	transaction *sql.Tx
}

var _ db.Transaction = (*transactionContext)(nil) // Interface guard

func (transactionCtx *transactionContext) Commit() error {
	return transactionCtx.transaction.Commit()
}

func (transactionCtx *transactionContext) Rollback() error {
	return transactionCtx.transaction.Rollback()
}
