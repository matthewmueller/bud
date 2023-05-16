package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Open a connection to a PostgreSQL database
func Open(url string) (*sql.DB, error) {
	return sql.Open("pgx/v5", url)
}
