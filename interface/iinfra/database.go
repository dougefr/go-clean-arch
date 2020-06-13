// Copyright (c) 2020. Douglas Rodrigues - All rights reserved.
// This file is licensed under the MIT License.
// License text available at https://opensource.org/licenses/MIT

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
