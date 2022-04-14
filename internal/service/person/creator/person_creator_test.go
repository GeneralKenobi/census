package creator

import (
	"context"
	"github.com/GeneralKenobi/census/internal/db"
	"github.com/GeneralKenobi/census/internal/db/model"
	"github.com/GeneralKenobi/census/pkg/api/apimodel"
	"github.com/GeneralKenobi/census/pkg/mdctx"
	"github.com/GeneralKenobi/census/pkg/util"
	"github.com/GeneralKenobi/census/test/testutil"
	"testing"
	"time"
)

func TestCreatePersonFromDto(t *testing.T) {
	originalTimeNowHook := timeNow
	defer func() { timeNow = originalTimeNowHook }()
	var mockedCreationTime = time.Now()
	timeNow = func() time.Time { return mockedCreationTime }

	personDto := apimodel.PersonCreate{
		Name:        "John",
		Surname:     "Smith",
		Email:       "john.smith@test.com",
		DateOfBirth: util.Date(2001, time.March, 23),
		Hobby:       "Literature",
	}

	expectedPerson := model.Person{
		Id:             "1ab7",
		Name:           personDto.Name,
		Surname:        personDto.Surname,
		Email:          personDto.Email,
		DateOfBirth:    personDto.DateOfBirth,
		Hobby:          personDto.Hobby,
		CreatedAt:      mockedCreationTime,
		LastModifiedAt: mockedCreationTime,
	}

	repository := repositoryMock{
		findPersonByEmail: func(ctx context.Context, email string) (model.Person, error) {
			if email != personDto.Email {
				t.Errorf("Expected email %s but got %s", personDto.Email, email)
			}
			return model.Person{}, db.ErrNoRows
		},
		insertPerson: func(ctx context.Context, person model.Person) (model.Person, error) {
			person.Id = expectedPerson.Id
			if person != expectedPerson {
				t.Errorf("Expected person like %#v but got %#v", expectedPerson, person)
			}
			return person, nil
		},
	}

	testObject := New(repository)
	result, err := testObject.CreateFromDto(mdctx.New(), personDto)
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if result != expectedPerson {
		t.Errorf("Expected person like %#v but got %#v", expectedPerson, result)
	}
}

func TestCreatePersonFromDtoShouldReturnErrorForDuplicateEmail(t *testing.T) {
	personDefinition := apimodel.PersonCreate{
		Name:        "John",
		Surname:     "Smith",
		Email:       "john.smith@test.com",
		DateOfBirth: util.Date(25, time.March, 23),
		Hobby:       "Literature",
	}
	personWithSameEmail := model.Person{
		Id:    "51ad7",
		Email: personDefinition.Email,
	}

	repository := repositoryMock{
		findPersonByEmail: func(ctx context.Context, email string) (model.Person, error) {
			return personWithSameEmail, nil
		},
		insertPerson: func(ctx context.Context, person model.Person) (model.Person, error) {
			t.Fatalf("Shouldn't be called because email is already in use")
			return model.Person{}, nil
		},
	}

	testObject := New(repository)
	_, err := testObject.CreateFromDto(mdctx.New(), personDefinition)
	if !testutil.IsBadInputError(err) {
		t.Errorf("Expected bad input status error but got %v", err)
	}
}

type repositoryMock struct {
	findPersonByEmail func(ctx context.Context, email string) (model.Person, error)
	insertPerson      func(ctx context.Context, person model.Person) (model.Person, error)
}

func (mock repositoryMock) FindPersonByEmail(ctx context.Context, email string) (model.Person, error) {
	return mock.findPersonByEmail(ctx, email)
}
func (mock repositoryMock) InsertPerson(ctx context.Context, person model.Person) (model.Person, error) {
	return mock.insertPerson(ctx, person)
}
