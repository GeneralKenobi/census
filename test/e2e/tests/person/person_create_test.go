package person

import (
	"github.com/GeneralKenobi/census/pkg/api/apimodel"
	"github.com/GeneralKenobi/census/pkg/util"
	"github.com/GeneralKenobi/census/test/e2e"
	"github.com/GeneralKenobi/census/test/e2e/e2eutil"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	api := e2e.GetApi()
	ctx := e2e.Context(t)

	person := apimodel.PersonCreate{
		Name:        "John",
		Surname:     "Smith",
		Email:       e2eutil.RandomEmail(),
		DateOfBirth: util.Date(1995, time.April, 13),
		Hobby:       "Jogging",
	}

	personCreated, err := api.CreatePerson(ctx, person)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}
	if personCreated.Id == "" {
		t.Fatalf("Expected to get a non-empty ID for the created person")
	}
}

func TestCreateShouldFailForDuplicateEmail(t *testing.T) {
	api := e2e.GetApi()
	ctx := e2e.Context(t)

	// Creating the first person should succeed.
	person1 := apimodel.PersonCreate{
		Name:        "John",
		Surname:     "Smith",
		Email:       e2eutil.RandomEmail(),
		DateOfBirth: util.Date(1995, time.August, 13),
		Hobby:       "Jogging",
	}
	_, err := api.CreatePerson(ctx, person1)
	if err != nil {
		t.Fatalf("Expected no error but got: %v", err)
	}

	// Creating the second person with the same email as the first person should fail with bad request.
	person2 := apimodel.PersonCreate{
		Name:        "Jane",
		Surname:     "Doe",
		Email:       person1.Email,
		DateOfBirth: util.Date(1988, time.January, 11),
	}
	_, err = api.CreatePerson(ctx, person2)
	if !e2eutil.IsBadRequest(err) {
		t.Fatalf("Expected a bad request error but got %v", err)
	}
}
