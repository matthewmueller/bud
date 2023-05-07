package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/matthewmueller/bud/db"
)

// Open a connection to a PostgreSQL database
func Open(url string) (db.DB, error) {
	return sql.Open("pgx/v5", url)
}

var _ db.DB = (*sql.DB)(nil)
