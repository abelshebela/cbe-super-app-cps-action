package passwordrule

import (
	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
)

type PasswordRuleUpdate struct {
	Rule struct {
		Name                   string `bson:"name" json:"name"`
		MinLength              int    `bson:"min_length" json:"min_length"`
		MaxLength              int    `bson:"max_length" json:"max_length"`
		Numbers                *bool  `bson:"numbers" json:"numbers"`
		CapitalLetters         *bool  `bson:"capital_letters" json:"capital_letters"`
		AllowSequentialNumbers *bool  `bson:"allow_sequential_numbers" json:"allow_sequential_numbers"`
		SmallLetters           *bool  `bson:"small_letters" json:"small_letters"`
		Characters             *bool  `bson:"characters" json:"characters"`
		IsSpacedAllowed        *bool  `json:"is_spaced_allowed" bson:"is_spaced_allowed"`
	} `json:"rule"`
}

type CheckPasswordDTO struct {
	Password string `json:"password"`
}

type PaginatedPasswordRulesResponse struct {
	Data []model.PasswordRule `json:"docs"`
	Meta types.PaginationMeta `json:"meta"`
}
