package db

import (
	"context"
	"fmt"
)

// Context is the top-level interface implemented by db providers with support for both transactional and no-transaction-guarantee
// repositories.
type Context interface {
	Querent
	Transactioner
}

// Querent is a db manager that creates repositories without any transaction guarantees. Queries may be executed without a
// transaction or within a default one, depending on the underlying DB technology.
type Querent interface {
	// Repository creates a Repository that executes queries without a transaction or within a default one, depending on the underlying
	// DB technology.
	Repository(ctx context.Context) (Repository, error)
}

// Transactioner is a db manager that creates transaction-scoped repositories.
type Transactioner interface {
	// TransactionalRepository creates a Repository that runs all queries within the returned Transaction. It's not allowed to use the
	// repository after committing or rolling back the transaction.
	TransactionalRepository(ctx context.Context) (Repository, Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}

// Repository aggregates all queries implemented by db providers.
type Repository interface {
}

var (
	// ErrNoRows is returned from queries that returned/affected 0 rows but at least 1 was expected (e.g. select one found no rows).
	ErrNoRows = fmt.Errorf("no row matched the query")
	// ErrTooManyRows is returned from queries that returned/affected more rows than expected (e.g. select one found 2 rows).
	ErrTooManyRows = fmt.Errorf("more rows than expected matched the query")
)
