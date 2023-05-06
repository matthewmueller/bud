package migrate_test

import (
	"context"
	"os"
	"testing"

	"github.com/livebud/buddy/db"
	"github.com/livebud/buddy/db/postgres"
	"github.com/livebud/buddy/db/sqlite"
	"github.com/matryer/is"
	"github.com/tj/assert"
)

const tableName = "migrate"

func TestPostgres(t *testing.T) {
	url := "postgres://localhost:5432/migrate-test?sslmode=disable"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)
			db, err := postgres.Open(url)
			is.NoErr(err)
			defer db.Close()
			test.fn(t, db)
		})
	}
}

func TestSQLite(t *testing.T) {
	is := is.New(t)
	url := "sqlite:///tmp.db"
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			is := is.New(t)
			is.NoErr(os.RemoveAll("./tmp.db"))
			db, err := sqlite.Open(url)
			is.NoErr(err)
			defer db.Close()
			test.fn(t, db)
		})
	}
}

// func connect(t testing.TB, url string) (*sql.DB, func()) {
// 	sql, err := db.Open(url)
// 	is.NoErr(err)
// 	return db, func() {
// 		is.NoErr(t, db.Close())
// 	}
// }

func drop(t testing.TB, db db.DB) {
	_, err := db.QueryContext(context.Background(), `
		drop table if exists migrate;
		drop table if exists users;
		drop table if exists teams;
	`)
	assert.NoError(t, err)
}

// func exists(t testing.TB, path string) {
// 	_, err := os.Stat(path)
// 	assert.NoError(t, err)
// }

// func notExists(err error, name string) bool {
// 	return strings.Contains(err.Error(), fmt.Sprintf("relation \"%s\" does not exist", name)) ||
// 		strings.Contains(err.Error(), fmt.Sprintf("no such table: %s", name))
// }

// func syntaxError(err error, name string) bool {
// 	return strings.Contains(err.Error(), fmt.Sprintf(`syntax error at or near "%s"`, name)) ||
// 		strings.Contains(err.Error(), fmt.Sprintf(`near "%s": syntax error`, name))
// }

var tests = []struct {
	name string
	fn   func(t testing.TB, db db.DB)
}{
	// {
	// 	name: "no migrations",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, db)
	// 		fs := fstest.MapFS{}
	// 		db, close := connect(t, url)
	// 		defer close()
	// 		err := migrate.Up(l, db, fs, tableName)
	// 		assert.Equal(t, migrate.ErrNoMigrations, err)
	// 	},
	// },
	// {
	// 	name: "no matching migrations",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)
	// 		fs := fstest.MapFS{
	// 			"migrate/001_init.up.sql":   {Data: []byte(``)},
	// 			"migrate/001_init.down.sql": {Data: []byte(``)},
	// 		}
	// 		db, close := connect(t, url)
	// 		defer close()
	// 		err := migrate.Up(l, db, fs, tableName)
	// 		assert.Equal(t, migrate.ErrNoMigrations, err)
	// 	},
	// },
	// {
	// 	name: "no migrations down",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)
	// 		fs := fstest.MapFS{}
	// 		db, close := connect(t, url)
	// 		defer close()
	// 		err := migrate.Down(l, db, fs, tableName)
	// 		assert.Equal(t, migrate.ErrNoMigrations, err)
	// 	},
	// },
	// {
	// 	name: "no matching migrations down",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)
	// 		fs := fstest.MapFS{
	// 			"migrate/001_init.up.sql":   {Data: []byte(``)},
	// 			"migrate/001_init.down.sql": {Data: []byte(``)},
	// 		}
	// 		db, close := connect(t, url)
	// 		defer close()
	// 		err := migrate.Down(l, db, fs, tableName)
	// 		assert.Equal(t, migrate.ErrNoMigrations, err)
	// 	},
	// },
	// {
	// 	name: "zeroth",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)
	// 		fs := fstest.MapFS{
	// 			"000_init.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists teams (
	// 						id serial primary key not null,
	// 						name text not null
	// 					);
	// 					create table if not exists users (
	// 						id serial primary key not null,
	// 						email text not null,
	// 						created_at time with time zone not null,
	// 						updated_at time with time zone not null
	// 					);
	// 				`),
	// 			},
	// 			"000_init.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists users;
	// 					drop table if exists teams;
	// 				`),
	// 			},
	// 		}
	// 		db, close := connect(t, url)
	// 		defer close()
	// 		err := migrate.Up(l, db, fs, tableName)
	// 		assert.Equal(t, migrate.ErrZerothMigration, err)
	// 	},
	// },
	// {
	// 	name: "up down",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)

	// 		fs := fstest.MapFS{
	// 			"001_init.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists teams (
	// 						id serial primary key not null,
	// 						name text not null
	// 					);
	// 					create table if not exists users (
	// 						id serial primary key not null,
	// 						email text not null,
	// 						created_at time with time zone not null,
	// 						updated_at time with time zone not null
	// 					);
	// 				`),
	// 			},
	// 			"001_init.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists users;
	// 					drop table if exists teams;
	// 				`),
	// 			},
	// 		}

	// 		db, close := connect(t, url)
	// 		defer close()

	// 		err := migrate.Up(l, db, fs, tableName)
	// 		assert.NoError(t, err)

	// 		rows, err := db.Query(`insert into teams (id, name) values (1, 'jack')`)
	// 		assert.NoError(t, err)
	// 		for rows.Next() {
	// 			var id int
	// 			var name string
	// 			err := rows.Scan(&id, &name)
	// 			assert.NoError(t, err)
	// 			assert.Equal(t, 1, id)
	// 			assert.Equal(t, "jack", name)
	// 		}
	// 		assert.NoError(t, rows.Err())

	// 		err = migrate.Down(l, db, fs, tableName)
	// 		assert.NoError(t, err)

	// 		_, err = db.Query(`insert into teams (id, name) values (2, 'jack')`)
	// 		assert.NotNil(t, err)
	// 		assert.Contains(t, err.Error(), "teams")
	// 		assert.True(t, notExists(err, "teams"), err.Error())
	// 	},
	// },
	// {
	// 	name: "up down no logger",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)

	// 		fs := fstest.MapFS{
	// 			"001_init.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists teams (
	// 						id serial primary key not null,
	// 						name text not null
	// 					);
	// 					create table if not exists users (
	// 						id serial primary key not null,
	// 						email text not null,
	// 						created_at time with time zone not null,
	// 						updated_at time with time zone not null
	// 					);
	// 				`),
	// 			},
	// 			"001_init.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists users;
	// 					drop table if exists teams;
	// 				`),
	// 			},
	// 		}

	// 		db, close := connect(t, url)
	// 		defer close()

	// 		err := migrate.Up(nil, db, fs, tableName)
	// 		assert.NoError(t, err)

	// 		rows, err := db.Query(`insert into teams (id, name) values (1, 'jack')`)
	// 		assert.NoError(t, err)
	// 		for rows.Next() {
	// 			var id int
	// 			var name string
	// 			err := rows.Scan(&id, &name)
	// 			assert.NoError(t, err)
	// 			assert.Equal(t, 1, id)
	// 			assert.Equal(t, "jack", name)
	// 		}
	// 		assert.NoError(t, rows.Err())

	// 		err = migrate.Down(nil, db, fs, tableName)
	// 		assert.NoError(t, err)

	// 		_, err = db.Query(`insert into teams (id, name) values (2, 'jack')`)
	// 		assert.NotNil(t, err)
	// 		assert.Contains(t, err.Error(), "teams")
	// 		assert.True(t, notExists(err, "teams"), err.Error())
	// 	},
	// },
	// {
	// 	name: "upupdowndown",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)

	// 		fs := fstest.MapFS{
	// 			"001_init.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists teams (
	// 						id serial primary key not null,
	// 						name text not null
	// 					);
	// 				`),
	// 			},
	// 			"001_init.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists teams;
	// 				`),
	// 			},
	// 			"002_users.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists users (
	// 						id serial primary key not null,
	// 						email text not null
	// 					);
	// 				`),
	// 			},
	// 			"002_users.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists users;
	// 				`),
	// 			},
	// 		}

	// 		db, close := connect(t, url)
	// 		defer close()

	// 		err := migrate.Up(l, db, fs, tableName)
	// 		assert.NoError(t, err)
	// 		err = migrate.Up(l, db, fs, tableName)
	// 		assert.NoError(t, err)

	// 		rows, err := db.Query(`insert into users (id, email) values (1, 'jack')`)
	// 		assert.NoError(t, err)
	// 		for rows.Next() {
	// 			var id int
	// 			var email string
	// 			err := rows.Scan(&id, &email)
	// 			assert.NoError(t, err)
	// 			assert.Equal(t, 1, id)
	// 			assert.Equal(t, "jack", email)
	// 		}
	// 		assert.NoError(t, rows.Err())

	// 		err = migrate.Down(l, db, fs, tableName)
	// 		assert.NoError(t, err)
	// 		err = migrate.Down(l, db, fs, tableName)
	// 		assert.NoError(t, err)

	// 		_, err = db.Query(`insert into users (id, email) values (2, 'jack')`)
	// 		assert.NotNil(t, err)
	// 		assert.Contains(t, err.Error(), "users")
	// 		assert.True(t, notExists(err, "users"), err.Error())
	// 	},
	// },
	// {
	// 	name: "upbydownby",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)

	// 		fs := fstest.MapFS{
	// 			"001_init.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists teams (
	// 						id serial primary key not null,
	// 						name text not null
	// 					);
	// 				`),
	// 			},
	// 			"001_init.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists teams;
	// 				`),
	// 			},
	// 			"002_users.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists users (
	// 						id serial primary key not null,
	// 						email text not null
	// 					);
	// 				`),
	// 			},
	// 			"002_users.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists users;
	// 				`),
	// 			},
	// 		}

	// 		db, close := connect(t, url)
	// 		defer close()

	// 		err := migrate.UpBy(l, db, fs, tableName, 1)
	// 		assert.NoError(t, err)

	// 		_, err = db.Query(`insert into teams (name) values ('jack')`)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into users (email) values ('jack')`)
	// 		assert.NotNil(t, err)
	// 		assert.Contains(t, err.Error(), "users")
	// 		assert.True(t, notExists(err, "users"), err.Error())

	// 		err = migrate.UpBy(l, db, fs, tableName, 1)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into teams (name) values ('jack')`)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into users (email) values ('jack')`)
	// 		assert.NoError(t, err)

	// 		err = migrate.UpBy(l, db, fs, tableName, 1)
	// 		assert.NoError(t, err)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into teams (name) values ('jack')`)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into users (email) values ('jack')`)
	// 		assert.NoError(t, err)

	// 		err = migrate.DownBy(l, db, fs, tableName, 1)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into teams (name) values ('jack')`)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into users (email) values ('jack')`)
	// 		assert.NotNil(t, err)
	// 		assert.Contains(t, err.Error(), "users")
	// 		assert.True(t, notExists(err, "users"), err.Error())

	// 		err = migrate.DownBy(l, db, fs, tableName, 1)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into teams (name) values ('jack')`)
	// 		assert.Contains(t, err.Error(), "teams")
	// 		assert.True(t, notExists(err, "teams"), err.Error())
	// 		_, err = db.Query(`insert into users (email) values ('jack')`)
	// 		assert.Contains(t, err.Error(), "users")
	// 		assert.True(t, notExists(err, "users"), err.Error())

	// 		err = migrate.DownBy(l, db, fs, tableName, 1)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into teams (name) values ('jack')`)
	// 		assert.Contains(t, err.Error(), "teams")
	// 		assert.True(t, notExists(err, "teams"), err.Error())
	// 		_, err = db.Query(`insert into users (email) values ('jack')`)
	// 		assert.Contains(t, err.Error(), "users")
	// 		assert.True(t, notExists(err, "users"), err.Error())
	// 	},
	// },
	// {
	// 	name: "uprollback",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)

	// 		fs := fstest.MapFS{
	// 			"001_init.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists teams (
	// 						id serial primary key not null,
	// 						name text not null
	// 					);
	// 				`),
	// 			},
	// 			"001_init.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists teams;
	// 				`),
	// 			},
	// 			"002_users.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists users (
	// 						id serial primary key not null -- intentionally missing comma
	// 						email text not null
	// 					);
	// 				`),
	// 			},
	// 			"002_users.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists users;
	// 				`),
	// 			},
	// 		}

	// 		db, close := connect(t, url)
	// 		defer close()

	// 		err := migrate.UpBy(l, db, fs, tableName, 1)
	// 		assert.NoError(t, err)

	// 		_, err = db.Query(`insert into teams (name) values ('jack')`)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into users (email) values ('jack')`)
	// 		assert.NotNil(t, err)
	// 		assert.Contains(t, err.Error(), "users")
	// 		assert.True(t, notExists(err, "users"), err.Error())

	// 		err = migrate.UpBy(l, db, fs, tableName, 1)
	// 		assert.NotNil(t, err)
	// 		assert.True(t, syntaxError(err, "email"), err.Error())

	// 		_, err = db.Query(`insert into teams (name) values ('jack')`)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into users (email) values ('jack')`)
	// 		assert.Contains(t, err.Error(), "users")
	// 		assert.True(t, notExists(err, "users"), err.Error())
	// 	},
	// },
	// {
	// 	name: "downrollback",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)

	// 		fs := fstest.MapFS{
	// 			"001_init.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists teams (
	// 						id serial primary key not null,
	// 						name text not null
	// 					);
	// 				`),
	// 			},
	// 			"001_init.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exis teams; -- intentional syntax error
	// 				`),
	// 			},
	// 			"002_users.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists users (
	// 						id serial primary key not null,
	// 						email text not null
	// 					);
	// 				`),
	// 			},
	// 			"002_users.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists users;
	// 				`),
	// 			},
	// 		}

	// 		db, close := connect(t, url)
	// 		defer close()

	// 		// setup
	// 		err := migrate.Up(l, db, fs, tableName)
	// 		assert.NoError(t, err)

	// 		err = migrate.DownBy(l, db, fs, tableName, 1)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into teams (name) values ('jack')`)
	// 		assert.NoError(t, err)
	// 		_, err = db.Query(`insert into users (email) values ('jack')`)
	// 		assert.NotNil(t, err)
	// 		assert.Contains(t, err.Error(), "users")
	// 		assert.True(t, notExists(err, "users"))

	// 		err = migrate.DownBy(l, db, fs, tableName, 1)
	// 		assert.NotNil(t, err)
	// 		assert.True(t, syntaxError(err, "exis"), err.Error())

	// 		_, err = db.Query(`insert into teams (name) values ('jack')`)
	// 		assert.NoError(t, err)
	// 	},
	// },
	// {
	// 	name: "new",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)

	// 		// cleanup
	// 		assert.NoError(t, os.RemoveAll("migrate"))

	// 		err := migrate.New(l, "migrate", "setup")
	// 		assert.NoError(t, err)
	// 		exists(t, "migrate/001_setup.up.sql")
	// 		exists(t, "migrate/001_setup.down.sql")

	// 		err = migrate.New(l, "migrate", "create teams")
	// 		assert.NoError(t, err)
	// 		exists(t, "migrate/002_create_teams.up.sql")
	// 		exists(t, "migrate/002_create_teams.down.sql")

	// 		err = migrate.New(l, "migrate", "new-users")
	// 		assert.NoError(t, err)
	// 		exists(t, "migrate/003_new_users.up.sql")
	// 		exists(t, "migrate/003_new_users.down.sql")

	// 		if !t.Failed() {
	// 			assert.NoError(t, os.RemoveAll("migrate"))
	// 		}
	// 	},
	// },
	// {
	// 	name: "remoteversion",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)

	// 		fs := fstest.MapFS{
	// 			"001_init.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists teams (
	// 						id serial primary key not null,
	// 						name text not null
	// 					);
	// 				`),
	// 			},
	// 			"001_init.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists teams;
	// 				`),
	// 			},
	// 			"002_users.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists users (
	// 						id serial primary key not null,
	// 						email text not null
	// 					);
	// 				`),
	// 			},
	// 			"002_users.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists users;
	// 				`),
	// 			},
	// 		}

	// 		db, close := connect(t, url)
	// 		defer close()

	// 		// setup
	// 		err := migrate.Up(l, db, fs, tableName)
	// 		assert.NoError(t, err)

	// 		name, err := migrate.RemoteVersion(db, fs, tableName)
	// 		assert.NoError(t, err)
	// 		assert.Equal(t, `002_users.up.sql`, name)

	// 		// teardown
	// 		err = migrate.Down(l, db, fs, tableName)
	// 		assert.NoError(t, err)

	// 		_, err = migrate.RemoteVersion(db, fs, tableName)
	// 		assert.Equal(t, migrate.ErrNoMigrations, err)
	// 	},
	// },
	// {
	// 	name: "localversion",
	// 	fn: func(t testing.TB, db db.DB) {
	// 		drop(t, url)

	// 		fs := fstest.MapFS{
	// 			"001_init.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists teams (
	// 						id serial primary key not null,
	// 						name text not null
	// 					);
	// 				`),
	// 			},
	// 			"001_init.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists teams;
	// 				`),
	// 			},
	// 			"002_users.up.sql": {
	// 				Data: []byte(`
	// 					create table if not exists users (
	// 						id serial primary key not null,
	// 						email text not null
	// 					);
	// 				`),
	// 			},
	// 			"002_users.down.sql": {
	// 				Data: []byte(`
	// 					drop table if exists users;
	// 				`),
	// 			},
	// 		}

	// 		name, err := migrate.LocalVersion(fs)
	// 		assert.NoError(t, err)
	// 		assert.Equal(t, `002_users.up.sql`, name)
	// 	},
	// },
}
