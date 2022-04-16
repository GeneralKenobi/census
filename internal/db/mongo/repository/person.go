package repository

import (
	"context"
	"github.com/GeneralKenobi/census/internal/db/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (repository Repository) FindPersonById(ctx context.Context, id string) (model.Person, error) {
	ctx = repository.queryCtx.QueryContext(ctx)
	objectId, err := idAsObjectId(id)
	if err != nil {
		return model.Person{}, err
	}

	result := repository.personCollection().FindOne(ctx, filterById(objectId))
	return decodeSingleResult[personBson, model.Person]("find person by ID", result, personFromBson)
}

func (repository Repository) FindPersonByEmail(ctx context.Context, email string) (model.Person, error) {
	ctx = repository.queryCtx.QueryContext(ctx)

	result := repository.personCollection().FindOne(ctx, bson.M{"email": email})
	return decodeSingleResult[personBson, model.Person]("find person by email", result, personFromBson)
}

func (repository Repository) InsertPerson(ctx context.Context, person model.Person) (model.Person, error) {
	ctx = repository.queryCtx.QueryContext(ctx)
	person.Id = ""

	result, err := repository.personCollection().InsertOne(ctx, personToBson(person))
	insertedDocumentId, err := processInsertOneResult("insert person", result, err)
	if err != nil {
		return model.Person{}, err
	}

	person.Id = insertedDocumentId.Hex()
	return person, nil
}

func (repository Repository) UpdatePerson(ctx context.Context, person model.Person) (model.Person, error) {
	ctx = repository.queryCtx.QueryContext(ctx)
	objectId, err := idAsObjectId(person.Id)
	if err != nil {
		return model.Person{}, err
	}

	result, err := repository.personCollection().UpdateOne(ctx, filterById(objectId), updateBySet(personToBson(person)))
	err = processUpdateOneResult("update person", result, err)
	if err != nil {
		return model.Person{}, err
	}

	return person, nil
}

func (repository Repository) DeletePersonById(ctx context.Context, id string) error {
	ctx = repository.queryCtx.QueryContext(ctx)
	objectId, err := idAsObjectId(id)
	if err != nil {
		return err
	}

	result, err := repository.personCollection().DeleteOne(ctx, filterById(objectId))
	return processDeleteOneResult("delete person by ID", result, err)
}

func (repository Repository) personCollection() *mongo.Collection {
	return repository.queryCtx.Database().Collection("person")
}

// personBson is a BSON-mapped representation of model.Person.
type personBson struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`
	Name           string             `bson:"name,omitempty"`
	Surname        string             `bson:"surname,omitempty"`
	Email          string             `bson:"email,omitempty"`
	DateOfBirth    time.Time          `bson:"dateOfBirth,omitempty"`
	Hobby          string             `bson:"hobby,omitempty"`
	CreatedAt      time.Time          `bson:"createdAt,omitempty"`
	LastModifiedAt time.Time          `bson:"LastModifiedAt,omitempty"`
}

func personToBson(person model.Person) personBson {
	return personBson{
		Id:             idAsObjectIdOptional(person.Id),
		Name:           person.Name,
		Surname:        person.Surname,
		Email:          person.Email,
		DateOfBirth:    person.DateOfBirth,
		Hobby:          person.Hobby,
		CreatedAt:      person.CreatedAt,
		LastModifiedAt: person.LastModifiedAt,
	}
}

func personFromBson(person personBson) model.Person {
	return model.Person{
		Id:             person.Id.Hex(),
		Name:           person.Name,
		Surname:        person.Surname,
		Email:          person.Email,
		DateOfBirth:    person.DateOfBirth,
		Hobby:          person.Hobby,
		CreatedAt:      person.CreatedAt,
		LastModifiedAt: person.LastModifiedAt,
	}
}
