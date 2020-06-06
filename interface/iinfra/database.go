package iinfra

import (
	"context"
	"database/sql"
)

// Database ...
type Database interface {
	Session
	Query(context.Context, string, ...interface{}) (*sql.Rows, error)
	Exec(context.Context, string, ...interface{}) (sql.Result, error)
}
