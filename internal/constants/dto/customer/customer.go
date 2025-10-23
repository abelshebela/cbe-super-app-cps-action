package customer

type CustomerEnableDTO struct {
	UserOTP string `json:"user_otp"`
}
type CustomerDisableDTO struct {
	IsTemporary   *bool  `json:"is_temporary"`
	DisableReason string `json:"disable_reason"`
}

type CustomerEnableSessionResponse struct {
	Otp string `json:"otp"`
}
