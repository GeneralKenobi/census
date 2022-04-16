package request

import (
	"github.com/GeneralKenobi/census/pkg/mdctx"
	"github.com/gin-gonic/gin"
)

// ContextMiddleware creates and saves in gin's context an MDC-enhanced context for the request.
func ContextMiddleware(ginCtx *gin.Context) {
	ctx := mdctx.New()
	if correlationId := ginCtx.GetHeader("X-Correlation-ID"); correlationId != "" {
		ctx = mdctx.WithCorrelationId(ctx, correlationId)
	}
	ctx = mdctx.WithRequestMethod(ctx, ginCtx.Request.Method)
	ctx = mdctx.WithRequestUri(ctx, ginCtx.Request.RequestURI)
	ctx = mdctx.WithClientIp(ctx, ginCtx.ClientIP())
	ginCtx.Set(requestContextKey, ctx)
	ginCtx.Next()
}

func LogRequestProcessingMiddleware(ginCtx *gin.Context) {
	ctx := Context(ginCtx)
	mdctx.Infof(ctx, "Begin processing")
	ginCtx.Next()
	mdctx.Debugf(ctx, "End processing")
}
