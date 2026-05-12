package donation_category_oracle

import (
	"context"
	"database/sql"
)

// DBTX is the minimal surface a repository uses from *sql.DB.
// Both *sql.DB and *sql.Tx satisfy it, enabling WithTx.
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
