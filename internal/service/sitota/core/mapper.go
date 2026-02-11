package core

// func MapTransactionToSitota(tx any) *model.SitotaTransaction {
// 	sitota := &model.SitotaTransaction{}

// 	switch t := tx.(type) {
// 	case *transactionpb.Transaction:
// 		sitota.ID = t.GetTransactionId()
// 		sitota.SenderName = t.GetDebitAccountHolderName()
// 		sitota.SenderAccountNumber = t.GetDebitAccountNumber()
// 		sitota.SenderPhoneNumber = ""

// 		sitota.RecipientName = t.GetCreditAccountHolderName()
// 		sitota.RecipientAccountNumber = t.GetCreditAccountNumber()
// 		sitota.SenderPhoneNumber = ""

// 		sitota.Status = t.GetTransactionStatus()
// 		sitota.GLAccountNumber = ""

// 		if amountStr := t.GetAmount(); amountStr != "" {
// 			if amount, err := strconv.ParseFloat(amountStr, 64); err == nil {
// 				sitota.SitotaAmount = amount
// 			}
// 		}

// 		if createdAtStr := t.GetCreatedAt(); createdAtStr != "" {
// 			if tm, err := parseTime(createdAtStr); err == nil {
// 				sitota.CreatedAt = tm
// 			}
// 		}
// 		if lastModifiedStr := t.GetLastModifiedAt(); lastModifiedStr != "" {
// 			if tm, err := parseTime(lastModifiedStr); err == nil {
// 				sitota.UpdatedAt = tm
// 			}
// 		}
// 		if paidAtStr := t.GetPaidAt(); paidAtStr != "" {
// 			if tm, err := parseTime(paidAtStr); err == nil {
// 				sitota.ClaimedAt = &tm
// 			}
// 		}

// 	case *transactionpb.TransactionsList:
// 		sitota.ID = t.GetTransactionId()
// 		sitota.SenderName = t.GetDebitAccountHolderName()
// 		sitota.SenderAccountNumber = t.GetDebitAccountNumber()

// 		sitota.RecipientName = t.GetCreditAccountHolderName()
// 		sitota.RecipientAccountNumber = t.GetCreditAccountNumber()

// 		sitota.Status = t.GetTransactionStatus()
// 		sitota.GLAccountNumber = ""

// 		if totalStr := t.GetTotalAmount(); totalStr != "" {
// 			if total, err := strconv.ParseFloat(totalStr, 64); err == nil {
// 				sitota.SitotaAmount = total
// 			}
// 		}
// 		if createdAtStr := t.GetCreatedAt(); createdAtStr != "" {
// 			if tm, err := parseTime(createdAtStr); err == nil {
// 				sitota.CreatedAt = tm
// 			}
// 		}
// 	default:
// 		return nil
// 	}

// 	return sitota
// }

// func parseTime(timeStr string) (time.Time, error) {
// 	formats := []string{
// 		time.RFC3339,
// 		time.RFC3339Nano,
// 		"2006-01-02T15:04:05Z",
// 		"2006-01-02T15:04:05.000Z",
// 		"2006-01-02 15:04:05",
// 		"2006-01-02",
// 	}

// 	for _, format := range formats {
// 		if t, err := time.Parse(format, timeStr); err == nil {
// 			return t, nil
// 		}
// 	}

// 	return time.Time{}, errors.New("unable to parse time string")
// }
