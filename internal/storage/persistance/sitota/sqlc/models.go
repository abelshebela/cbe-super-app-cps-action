package sqlc

import (
	"database/sql"
)

type FindSitotaTransactionRow struct {
	ID                     string       `json:"id"`
	SenderName             string       `json:"sender_name"`
	SenderAccountNumber    string       `json:"sender_account_number"`
	SenderPhoneNumber      string       `json:"sender_phone_number"`
	RecipientName          string       `json:"recipient_name"`
	RecipientAccountNumber string       `json:"recipient_account_number"`
	RecipientPhoneNumber   string       `json:"recipient_phone_number"`
	SitotaAmount           float64      `json:"sitota_amount"`
	GLAccountNumber        string       `json:"gl_account_number"`
	Status                 string       `json:"status"`
	PaidAt                 sql.NullTime `json:"paid_at"`
	CreatedAt              sql.NullTime `json:"created_at"`
	UpdatedAt              sql.NullTime `json:"updated_at"`
	DeletedAt              sql.NullTime `json:"deleted_at,omitempty"`
	TotalCount             int64        `json:"total_count"`
}
