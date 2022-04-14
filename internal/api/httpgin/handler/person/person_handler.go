package person

import (
	"context"
	"github.com/GeneralKenobi/census/internal/api/httpgin/wrapper"
	"github.com/GeneralKenobi/census/internal/db"
	personcreator "github.com/GeneralKenobi/census/internal/service/person/creator"
	persondeleter "github.com/GeneralKenobi/census/internal/service/person/deleter"
	persongetter "github.com/GeneralKenobi/census/internal/service/person/getter"
	personupdater "github.com/GeneralKenobi/census/internal/service/person/updater"
	"github.com/GeneralKenobi/census/pkg/api/apimodel"
	"github.com/GeneralKenobi/census/pkg/mdctx"
	"github.com/gin-gonic/gin"
)

func NewHandler(querent db.Querent, transactioner db.Transactioner) *Handler {
	return &Handler{
		querent:       querent,
		transactioner: transactioner,
	}
}

type Handler struct {
	querent       db.Querent
	transactioner db.Transactioner
}

func (handler *Handler) CreateHandlerFunc(request *gin.Context) {
	wrapper.ForRequestRetV[apimodel.PersonCreated](request).Handle(func(ctx context.Context) (apimodel.PersonCreated, error) {
		ctx = mdctx.WithOperationName(ctx, "create person")
		return wrapper.WithBoundRequestBodyRetV(request, func(personDto apimodel.PersonCreate) (apimodel.PersonCreated, error) {
			return db.InTransactionRetV(ctx, handler.transactioner, func(repository db.Repository) (apimodel.PersonCreated, error) {
				personCreator := personcreator.New(repository)
				person, err := personCreator.CreateFromDto(ctx, personDto)
				if err != nil {
					return apimodel.PersonCreated{}, err
				}

				personCreated := apimodel.PersonCreated{Id: person.Id}
				return personCreated, nil
			})
		})
	})
}

func (handler *Handler) GetHandlerFunc(request *gin.Context) {
	wrapper.ForRequestRetV[apimodel.Person](request).Handle(func(ctx context.Context) (apimodel.Person, error) {
		ctx = mdctx.WithOperationName(ctx, "get person")
		return wrapper.WithRequiredPathParamRetV(request, "id", func(id string) (apimodel.Person, error) {
			return db.WithRepositoryRetV(ctx, handler.querent, func(repository db.Repository) (apimodel.Person, error) {
				personGetter := persongetter.New(repository)
				person, err := personGetter.FindById(ctx, id)
				if err != nil {
					return apimodel.Person{}, err
				}

				personDto := apimodel.Person{
					Id:             person.Id,
					Name:           person.Name,
					Surname:        person.Surname,
					Email:          person.Email,
					DateOfBirth:    person.DateOfBirth,
					Hobby:          person.Hobby,
					CreatedAt:      person.CreatedAt,
					LastModifiedAt: person.LastModifiedAt,
				}
				return personDto, nil
			})
		})
	})
}

func (handler *Handler) UpdateHandlerFunc(request *gin.Context) {
	wrapper.ForRequest(request).Handle(func(ctx context.Context) error {
		ctx = mdctx.WithOperationName(ctx, "update person")
		return wrapper.WithRequiredPathParam(request, "id", func(id string) error {
			return wrapper.WithBoundRequestBody(request, func(personDto apimodel.PersonUpdate) error {
				return db.InTransaction(ctx, handler.transactioner, func(repository db.Repository) error {
					personUpdater := personupdater.New(repository)
					_, err := personUpdater.UpdateFromDto(ctx, id, personDto)
					return err
				})
			})
		})
	})
}

func (handler *Handler) DeleteHandlerFunc(request *gin.Context) {
	wrapper.ForRequest(request).Handle(func(ctx context.Context) error {
		ctx = mdctx.WithOperationName(ctx, "delete person")
		return wrapper.WithRequiredPathParam(request, "id", func(id string) error {
			return db.InTransaction(ctx, handler.transactioner, func(repository db.Repository) error {
				personDeleter := persondeleter.New(repository)
				return personDeleter.DeleteById(ctx, id)
			})
		})
	})
}
