package person

import (
	"github.com/GeneralKenobi/census/pkg/api/apimodel"
	"github.com/GeneralKenobi/census/pkg/util"
	"github.com/GeneralKenobi/census/test/e2e"
	"github.com/GeneralKenobi/census/test/e2e/e2eutil"
	"testing"
	"time"
)

func TestDelete(t *testing.T) {
	api := e2e.GetApi()
	ctx := e2e.Context(t)

	person := apimodel.PersonCreate{
		Name:        "John",
		Surname:     "Smith",
		Email:       e2eutil.RandomEmail(),
		DateOfBirth: util.Date(1995, time.August, 13),
		Hobby:       "Jogging",
	}

	personCreated, err := api.CreatePerson(ctx, person)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	err = api.DeletePerson(ctx, personCreated.Id)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	_, err = api.GetPerson(ctx, personCreated.Id)
	if !e2eutil.IsNotFound(err) {
		t.Fatalf("Expected a not found error but got %v", err)
	}
}

func TestDeleteShouldFailForNotExistingPersonId(t *testing.T) {
	api := e2e.GetApi()
	ctx := e2e.Context(t)

	randomId := util.RandomAlphanumericString(8)
	err := api.DeletePerson(ctx, randomId)
	if !e2eutil.IsNotFound(err) {
		t.Fatalf("Expected a not found error but got %v", err)
	}
}
