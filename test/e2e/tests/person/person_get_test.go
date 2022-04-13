package person

import (
	"github.com/GeneralKenobi/census/pkg/api/apimodel"
	"github.com/GeneralKenobi/census/pkg/util"
	"github.com/GeneralKenobi/census/test/e2e"
	testutil2 "github.com/GeneralKenobi/census/test/e2e/e2eutil"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	api := e2e.GetApi()
	ctx := e2e.Context(t)

	personDefinition := apimodel.PersonCreate{
		Name:        "John",
		Surname:     "Smith",
		Email:       testutil2.RandomEmail(),
		DateOfBirth: util.Date(1995, time.September, 13),
		Hobby:       "Jogging",
	}
	personCreated, err := api.CreatePerson(ctx, personDefinition)
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	person, err := api.GetPerson(ctx, personCreated.Id)
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if personDefinition.Name != person.Name ||
		personDefinition.Surname != person.Surname ||
		personDefinition.Email != person.Email ||
		personDefinition.DateOfBirth != person.DateOfBirth ||
		personDefinition.Hobby != person.Hobby {
		t.Errorf("Expected person to be like %#v but got %#v", personDefinition, person)
	}

	const timeSinceCreatedAtThreshold = 2 * time.Minute
	timeSinceCreatedAt := time.Now().Sub(person.CreatedAt)
	if timeSinceCreatedAt > timeSinceCreatedAtThreshold {
		t.Errorf("Created at %v is too far in the past (%v), expected at most %v",
			person.CreatedAt, timeSinceCreatedAt, timeSinceCreatedAtThreshold)
	}

	if person.CreatedAt != person.LastModifiedAt {
		t.Errorf("Expected last modified at (%v) to be equal to created at (%v)", person.LastModifiedAt, person.CreatedAt)
	}
}

func TestGetShouldFailForNotExistingPersonId(t *testing.T) {
	api := e2e.GetApi()
	ctx := e2e.Context(t)

	randomId := util.RandomAlphanumericString(8)
	_, err := api.GetPerson(ctx, randomId)
	if !testutil2.IsNotFound(err) {
		t.Fatalf("Expected a not found error but got %v", err)
	}
}
