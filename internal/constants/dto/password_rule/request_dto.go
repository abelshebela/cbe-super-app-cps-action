package passwordrule

import (
	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

type PasswordRuleUpdate struct {
	Rule struct {
		Name           string `json:"name"`
		MinLength      int    `json:"min_length"`
		MaxLength      int    `json:"max_length"`
		Numbers        bool   `bson:"numbers" json:"numbers"`
		CapitalLetters bool   `json:"capital_letters"`
		SmallLetters   bool   `json:"small_letters"`
		Characters     bool   `json:"characters"`
	} `json:"rule"`
}

type CheckPasswordDTO struct {
	Password string `json:"password"`
}

type PaginatedPasswordRulesResponse struct {
	Data []model.PasswordRule `json:"docs"`
	Meta types.PaginationMeta `json:"meta"`
}
