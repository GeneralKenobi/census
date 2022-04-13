package updater

import (
	"context"
	"errors"
	"fmt"
	"github.com/GeneralKenobi/census/internal/api"
	"github.com/GeneralKenobi/census/internal/db"
	"github.com/GeneralKenobi/census/internal/db/model"
	"github.com/GeneralKenobi/census/pkg/api/apimodel"
	"github.com/GeneralKenobi/census/pkg/mdctx"
	"time"
)

type Repository interface {
	FindPersonById(ctx context.Context, id string) (model.Person, error)
	UpdatePerson(ctx context.Context, person model.Person) (model.Person, error)
}

func New(repository Repository) *PersonUpdater {
	return &PersonUpdater{repository: repository}
}

type PersonUpdater struct {
	repository Repository
}

func (updater *PersonUpdater) UpdateFromDto(ctx context.Context, id string, personDto apimodel.PersonUpdate) (model.Person, error) {
	mdctx.Debugf(ctx, "Fetching person %s for update", id)
	person, err := updater.repository.FindPersonById(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNoRows) {
			return model.Person{}, api.StatusNotFound.WithMessageAndCause(err, "person with ID %s doesn't exist", id)
		}
		return model.Person{}, fmt.Errorf("error getting person %s for update: %w", id, err)
	}

	mdctx.Debugf(ctx, "Updating person %s", id)
	person.Name = personDto.Name
	person.Surname = personDto.Surname
	person.DateOfBirth = personDto.DateOfBirth
	person.Hobby = personDto.Hobby
	person.LastModifiedAt = timeNow()

	person, err = updater.repository.UpdatePerson(ctx, person)
	if err != nil {
		return model.Person{}, fmt.Errorf("error updating person %s: %w", id, err)
	}

	mdctx.Infof(ctx, "Updated person %s", id)
	return person, nil
}

// Hook for mocking in unit tests.
var timeNow = time.Now
