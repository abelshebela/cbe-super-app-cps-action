package customer_oracle

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *customerOracleRepository {
	return &customerOracleRepository{db: db}
}

func (q *customerOracleRepository) WithTx(tx *sql.Tx) *customerOracleRepository {
	return &customerOracleRepository{
		db: tx,
	}
}
