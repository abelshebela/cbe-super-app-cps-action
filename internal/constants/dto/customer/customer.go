package customer

type CustomerEnableDTO struct {
	UserOTP string `json:"user_otp"`
}

type CustomerEnableSessionResponse struct {
	Otp string `json:"otp"`
}
