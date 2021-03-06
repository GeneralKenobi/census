package db

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/census/pkg/mdctx"
	"github.com/GeneralKenobi/census/pkg/util"
)

// InTransaction opens a transaction and runs the given function with transaction-scoped repository.
// If the function returns no errors then transaction is committed.
// If the function returns an error the transaction is rolled back and the error is returned.
// If the function panics the transaction is rolled back and this function re-panics with the original value.
// Error is also returned if it wasn't possible to open or commit a transaction.
func InTransaction(ctx context.Context, transactioner Transactioner, todo func(transactionalRepository Repository) error) error {
	_, err := InTransactionRetV(ctx, transactioner, func(transactionalRepository Repository) (any, error) {
		return nil, todo(transactionalRepository)
	})
	return err
}

// InTransactionRetV opens a transaction and runs the given function with transaction-scoped repository.
// If the function returns no errors then transaction is committed and the function's result is returned.
// If the function returns an error the transaction is rolled back and the error is returned.
// If the function panics the transaction is rolled back and this function re-panics with the original value.
// Error is also returned if it wasn't possible to open or commit a transaction.
func InTransactionRetV[V any](ctx context.Context, transactioner Transactioner, todo func(transactionalRepository Repository) (V, error)) (V, error) {

	mdctx.Debugf(ctx, "Opening transaction")
	transactionalRepository, transaction, err := transactioner.TransactionalRepository(ctx)
	if err != nil {
		return util.ZeroValue[V](), fmt.Errorf("error creating a transaction: %w", err)
	}

	value, err := todo(transactionalRepository)

	if panicErr := recover(); panicErr != nil {
		mdctx.Debugf(ctx, "Rolling back transaction after panic: %v", panicErr)
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			mdctx.Errorf(ctx, "Error rolling back transaction after panic: %v", rollbackErr)
		}
		panic(panicErr)
	}

	if err != nil {
		mdctx.Debugf(ctx, "Rolling back transaction after error: %v", err)
		rollbackErr := transaction.Rollback()
		if rollbackErr != nil {
			mdctx.Errorf(ctx, "Error rolling back transaction after error: %v", rollbackErr)
		}
		return util.ZeroValue[V](), err
	}

	mdctx.Debugf(ctx, "Committing transaction")
	err = transaction.Commit()
	if err != nil {
		return util.ZeroValue[V](), fmt.Errorf("error committing transaction: %w", err)
	}

	return value, nil
}

// WithRepository gets a repository from the querent and runs the given function with it. It returns an error without running the function
// if getting a repository failed.
//
// It's provided to allow for coding style analogous to InTransaction when transaction is not required.
func WithRepository(ctx context.Context, querent Querent, todo func(repository Repository) error) error {
	_, err := WithRepositoryRetV(ctx, querent, func(repository Repository) (any, error) {
		return nil, todo(repository)
	})
	return err
}

// WithRepositoryRetV gets a repository from the querent and runs the given function with it. It returns an error without running the
// function if getting a repository failed.
//
// It's provided to allow for coding style analogous to InTransactionRetV when transaction is not required.
func WithRepositoryRetV[V any](ctx context.Context, querent Querent, todo func(repository Repository) (V, error)) (V, error) {
	repository, err := querent.Repository(ctx)
	if err != nil {
		return util.ZeroValue[V](), fmt.Errorf("error creating a repository: %w", err)
	}

	value, err := todo(repository)
	if err != nil {
		return util.ZeroValue[V](), err
	}

	return value, nil
}
