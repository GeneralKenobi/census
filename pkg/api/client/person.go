package client

import (
	"context"
	"github.com/GeneralKenobi/census/pkg/api/apimodel"
)

func (api api) CreatePerson(ctx context.Context, person apimodel.PersonCreate) (apimodel.PersonCreated, error) {
	path := "/api/person"
	response, err := api.postPayload(ctx, path, person)
	if err != nil {
		return apimodel.PersonCreated{}, err
	}
	return processResponseWithPayload[apimodel.PersonCreated](response)
}

func (api api) GetPerson(ctx context.Context, id string) (apimodel.Person, error) {
	path := "/api/person/" + id
	response, err := api.get(ctx, path)
	if err != nil {
		return apimodel.Person{}, err
	}
	return processResponseWithPayload[apimodel.Person](response)
}

func (api api) UpdatePerson(ctx context.Context, id string, person apimodel.PersonUpdate) error {
	path := "/api/person/" + id
	response, err := api.putPayload(ctx, path, person)
	if err != nil {
		return err
	}
	return processResponse(response)
}

func (api api) DeletePerson(ctx context.Context, id string) error {
	path := "/api/person/" + id
	response, err := api.delete(ctx, path)
	if err != nil {
		return err
	}
	return processResponse(response)
}
