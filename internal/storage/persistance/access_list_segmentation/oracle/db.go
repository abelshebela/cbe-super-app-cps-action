package access_list_segmentation_oracle

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *accessListSegmentationOracle {
	return &accessListSegmentationOracle{db: db}
}

func (q *accessListSegmentationOracle) WithTx(tx *sql.Tx) *accessListSegmentationOracle {
	return &accessListSegmentationOracle{
		db: tx,
	}
}
