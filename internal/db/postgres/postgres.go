package postgres

import (
	"database/sql"
	"fmt"
	"github.com/GeneralKenobi/census/internal/config"
	"github.com/GeneralKenobi/census/internal/db"
	"github.com/GeneralKenobi/census/pkg/mdctx"
	"github.com/GeneralKenobi/census/pkg/shutdown"
	_ "github.com/lib/pq" // Postgres driver registration by import.
)

// NewContext creates a postgres DB context. The DB client is closed when context is canceled.
func NewContext(ctx shutdown.Context) (*Context, error) {
	sqlDb, err := sql.Open("postgres", connectionString())
	if err != nil {
		return nil, fmt.Errorf("connection configuration is invalid: %w", err)
	}

	dbCtx := Context{db: sqlDb}
	go shutdownDbOnContextCancellation(ctx, sqlDb)
	return &dbCtx, nil
}

func connectionString() string {
	cfg := config.Get().Postgres
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, sslMode(cfg.VerifyTls))
}

func sslMode(enable bool) string {
	if enable {
		return "verify-full"
	}
	return "disable"
}

func shutdownDbOnContextCancellation(ctx shutdown.Context, db *sql.DB) {
	defer ctx.Notify()

	<-ctx.Done()
	mdctx.Infof(nil, "DB context canceled")
	shutdownDb(db)
}

func shutdownDb(db *sql.DB) {
	mdctx.Infof(nil, "Shutting down DB connection")
	err := db.Close()
	if err != nil {
		mdctx.Errorf(nil, "Error closing DB connection: %v", err)
	}
	mdctx.Infof(nil, "DB connection closed")
}

// Context implements DB integration for postgres.
type Context struct {
	db *sql.DB
}

var _ db.Context = (*Context)(nil) // Interface guard
