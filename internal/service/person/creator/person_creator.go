package creator

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
	FindPersonByEmail(ctx context.Context, email string) (model.Person, error)
	InsertPerson(ctx context.Context, person model.Person) (model.Person, error)
}

func New(repository Repository) *PersonCreator {
	return &PersonCreator{repository: repository}
}

type PersonCreator struct {
	repository Repository
}

// CreateFromDto creates a person based on the properties in the DTO. If a person with the same email already exists an error is returned.
func (creator *PersonCreator) CreateFromDto(ctx context.Context, personDto apimodel.PersonCreate) (model.Person, error) {
	if err := creator.assertEmailIsNotUsed(ctx, personDto.Email); err != nil {
		return model.Person{}, err
	}

	mdctx.Debugf(ctx, "Creating new person")
	creationTimestamp := timeNow()
	person := model.Person{
		Name:           personDto.Name,
		Surname:        personDto.Surname,
		Email:          personDto.Email,
		DateOfBirth:    personDto.DateOfBirth,
		Hobby:          personDto.Hobby,
		CreatedAt:      creationTimestamp,
		LastModifiedAt: creationTimestamp,
	}

	person, err := creator.repository.InsertPerson(ctx, person)
	if err != nil {
		return model.Person{}, fmt.Errorf("error creating person: %w", err)
	}
	mdctx.Infof(ctx, "Created person %s", person.Id)
	return person, nil
}

func (creator *PersonCreator) assertEmailIsNotUsed(ctx context.Context, email string) error {
	person, err := creator.repository.FindPersonByEmail(ctx, email)
	if err == nil {
		mdctx.Debugf(ctx, "Email is already used by person %s", person.Id)
		return api.StatusBadInput.WithMessage("person with this email already exists")
	}
	if errors.Is(err, db.ErrNoRows) {
		mdctx.Debugf(ctx, "No person found - email is available")
		return nil
	}
	return fmt.Errorf("error checking if person with email already exists: %w", err)
}

// Hook for mocking in unit tests.
var timeNow = time.Now
