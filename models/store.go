package models

import (
	"context"
	"database/sql"
)

type SQL interface {
	Query
	BeginTx(context.Context, *sql.TxOptions) (Tx, error)
}

type Query interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Tx interface {
	Commit() error
	Rollback() error
	Query
}
