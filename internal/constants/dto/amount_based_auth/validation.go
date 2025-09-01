package amount_based_auth

func (r UpdateOpenTierRequest) Validate() bool {
	return r.MaxAmount > 0
}

func (r UpdatePinTierRequest) Validate() bool {
	return r.MinAmount > 0 && r.MaxAmount > 0
}

func (r UpdateOtpPinTierRequest) Validate() bool {
	return r.MinAmount > 0
}
