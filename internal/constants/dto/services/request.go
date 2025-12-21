package servicesdto

// CapRequest mirrors model.Cap with swagger examples for docs
// Use in request bodies only; service layer still uses model.Services.
type CapRequest struct {
	KYCLevel           string `json:"kyc_level" example:"Level1"`
	SingleCap          int64  `json:"single_cap" example:"1000000"`
	MinimumTransferCap int64  `json:"minimum_transfer_cap" example:"10"`
}

type TierRequest struct {
	FeeType   string `json:"fee_type" example:"PERCENT"`
	FeeAmount uint8  `json:"fee_amount" example:"1"`
	Min       int64  `json:"min" example:"0"`
	Max       int64  `json:"max" example:"5000"`
}

// CreateServiceRequest documents required fields for creating a Service
// Includes charge_code, commission_code, and cap.minimum_transfer_cap
// and uses cbe_gl_product_account for Product Account
type CreateServiceRequest struct {
	ServiceCode          string        `json:"service_code" example:"JKSJBDJB"`
	ServiceName          string        `json:"service_name" example:"Transfer To Other bank"`
	ProductAccount       string        `json:"product_account" example:"234353354"`
	ServiceType          string        `json:"service_type" example:"TRANSFER"`
	Key                  string        `json:"key" example:"transfer_to_other_bank"`
	ChargeCode           string        `json:"charge_code" example:"JKSJBDJB"`
	CommissionCode       string        `json:"commission_code" example:"DSFDDHJS"`
	Cap                  CapRequest    `json:"cap"`
	CbeGLProductAccount  string        `json:"cbe_gl_product_account" example:"234353354"`
	CbeIFBProductAccount string        `json:"cbe_ifb_product_account,omitempty" example:""`
	PaymentType          string        `json:"payment_type" example:"Tier Percentage"`
	Tiers                []TierRequest `json:"tiers"`
	AboveAmount          int64         `json:"above_amount" example:"10000"`
	AboveServiceFee      int64         `json:"above_service_fee" example:"500"`
	Enabled              bool          `json:"enabled" example:"true"`
}

// UpdateServiceRequest mirrors CreateServiceRequest for documentation
// Partial updates are allowed; examples remain the same.
type UpdateServiceRequest struct {
	ServiceCode          string        `json:"service_code" example:"JKSJBDJB"`
	ServiceName          string        `json:"service_name" example:"Transfer To Other bank"`
	ServiceType          string        `json:"service_type" example:"TRANSFER"`
	Key                  string        `json:"key" example:"transfer_to_other_bank"`
	ChargeCode           string        `json:"charge_code" example:"JKSJBDJB"`
	CommissionCode       string        `json:"commission_code" example:"DSFDDHJS"`
	Cap                  CapRequest    `json:"cap"`
	CbeGLProductAccount  string        `json:"cbe_gl_product_account" example:"234353354"`
	CbeIFBProductAccount string        `json:"cbe_ifb_product_account,omitempty" example:""`
	PaymentType          string        `json:"payment_type" example:"Tier Percentage"`
	Tiers                []TierRequest `json:"tiers"`
	AboveAmount          int64         `json:"above_amount" example:"10000"`
	AboveServiceFee      int64         `json:"above_service_fee" example:"500"`
	Enabled              bool          `json:"enabled" example:"true"`
}
