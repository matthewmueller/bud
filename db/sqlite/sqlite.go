package sqlite

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
	"github.com/livebud/buddy/db"
)

func Open(url string) (*sql.DB, error) {
	return sql.Open("sqlite", url)
}

var _ db.DB = (*sql.DB)(nil)
