package repository

import (
	"context"
	"github.com/GeneralKenobi/census/internal/db/model"
)

func (repository *Repository) FindPersonById(ctx context.Context, id string) (model.Person, error) {
	idInt, err := idAsInt(id)
	if err != nil {
		return model.Person{}, err
	}
	return selectingOne(ctx, "find person by ID", repository.sql, personRowScanSupplier,
		"SELECT id::text, name, surname, email, date_of_birth, hobby, created_at, last_modified_at FROM censusdb.person WHERE id = $1", idInt)
}

func (repository *Repository) FindPersonByEmail(ctx context.Context, email string) (model.Person, error) {
	return selectingOne(ctx, "find person by email", repository.sql, personRowScanSupplier,
		"SELECT id::text, name, surname, email, date_of_birth, hobby, created_at, last_modified_at FROM censusdb.person WHERE email = $1", email)
}

func (repository *Repository) InsertPerson(ctx context.Context, person model.Person) (model.Person, error) {
	return selectingOne(ctx, "insert person", repository.sql, personRowScanSupplier,
		"INSERT INTO censusdb.person(name, surname, email, date_of_birth, hobby, created_at, last_modified_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id::text, name, surname, email, date_of_birth, hobby, created_at, last_modified_at",
		person.Name, person.Surname, person.Email, person.DateOfBirth, person.Hobby, person.CreatedAt, person.LastModifiedAt)
}

func (repository *Repository) UpdatePerson(ctx context.Context, person model.Person) (model.Person, error) {
	idInt, err := idAsInt(person.Id)
	if err != nil {
		return model.Person{}, err
	}
	return selectingOne(ctx, "update person", repository.sql, personRowScanSupplier,
		"UPDATE censusdb.person SET (name, surname, email, date_of_birth, hobby, created_at, last_modified_at) = ($2, $3, $4, $5, $6, $7, $8) WHERE id = $1 RETURNING id::text, name, surname, email, date_of_birth, hobby, created_at, last_modified_at",
		idInt, person.Name, person.Surname, person.Email, person.DateOfBirth, person.Hobby, person.CreatedAt, person.LastModifiedAt)
}

func (repository *Repository) DeletePersonById(ctx context.Context, id string) error {
	idInt, err := idAsInt(id)
	if err != nil {
		return err
	}
	return affectingOne(ctx, "delete person by ID", repository.sql,
		"DELETE FROM censusdb.person WHERE id = $1", idInt)
}

func personRowScanSupplier() (*model.Person, []any) {
	var person model.Person
	return &person, []any{
		&person.Id,
		&person.Name,
		&person.Surname,
		&person.Email,
		&person.DateOfBirth,
		&person.Hobby,
		&person.CreatedAt,
		&person.LastModifiedAt,
	}
}
