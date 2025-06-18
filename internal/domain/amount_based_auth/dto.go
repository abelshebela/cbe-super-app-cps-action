package amount_based_auth_domain

type AmountBasedAuthRequest struct {
	Id        string `json:"id"`
	MinAmount int    `json:"minAmount"`
	MaxAmount int    `json:"maxAmount"`
}
