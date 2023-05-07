package migrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/matthewmueller/bud/db"
	"github.com/matthewmueller/bud/log"
)

// ErrZerothMigration occurs when the migrations start at 000
var ErrZerothMigration = errors.New("migrations should start at 001 not 000")

// ErrNoMigrations happens when there are no migrations
var ErrNoMigrations = errors.New("no migrations")

// ErrNotEnoughMigrations happens when your migrations folder has less migrations than remote's version
var ErrNotEnoughMigrations = errors.New("remote migration version greater than the number of migrations you have")

type Migration interface {
	Up(ctx context.Context, db *sql.DB) error
	Down(ctx context.Context, db *sql.DB) error
}

type Status struct {
	Migrations []struct {
		Name    string
		Applied bool
	}
}

type Interface interface {
	Up(ctx context.Context) error
	UpBy(ctx context.Context, n int) error
	Down(ctx context.Context) error
	DownBy(ctx context.Context, n int) error
	Status(ctx context.Context) (*Status, error)
}

// func Migrate(m Migration)

// type Upper interface {
// 	Up(ctx context.Context, schema Schema) error
// }

// type Downer interface {
// 	Down(ctx context.Context, schema Schema) error
// }

// type Changer interface {
// 	Change(ctx context.Context, schema Schema) error
// }

// type UpDowner interface {
// 	Upper
// 	Downer
// }

func New(db db.DB, log log.Log) *Migrate {
	return &Migrate{db, log, nil}
}

// func Add(updown UpDowner) Migration {

// }

// type upDown struct {
// 	updown UpDown
// }

// func (u *upDown) Up(ctx context.Context) error {
// 	return u.updown.Up(ctx, nil)
// }

// func (u *upDown) Down(ctx context.Context) error {
// 	return u.updown.Down(ctx, nil)
// }

type Migrate struct {
	db         db.DB
	log        log.Log
	migrations []Migration
}

// Up migrates the database to the latest migration
func (m *Migrate) Up(ctx context.Context) error {

	m.log.Info("migrating up")
	return nil
}

// Down migrates the database back to the start
func (m *Migrate) Down(ctx context.Context) error {
	m.log.Info("migrating down")
	return nil
}

func AddGo(m Migration) Migration {
	// return m
	panic("Not implemented")
}

func (m *Migrate) AddSQL(mg Migration) Migration {
	// return mnil
	panic("Not implemented")
}

// // SQL migration
// func SQL(fsys fs.FS, path string) Migration {
// 	return nil
// }

// func Add(ud UpDowner) Migration {
// 	return &upDown{ud}
// }

// type upDown struct {
// 	updown UpDowner
// }

// func (u *upDown) Up(ctx context.Context, db *sql.DB) error {
// 	return u.updown.Up(ctx, nil)
// }

// func (u *upDown) Down(ctx context.Context, db *sql.DB) error {
// 	return u.updown.Down(ctx, nil)
// }

func (m *Migrate) Status(ctx context.Context) (*Status, error) {
	return nil, fmt.Errorf("unimplemented")
}

// ensure the table exists
// func ensureTableExists(db *sql.DB, tableName string) error {
// 	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS " + tableName + " (version bigint not null primary key);"); err != nil {
// 		return err
// 	}
// 	return nil
// }

// // Version gets the version from postgres
// func getRemoteVersion(db *sql.DB, tableName string) (version uint, err error) {
// 	err = db.QueryRow("SELECT version FROM " + tableName + " ORDER BY version DESC LIMIT 1").Scan(&version)
// 	switch {
// 	case err == sql.ErrNoRows:
// 		return 0, nil
// 	case err != nil:
// 		return 0, err
// 	default:
// 		return version, nil
// 	}
// }

// func (m)

// func File(fsys fs.FS, path string) Migration {

// }

// type sqlFile struct {
// 	db   *sql.DB
// 	fsys fs.FS
// 	path string
// }

// func (s *sqlFile) Up(ctx context.Context) error {
// 	return nil
// }

// func (s *sqlFile) Down(ctx context.Context) error {
// 	return nil
// }
