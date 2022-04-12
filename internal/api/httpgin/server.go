package httpgin

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/census/internal/api/httpgin/handler/health"
	"github.com/GeneralKenobi/census/internal/api/httpgin/request"
	"github.com/GeneralKenobi/census/internal/config"
	"github.com/GeneralKenobi/census/internal/db"
	"github.com/GeneralKenobi/census/pkg/mdctx"
	"github.com/GeneralKenobi/census/pkg/shutdown"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewServer(dbCtx db.Context) *Server {
	server := Server{
		dbCtx: dbCtx,
	}
	server.configure()
	return &server
}

type Server struct {
	dbCtx      db.Context
	httpServer *http.Server
}

// Run starts the HTTP server and shuts it down gracefully when ctx is cancelled.
func (server *Server) Run(ctx shutdown.Context) {
	go server.listenAndServe()
	server.shutdownOnContextCancellation(ctx)
}

// configure creates a ready-to-use server and stores it in Server.httpServer.
func (server *Server) configure() {
	httpCfg := config.Get().HttpServer
	ginEngine := server.setupGinEngine()
	server.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", httpCfg.Port),
		Handler: ginEngine,
	}
}

// setupGinEngine configures routing, middleware and handlers.
func (server *Server) setupGinEngine() *gin.Engine {
	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery(), request.ContextMiddleware, request.LogRequestProcessingMiddleware)

	ginEngine.GET("/health", health.HandlerFunc)

	return ginEngine
}

func (server *Server) listenAndServe() {
	mdctx.Infof(nil, "Starting HTTP server on address %s", server.httpServer.Addr)
	err := server.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		mdctx.Errorf(nil, "HTTP server exited with error: %v", err)
	} else {
		mdctx.Infof(nil, "HTTP server exited")
	}
}

func (server *Server) shutdownOnContextCancellation(ctx shutdown.Context) {
	defer ctx.Notify()

	<-ctx.Done()
	mdctx.Infof(nil, "Context canceled - shutting down HTTP server")
	serverShutdownCtx, cancel := context.WithTimeout(context.Background(), ctx.Timeout())
	defer cancel()

	err := server.httpServer.Shutdown(serverShutdownCtx)
	if err != nil {
		mdctx.Errorf(nil, "Error shutting down HTTP server: %v", err)
	} else {
		mdctx.Infof(nil, "HTTP server shutdown completed")
	}
}
