package wrap

import "context"

// SimpleHandler wraps a request handler that doesn't return a body in its response. Use Request for instantiating it.
type SimpleHandler interface {
	OnSuccess(onSuccess func(ctx context.Context)) SimpleHandler
	OnError(onError func(ctx context.Context, handlerErr error)) SimpleHandler
	Handle(handler func(ctx context.Context, request RequestData) error)
}

// SimpleHandlerRetV wraps a request handler that returns a body in its response. Use RequestRetV for instantiating it.
type SimpleHandlerRetV[V any] interface {
	OnSuccess(onSuccess func(ctx context.Context, handlerResult V)) SimpleHandlerRetV[V]
	OnError(onError func(ctx context.Context, handlerErr error)) SimpleHandlerRetV[V]
	Handle(handler func(ctx context.Context, request RequestData) (V, error))
}

// RequestData is a wrapper around the gin.Context instance wrapped in a SimpleHandler or SimpleHandlerRetV. It contains utility methods for
// getting data out of the request - binding request body, getting path params, etc.
type RequestData interface {
	BindBody(target any) error
	RequirePathParam(paramName string) (string, error)
	RequireIntPathParam(paramName string) (int, error)
}
