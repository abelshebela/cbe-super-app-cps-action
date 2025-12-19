package archived_linked_account

import (
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// LinkedAccountMapper maps a LinkedAccount model to a bson.M for updates
func LinkedAccountMapper(account model.LinkedAccount) bson.M {
	return bson.M{
		"$set": bson.M{
			"user_id":             account.UserID,
			"customer_number":     account.CustomerNumber,
			"account_number":      account.AccountNumber,
			"account_holder_name": account.AccountHolderName,
			"account_type":        account.AccountType,
			"branch_code":         account.BranchCode,
			"registration_type":   account.RegistrationType,
			"and_or_status":       account.AndOrStatus,
			"currency_code":       account.CurrencyCode,
			"is_main":             account.IsMain,
			"updated_at":          time.Now(),
		},
	}
}
