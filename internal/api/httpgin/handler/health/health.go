package health

import (
	"context"
	"github.com/GeneralKenobi/census/internal/api/httpgin/wrap"
	"github.com/GeneralKenobi/census/pkg/mdctx"
	"github.com/gin-gonic/gin"
)

// HandlerFunc always responds with HTTP200 (if the server is up the application is healthy).
func HandlerFunc(ginCtx *gin.Context) {
	wrap.Request(ginCtx).Handle(func(ctx context.Context, request wrap.RequestData) error {
		mdctx.Debugf(ctx, "Health probe")
		return nil
	})
}
