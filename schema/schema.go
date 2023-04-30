package schema

import (
	"context"
	"database/sql"
)

func New(db *sql.DB) *Schema {
	return &Schema{db}
}

type Schema struct {
	db *sql.DB
}

func Create(ctx context.Context, tableName string, fn func(t *Table)) error {
	return nil
}

type Table struct {
}

// type Schema interface {
// 	Create(name string) Table
// 	Delete(name string)
// 	Query(query string, args ...interface{}) error
// }

// type Table interface {
// 	ID(name string) IntColumn
// 	String(name string) StringColumn
// 	Email(name string) StringColumn
// 	Password(name string) StringColumn
// 	Int(name string) IntColumn
// 	Bool(name string) BoolColumn
// 	Float32(name string) Float32Column
// }

// type Column interface {
// 	Primary() Column
// 	Nullable() Column
// 	Unique() Column
// }

// type StringColumn interface {
// 	Column
// 	Default(value string) StringColumn
// }

// type IntColumn interface {
// 	Column
// 	Default(value int) IntColumn
// }

// type BoolColumn interface {
// 	Column
// 	Default(value bool) BoolColumn
// }

// type Float32Column interface {
// 	Column
// 	Default(value float32) Float32Column
// }
