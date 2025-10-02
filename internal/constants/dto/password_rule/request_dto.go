package passwordrule

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
