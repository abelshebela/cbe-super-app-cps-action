package amount_based_auth_domain

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AmountBasedAuthRequest struct {
	Id        string `json:"id"`
	MinAmount int    `json:"minAmount"`
	MaxAmount int    `json:"maxAmount"`
	Method    Method `json:"method"`
}

type AmountBasedAuthRespose struct {
	Page            int         `json:"page"`
	AmountBasedAuth []*AuthTier `json:"amount_based_auth"`
	Limit           int         `json:"limit"`
	Total           int64       `json:"total"`
}

type UpdateAmountBasedAuth struct {
	Id        string `json:"id" bson:"id"`
	MinAmount int    `json:"min_amount" bson:"min_amount"`
	MaxAmount int    `json:"max_amount" bson:"max_amount"`
	Method    Method `json:"method" bson:"method"`
}

func (u UpdateAmountBasedAuth) Valiadate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Method, validation.Required, validation.In(OPEN, PIN, OTPANDPIN),
			validation.By(validateRequest(u)),
		))
}

func validateRequest(u UpdateAmountBasedAuth) validation.RuleFunc {
	return func(_ interface{}) error {
		if u.Method == OPEN && u.MaxAmount == 0 {
			return fmt.Errorf("maximum amount is required for open method")
		}
		if u.Method == OTPANDPIN && u.MinAmount == 0 {
			return fmt.Errorf("minimum amount is required for OTP and PIN method")
		}
		if u.Method == PIN && (u.MaxAmount == 0 && u.MinAmount == 0) {
			return fmt.Errorf("either minimum or maximum amount are required for PIN method")
		}
		return nil
	}
}
