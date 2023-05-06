package db

import (
	"context"
	"database/sql"
)

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
