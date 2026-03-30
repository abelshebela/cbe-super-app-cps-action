package sqlc

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *WalletStorage {
	return &WalletStorage{db: db}
}

func (q *WalletStorage) WithTx(tx *sql.Tx) *WalletStorage {
	return &WalletStorage{
		db: tx,
	}
}
