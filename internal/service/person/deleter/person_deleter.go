package deleter

import (
	"context"
	"errors"
	"fmt"
	"github.com/GeneralKenobi/census/internal/api"
	"github.com/GeneralKenobi/census/internal/db"
	"github.com/GeneralKenobi/census/pkg/mdctx"
)

type Repository interface {
	DeletePersonById(ctx context.Context, id string) error
}

func New(repository Repository) *PersonDeleter {
	return &PersonDeleter{repository: repository}
}

type PersonDeleter struct {
	repository Repository
}

func (deleter *PersonDeleter) DeleteById(ctx context.Context, id string) error {
	mdctx.Infof(ctx, "Deleting person %s", id)
	err := deleter.repository.DeletePersonById(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return api.StatusNotFound.WithMessageAndCause(err, "person with ID %s doesn't exist", id)
		}
		return fmt.Errorf("error getting person %s: %w", id, err)
	}

	mdctx.Debugf(ctx, "Deleted person %s", id)
	return nil
}
