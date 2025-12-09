package transaction_dto

import (
	"encoding/json"
	"time"
)

type FullTransaction struct {
	ID                      string
	TransactionID           string
	FTNumber                string
	DebitBranchCode         string
	DebitDistrictCode       string
	DebitUserID             string
	DebitAccountNumber      string
	DebitAccountHolderName  string
	CreditUserID            string
	CreditAccountNumber     string
	CreditAccountHolderName string
	InstitutionCode         string
	InstitutionName         string
	Currency                string
	ServiceFee              string
	TipAmount               string
	PaidAmount              string
	VAT                     string
	Amount                  string
	TotalAmount             string
	ExternalReference       string
	TransactionReason       string
	TransactionType         string
	TransactionStatus       string
	IsIFB                   bool
	IsReversed              bool
	PaidAt                  *time.Time
	ReversedAt              *time.Time
	Metadata                json.RawMessage
	CreatedAt               *time.Time
	LastModifiedAt          *time.Time
}
