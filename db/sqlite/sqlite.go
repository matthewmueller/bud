package sqlite

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

func Open(url string) (*sql.DB, error) {
	return sql.Open("sqlite", url)
}
