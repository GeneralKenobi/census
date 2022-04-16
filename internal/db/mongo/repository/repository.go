package repository

import (
	"context"
	"github.com/GeneralKenobi/census/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

// QueryContext abstract common mongo driver functions for standalone and session-based queries.
type QueryContext interface {
	Database() *mongo.Database
	// QueryContext enhances ctx with session information if a session is associated with this QueryContext.
	QueryContext(ctx context.Context) context.Context
}

func New(queryCtx QueryContext) *Repository {
	return &Repository{queryCtx: queryCtx}
}

type Repository struct {
	queryCtx QueryContext
}

var _ db.Repository = (*Repository)(nil) // Interface guard
