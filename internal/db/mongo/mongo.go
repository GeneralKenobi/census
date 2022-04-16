package mongo

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/census/internal/config"
	"github.com/GeneralKenobi/census/internal/db"
	"github.com/GeneralKenobi/census/pkg/mdctx"
	"github.com/GeneralKenobi/census/pkg/shutdown"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// NewContext creates a mongo DB context. The DB client is closed when context is canceled.
func NewContext(shutdownCtx shutdown.Context) (*Context, error) {
	ctx, cancel := contextWithTimeout()
	defer cancel()

	clientOptions := options.Client().ApplyURI(connectionString())
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("connection configuration is invalid: %w", err)
	}

	dbCtx := Context{
		client:       client,
		databaseName: config.Get().Mongo.Database,
	}
	go disconnectClientOnContextCancellation(shutdownCtx, client)
	return &dbCtx, nil
}

func connectionString() string {
	cfg := config.Get().Mongo
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/?ssl=%t&replicaSet=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.VerifyTls, cfg.ReplicaSet)
}

func disconnectClientOnContextCancellation(ctx shutdown.Context, client *mongo.Client) {
	defer ctx.Notify()

	<-ctx.Done()
	mdctx.Infof(nil, "DB context canceled")
	disconnectClient(client)
}

func disconnectClient(client *mongo.Client) {
	mdctx.Infof(nil, "Shutting down DB connection")

	ctx, cancel := contextWithTimeout()
	defer cancel()
	err := client.Disconnect(ctx)
	if err != nil {
		mdctx.Errorf(nil, "Error closing DB connection: %v", err)
	}

	mdctx.Infof(nil, "DB connection closed")
}

func contextWithTimeout() (context.Context, context.CancelFunc) {
	timeout := time.Duration(config.Get().Mongo.TimeoutSeconds) * time.Second
	return context.WithTimeout(context.Background(), timeout)
}

// Context implements DB integration for mongo.
type Context struct {
	client       *mongo.Client
	databaseName string
}

var _ db.Context = (*Context)(nil) // Interface guard

func (mongoCtx *Context) database() *mongo.Database {
	return mongoCtx.client.Database(mongoCtx.databaseName)
}
