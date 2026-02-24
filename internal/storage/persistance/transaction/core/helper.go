package transaction_core

import (
	transaction_dto "cbe-super-app-cps-action/internal/constants/dto/transaction"
	"encoding/json"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

func MapFullTransactionToNative(ft model.TransactionModel) transaction_dto.FullTransaction {
	var paidAt, reversedAt, createdAt, lastModifiedAt *time.Time
	if ft.PaidAt.Valid {
		paidAt = &ft.PaidAt.Time
	}
	if ft.ReversedAt.Valid {
		reversedAt = &ft.ReversedAt.Time
	}
	if ft.CreatedAt.Valid {
		createdAt = &ft.CreatedAt.Time
	}
	if ft.LastModifiedAt.Valid {
		lastModifiedAt = &ft.LastModifiedAt.Time
	}

	var metadataBytes json.RawMessage
	if ft.Metadata != nil {
		metadataBytes = ft.Metadata
	}

	fee := ft.ServiceFee
	TipAmount := ft.TipAmount
	paidAmount := ft.PaidAmount
	vat := ft.VAT
	amount := ft.Amount
	totalAmount := ft.TotalAmount

	return transaction_dto.FullTransaction{
		ID:                      ft.ID,
		TransactionID:           ft.TransactionID,
		FTNumber:                ft.FTNumber,
		DebitBranchCode:         ft.DebitBranchCode,
		DebitDistrictCode:       ft.DebitDistrictCode,
		DebitUserID:             ft.DebitUserID,
		DebitAccountNumber:      ft.DebitAccountNumber,
		DebitAccountHolderName:  ft.DebitAccountHolderName,
		CreditUserID:            ft.CreditUserID,
		CreditAccountNumber:     ft.CreditAccountNumber,
		CreditAccountHolderName: ft.CreditAccountHolderName,
		InstitutionCode:         ft.InstitutionCode,
		InstitutionName:         ft.InstitutionName,
		Currency:                string(ft.Currency),
		ServiceFee:              fee.String(),
		TipAmount:               TipAmount.String(),
		PaidAmount:              paidAmount.String(),
		VAT:                     vat.String(),
		Amount:                  amount.String(),
		TotalAmount:             totalAmount.String(),
		ExternalReference:       ft.ExternalReference,
		TransactionReason:       ft.TransactionReason,
		TransactionType:         string(ft.TransactionType),
		TransactionStatus:       string(ft.TransactionStatus),
		IsIFB:                   ft.IsIFB == "Y" || ft.IsIFB == "true" || ft.IsIFB == "1",
		PaidAt:                  paidAt,
		ReversedAt:              reversedAt,
		Metadata:                metadataBytes,
		CreatedAt:               createdAt,
		LastModifiedAt:          lastModifiedAt,
	}
}
