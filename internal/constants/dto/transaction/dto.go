package transaction_dto

import (
	"encoding/json"
	"time"
)

type FullTransaction struct {
	ID                      string          `json:"id"`
	TransactionID           string          `json:"transaction_id"`
	FTNumber                string          `json:"ft_number"`
	DebitBranchCode         string          `json:"debit_branch_code"`
	DebitDistrictCode       string          `json:"debit_district_code"`
	DebitUserID             string          `json:"debit_user_id"`
	DebitAccountNumber      string          `json:"debit_account_number"`
	DebitAccountHolderName  string          `json:"debit_account_holder_name"`
	CreditUserID            string          `json:"credit_user_id"`
	CreditAccountNumber     string          `json:"credit_account_number"`
	CreditAccountHolderName string          `json:"credit_account_holder_name"`
	InstitutionCode         string          `json:"institution_code"`
	InstitutionName         string          `json:"institution_name"`
	Currency                string          `json:"currency"`
	ServiceFee              string          `json:"service_fee"`
	TipAmount               string          `json:"tip_amount"`
	PaidAmount              string          `json:"paid_amount"`
	VAT                     string          `json:"vat"`
	Amount                  string          `json:"amount"`
	TotalAmount             string          `json:"total_amount"`
	ExternalReference       string          `json:"external_reference"`
	TransactionReason       string          `json:"transaction_reason"`
	TransactionType         string          `json:"transaction_type"`
	TransactionStatus       string          `json:"transaction_status"`
	IsIFB                   bool            `json:"is_ifb"`
	IsReversed              bool            `json:"is_reversed"`
	PaidAt                  *time.Time      `json:"paid_at"`
	ReversedAt              *time.Time      `json:"reversed_at"`
	Metadata                json.RawMessage `json:"metadata" swaggertype:"object" format:"byte"`
	CreatedAt               *time.Time      `json:"created_at"`
	LastModifiedAt          *time.Time      `json:"last_modified_at"`
}
