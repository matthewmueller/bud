package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/matthewmueller/bud/db/postgres"
	"github.com/matthewmueller/bud/db/sqlite"
	"github.com/matthewmueller/bud/di"
	"github.com/xo/dburl"
)

func Provider(in di.Injector) {
	di.Provide[URL](in, provideURL)
	di.Provide[DB](in, provideDB)
	di.Provide[*pgx.Conn](in, providePGX)
}

func provideURL(in di.Injector) (URL, error) {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return URL(url), nil
	}
	return defaultURL, nil
}

func provideDB(in di.Injector) (DB, error) {
	url, err := di.Load[URL](in)
	if err != nil {
		return nil, err
	}
	if url == "" {
		return nil, fmt.Errorf("missing db url")
	}
	return Open(context.Background(), string(url))
}

func providePGX(in di.Injector) (*pgx.Conn, error) {
	url, err := di.Load[URL](in)
	if err != nil {
		return nil, err
	}
	if url == "" {
		return nil, fmt.Errorf("missing db url")
	}
	return pgx.Connect(context.Background(), string(url))
}

// URL is a database connection string
type URL string

var defaultURL URL = "sqlite://:memory:"

func Open(ctx context.Context, url string) (*sql.DB, error) {
	u, err := dburl.Parse(url)
	if err != nil {
		return nil, err
	}
	switch u.Driver {
	case "sqlite3":
		return sqlite.Open(u.DSN)
	case "postgres":
		return postgres.Open(u.DSN)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", u.Driver)
	}
}

var _ DB = (*sql.DB)(nil)

// DB is an interface to a sql database. It is a wrapper for the golang sql/db builtin
type DB interface {
	Stmt
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	Close() error
}

// Stmt is a sql prepared statement
type Stmt interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

// Tx is a database transaction
type Tx interface {
	Stmt
	Commit() error
	Rollback() error
}

// Rows is an iterator for sql.Query results
type Rows interface {
	Close() error
	Err() error
	Next() bool
	Row
}

// Rows is a result for for sql.QueryRow results
type Row interface {
	Scan(dest ...interface{}) error
}
