package services

type CapRequest struct {
	SingleCap          int64 `json:"single_cap" example:"1000000"`
	MinimumTransferCap int64 `json:"minimum_transfer_cap" example:"10"`
}

type TierRequest struct {
	FeeType   string `json:"fee_type" example:"PERCENT"`
	FeeAmount uint8  `json:"fee_amount" example:"1"`
	Min       int64  `json:"min" example:"0"`
	Max       int64  `json:"max" example:"5000"`
}

type CreateServiceRequest struct {
	ServiceName      string        `json:"service_name" example:"Transfer To Other bank"`
	ServiceCode      string        `json:"service_code" example:"JKSJBDJB"`
	ProductAccount   string        `json:"cbe_gl_product_account" example:"234353354"`
	ChargeCode       string        `json:"charge_code" example:"JKSJBDJB"`
	CommissionCode   string        `json:"commission_code" example:"DSFDDHJS"`
	Cap              CapRequest    `json:"cap"`
	Tiers            []TierRequest `json:"tiers"`
	AboveAmount      int64         `json:"above_amount" example:"10000"`
	AboveServiceFee  int64         `json:"above_service_fee" example:"500"`
	AbovePaymentType string        `json:"above_payment_type"`
}

type UpdateServiceRequest struct {
	ServiceName      string        `json:"service_name" example:"Transfer To Other bank"`
	ServiceCode      string        `json:"service_code" example:"JKSJBDJB"`
	ProductAccount   string        `json:"cbe_gl_product_account" example:"234353354"`
	ChargeCode       string        `json:"charge_code" example:"JKSJBDJB"`
	CommissionCode   string        `json:"commission_code" example:"DSFDDHJS"`
	Cap              CapRequest    `json:"cap"`
	Tiers            []TierRequest `json:"tiers"`
	AboveAmount      int64         `json:"above_amount" example:"10000"`
	AboveServiceFee  int64         `json:"above_service_fee" example:"500"`
	AbovePaymentType string        `json:"above_payment_type"`
}
