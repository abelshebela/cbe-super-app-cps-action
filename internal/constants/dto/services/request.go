package services

type CapRequest struct {
	SingleCap          *float64 `json:"single_cap" example:"1000000"`
	MinimumTransferCap *float64 `json:"minimum_transfer_cap" example:"10"`
}

type TierRequest struct {
	FeeType   string   `json:"fee_type" example:"PERCENT"`
	FeeAmount *float64 `json:"fee_amount" example:"1"`
	Min       *float64 `json:"min" example:"0"`
	Max       *float64 `json:"max" example:"5000"`
}

type ServiceList struct {
	ServiceName             string        `bson:"service_name" json:"service_name"`
	ServiceKey              string        `bson:"service_key" json:"service_key"`
	OverideCap              CapRequest    `bson:"overide_cap" json:"overide_cap"`
	OverideProductGlAccount string        `bson:"overide_product_gl_account" json:"overide_product_gl_account"`
	OverideTiers            []TierRequest `bson:"overide_tiers" json:"overide_tiers"`
	HaveAnOverideTiers      bool          `bson:"have_an_overide_tiers" json:"have_an_overide_tiers"`
	IsEnabled               *bool         `bson:"is_enabled" json:"is_enabled"`
}

type CreateServiceRequest struct {
	ServiceName      string        `json:"service_name" example:"Transfer To Other bank"`
	ServiceKey       string        `json:"service_key" example:"Transfer To Other bank"`
	ServiceCode      string        `json:"service_code" example:"JKSJBDJB"`
	ProductGlAccount string        `json:"cbe_gl_product_account" example:"234353354"`
	ServiceList      []ServiceList `bson:"service_list" json:"service_list"`
	Cap              CapRequest    `json:"cap"`
	Tiers            []TierRequest `json:"tiers"`
	HaveATier        bool          `json:"have_a_tier"`
	HaveAChild       bool          `json:"have_a_child"`
	Enabled          *bool         `bson:"enabled" json:"enabled"`
	IsDeleted        *bool         `bson:"is_deleted" json:"is_deleted"`
}

type UpdateServiceRequest struct {
	ServiceName      string        `json:"service_name" example:"Transfer To Other bank"`
	ServiceKey       string        `json:"service_key" example:"Transfer To Other bank"`
	ServiceCode      string        `json:"service_code" example:"JKSJBDJB"`
	ProductGlAccount string        `json:"cbe_gl_product_account" example:"234353354"`
	ServiceList      []ServiceList `bson:"service_list" json:"service_list"`
	Cap              CapRequest    `json:"cap"`
	Tiers            []TierRequest `json:"tiers"`
	HaveATier        bool          `json:"have_a_tier"`
	HaveAChild       bool          `json:"have_a_child"`
	Enabled          *bool         `bson:"enabled" json:"enabled"`
	IsDeleted        *bool         `bson:"is_deleted" json:"is_deleted"`
}
