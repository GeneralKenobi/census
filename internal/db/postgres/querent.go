package postgres

import (
	"context"
	"github.com/GeneralKenobi/census/internal/db"
	"github.com/GeneralKenobi/census/internal/db/postgres/repository"
)

func (postgresCtx *Context) Repository(ctx context.Context) (db.Repository, error) {
	return repository.New(postgresCtx.db), nil
}
