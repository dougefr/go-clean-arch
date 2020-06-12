package iinfra

import (
	"context"
	"database/sql"
)

// ContextKeyTx ...
const ContextKeyTx string = "ContextKeyTx"

type (
	// Tx ...
	Tx interface{}

	// Session ...
	Session interface {
		BeginTx() (Tx, error)
		CommitTx(tx Tx) error
		RollbackTx(tx Tx) error
	}

	// Database ...
	Database interface {
		Session
		Query(context.Context, string, ...interface{}) (*sql.Rows, error)
		Exec(context.Context, string, ...interface{}) (sql.Result, error)
	}
)
