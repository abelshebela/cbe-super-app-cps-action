package dto

type ServiceResponse struct{}
type UpdateServiceRequest struct{}

type Tier struct {
	Min       float64 `json:"min" bson:"min" validate:"required,gte=0"`
	Max       float64 `json:"max" bson:"max" validate:"required,gtfield=Min"`
	FeeAmount float64 `json:"fee_amount" bson:"fee_amount" validate:"required,gte=0"`
}
