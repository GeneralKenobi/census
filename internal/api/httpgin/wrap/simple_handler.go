package wrap

import (
	"context"
	"github.com/GeneralKenobi/census/internal/api/httpgin/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Request creates a convenience standard handler which writes status HTTP200 on success and writes an error response on error.
// SimpleHandler.OnSuccess and SimpleHandler.OnError can be used to override the default success/error handlers.
// The request is processed by calling SimpleHandler.Handle which also accepts the handler function.
func Request(ginCtx *gin.Context) SimpleHandler {
	return &simpleHandler{
		ginCtx: ginCtx,
		onSuccess: func(_ context.Context) {
			ginCtx.Status(http.StatusOK)
		},
		onError: func(ctx context.Context, handlerErr error) {
			request.WriteErrorResponse(ctx, ginCtx, handlerErr)
		},
	}
}

type simpleHandler struct {
	ginCtx    *gin.Context
	onSuccess func(ctx context.Context)
	onError   func(ctx context.Context, handlerErr error)
}

var _ SimpleHandler = (*simpleHandler)(nil) // Interface guard

func (handler *simpleHandler) OnSuccess(onSuccess func(ctx context.Context)) SimpleHandler {
	handler.onSuccess = onSuccess
	return handler
}

func (handler *simpleHandler) OnError(onError func(ctx context.Context, handlerErr error)) SimpleHandler {
	handler.onError = onError
	return handler
}

func (handler *simpleHandler) Handle(handlerFunc func(ctx context.Context, requestData RequestData) error) {
	ctx := request.Context(handler.ginCtx)
	requestData := requestDataParser{ginCtx: handler.ginCtx}

	err := handlerFunc(ctx, &requestData)

	if err == nil {
		handler.onSuccess(ctx)
	} else {
		handler.onError(ctx, err)
	}
}
