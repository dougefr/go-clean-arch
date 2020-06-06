package sqlite3

import (
	"context"
	"database/sql"
	"github.com/dougefr/go-clean-arch/interface/iinfra"
)

type sqlite3 struct {
	db *sql.DB
}

// NewDatabase ...
func NewDatabase() (iinfra.Database, error) {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		return nil, err
	}

	return sqlite3{
		db: db,
	}, nil
}

// Query ...
func (s sqlite3) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if tx, ok := ctx.Value(iinfra.ContextKeyTx).(*sql.Tx); ok {
		return tx.Query(query, args...)
	}

	return s.db.Query(query, args...)
}

// Query ...
func (s sqlite3) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if tx, ok := ctx.Value(iinfra.ContextKeyTx).(*sql.Tx); ok {
		return tx.Exec(query, args...)
	}

	return s.db.Exec(query, args...)
}

// BeginTx ...
func (s sqlite3) BeginTx() (iinfra.Tx, error) {
	return s.db.Begin()
}

// CommitTx ...
func (s sqlite3) CommitTx(tx iinfra.Tx) error {
	if db, ok := tx.(*sql.Tx); ok {
		return db.Commit()
	}
	return nil
}

// RollbackTx ...
func (s sqlite3) RollbackTx(tx iinfra.Tx) error {
	if db, ok := tx.(*sql.Tx); ok {
		return db.Rollback()
	}
	return nil
}
