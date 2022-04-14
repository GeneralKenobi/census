package person

import (
	"github.com/GeneralKenobi/census/pkg/api/apimodel"
	"github.com/GeneralKenobi/census/pkg/util"
	"github.com/GeneralKenobi/census/test/e2e"
	"github.com/GeneralKenobi/census/test/e2e/e2eutil"
	"testing"
	"time"
)

func TestUpdate(t *testing.T) {
	api := e2e.GetApi()
	ctx := e2e.Context(t)

	personDefinition := apimodel.PersonCreate{
		Name:        "John",
		Surname:     "Smith",
		Email:       e2eutil.RandomEmail(),
		DateOfBirth: util.Date(1995, time.August, 13),
		Hobby:       "Jogging",
	}
	personCreated, err := api.CreatePerson(ctx, personDefinition)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	personBeforeUpdate, err := api.GetPerson(ctx, personCreated.Id)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	personUpdate := apimodel.PersonUpdate{
		Name:        "James",
		Surname:     "Doe",
		DateOfBirth: time.Date(1965, 6, 11, 0, 0, 0, 0, time.UTC),
		Hobby:       "Reading",
	}
	err = api.UpdatePerson(ctx, personCreated.Id, personUpdate)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	personAfterUpdate, err := api.GetPerson(ctx, personCreated.Id)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}
	if personUpdate.Name != personAfterUpdate.Name ||
		personUpdate.Surname != personAfterUpdate.Surname ||
		personUpdate.DateOfBirth != personAfterUpdate.DateOfBirth ||
		personUpdate.Hobby != personAfterUpdate.Hobby {
		t.Errorf("Expected person to be like %#v but got %#v", personUpdate, personAfterUpdate)
	}

	if personBeforeUpdate.Email != personAfterUpdate.Email {
		t.Errorf("Expected email to be unchanged, but it changed from %s to %s", personBeforeUpdate.Email, personAfterUpdate.Email)
	}
	if personBeforeUpdate.CreatedAt != personAfterUpdate.CreatedAt {
		t.Errorf("Expected created at to be unchanged, but it changed from %s to %s",
			personBeforeUpdate.CreatedAt, personAfterUpdate.CreatedAt)
	}

	if !personAfterUpdate.LastModifiedAt.After(personBeforeUpdate.LastModifiedAt) {
		t.Errorf("Expected last modified at after update (%v) to be newer than before update (%v)",
			personAfterUpdate.LastModifiedAt, personBeforeUpdate.LastModifiedAt)
	}

	const timeSinceLastModifiedAtThreshold = 2 * time.Minute
	timeSinceLastModifiedAt := time.Now().Sub(personBeforeUpdate.LastModifiedAt)
	if timeSinceLastModifiedAt > timeSinceLastModifiedAtThreshold {
		t.Errorf("Last modified at %v is too far in the past (%v), expected at most %v",
			personBeforeUpdate.LastModifiedAt, timeSinceLastModifiedAt, timeSinceLastModifiedAtThreshold)
	}
}

func TestUpdateShouldFailForNotExistingPersonId(t *testing.T) {
	api := e2e.GetApi()
	ctx := e2e.Context(t)

	randomId := util.RandomAlphanumericString(8)
	personUpdate := apimodel.PersonUpdate{
		Name:        "James",
		Surname:     "Doe",
		DateOfBirth: time.Date(1965, 6, 11, 0, 0, 0, 0, time.UTC),
		Hobby:       "Reading",
	}
	err := api.UpdatePerson(ctx, randomId, personUpdate)
	if !e2eutil.IsNotFound(err) {
		t.Fatalf("Expected a not found error but got %v", err)
	}
}
