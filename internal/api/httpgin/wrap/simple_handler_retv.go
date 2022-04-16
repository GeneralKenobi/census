package wrap

import (
	"context"
	"github.com/GeneralKenobi/census/internal/api/httpgin/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RequestRetV creates a convenience standard handler which writes status HTTP200 and V marshalled to JSON on success and writes an error
// response on error. SimpleHandler.OnSuccess and SimpleHandler.OnError can be used to override the default success/error handlers.
// The request is processed by calling SimpleHandler.Handle which also accepts the handler function.
func RequestRetV[V any](ginCtx *gin.Context) SimpleHandlerRetV[V] {
	return &simpleHandlerRetV[V]{
		ginCtx: ginCtx,
		onSuccess: func(_ context.Context, handlerResult V) {
			ginCtx.JSON(http.StatusOK, handlerResult)
		},
		onError: func(ctx context.Context, handlerErr error) {
			request.WriteErrorResponse(ctx, ginCtx, handlerErr)
		},
	}
}

type simpleHandlerRetV[V any] struct {
	ginCtx    *gin.Context
	onSuccess func(ctx context.Context, handlerResult V)
	onError   func(ctx context.Context, handlerErr error)
}

var _ SimpleHandlerRetV[any] = (*simpleHandlerRetV[any])(nil) // Interface guard

func (handler *simpleHandlerRetV[V]) OnSuccess(onSuccess func(ctx context.Context, handlerResult V)) SimpleHandlerRetV[V] {
	handler.onSuccess = onSuccess
	return handler
}

func (handler *simpleHandlerRetV[V]) OnError(onError func(ctx context.Context, handlerErr error)) SimpleHandlerRetV[V] {
	handler.onError = onError
	return handler
}

func (handler *simpleHandlerRetV[V]) Handle(handlerFunc func(ctx context.Context, requestData RequestData) (V, error)) {
	ctx := request.Context(handler.ginCtx)
	requestData := requestDataParser{ginCtx: handler.ginCtx}

	value, err := handlerFunc(ctx, &requestData)

	if err == nil {
		handler.onSuccess(ctx, value)
	} else {
		handler.onError(ctx, err)
	}
}
