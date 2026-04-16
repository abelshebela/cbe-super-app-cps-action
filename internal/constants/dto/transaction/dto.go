package transaction_dto

import (
	"time"
)

type VaultTransaction struct {
	// Primary Key & Identifiers
	ID            string `json:"id"`
	TransactionID string `json:"transaction_id"`
	FtNumber      string `json:"ft_number"`
	VaultID       string `json:"vault_id"`
	VaultTxType   string `json:"vault_tx_type"`

	// Financial Deltas & Balances
	PrincipalDelta string `json:"principal_delta"`
	InterestDelta  string `json:"interest_delta"`
	BalanceAfter   string `json:"balance_after"`

	// Reference Info
	ReferenceType           string `json:"reference_type"`
	ReferenceID             string `json:"reference_id"`
	DebitUserID             string `json:"debit_user_id"`
	DebitAccountNumber      string `json:"debit_account_number"`
	DebitAccountHolderName  string `json:"debit_account_holder_name"`
	CreditAccountNumber     string `json:"credit_account_number"`
	CreditAccountHolderName string `json:"credit_account_holder_name"`

	// Financial Amounts & Details
	Currency    string `json:"currency"`
	ServiceFee  string `json:"service_fee"`
	PaidAmount  string `json:"paid_amount"`
	VAT         string `json:"vat"`
	Amount      string `json:"amount"`
	TotalAmount string `json:"total_amount"`

	// Metadata
	ExternalReference string    `json:"external_reference"`
	TransactionReason string    `json:"transaction_reason"`
	ReceiptLink       string    `json:"receipt_link"`
	TransactionType   string    `json:"transaction_type"`
	IsIFB             int       `json:"is_ifb"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}
