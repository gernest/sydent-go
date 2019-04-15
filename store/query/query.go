package query

import (
	"context"
	"database/sql"

	"github.com/gernest/sydent-go/models"
)

var _ models.SQL = (*SQL)(nil)

// SQL implements models.SQL interface based on *sql.DB adding instrumentation
// and logging.
type SQL struct {
	db *sql.DB
}

func New(db *sql.DB) *SQL {
	return &SQL{db: db}
}

func (s *SQL) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *SQL) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *SQL) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return s.db.QueryRowContext(ctx, query, args...)
}

func (s *SQL) BeginTx(ctx context.Context, opts *sql.TxOptions) (models.Tx, error) {
	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

type Tx struct {
	tx *sql.Tx
}

func (t *Tx) Commit() error {
	return t.tx.Commit()
}

func (t *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

func (t *Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

func (t *Tx) Rollback() error {
	return t.tx.Rollback()
}
