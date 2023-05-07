package sqlite

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
	"github.com/matthewmueller/bud/db"
)

func Open(url string) (db.DB, error) {
	return sql.Open("sqlite", url)
}

var _ db.DB = (*sql.DB)(nil)
