package getter

import (
	"context"
	"errors"
	"fmt"
	"github.com/GeneralKenobi/census/internal/api"
	"github.com/GeneralKenobi/census/internal/db"
	"github.com/GeneralKenobi/census/internal/db/model"
	"github.com/GeneralKenobi/census/pkg/mdctx"
)

type Repository interface {
	FindPersonById(ctx context.Context, id string) (model.Person, error)
}

func New(repository Repository) *PersonGetter {
	return &PersonGetter{repository: repository}
}

type PersonGetter struct {
	repository Repository
}

func (getter *PersonGetter) FindById(ctx context.Context, id string) (model.Person, error) {
	mdctx.Debugf(ctx, "Fetching person %s", id)
	person, err := getter.repository.FindPersonById(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return model.Person{}, api.StatusNotFound.WithMessageAndCause(err, "person with ID %s doesn't exist", id)
		}
		return model.Person{}, fmt.Errorf("error getting person %s: %w", id, err)
	}

	return person, nil
}
