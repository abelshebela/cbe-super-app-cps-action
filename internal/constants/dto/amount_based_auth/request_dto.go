package amount_based_auth

// UpdateOpenTierRequest carries fields allowed for OPEN tier update
type UpdateOpenTierRequest struct {
	MaxAmount uint64 `json:"max_amount"`
}

// UpdatePinTierRequest carries fields allowed for PIN tier update
type UpdatePinTierRequest struct {
	MinAmount uint64 `json:"min_amount"`
	MaxAmount uint64 `json:"max_amount"`
}

// UpdateOtpPinTierRequest carries fields allowed for OTP_PIN tier update
type UpdateOtpPinTierRequest struct {
	MinAmount uint64 `json:"min_amount"`
}

