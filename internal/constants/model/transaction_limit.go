package model

import "time"

type TransactionLimitEntry struct {
	ServiceName          string  `json:"service_name" bson:"service_name"`
	MaximumAmount        float64 `json:"maximum_amount" bson:"maximum_amount"`
	TransactionFrequency int     `json:"transaction_frequency" bson:"transaction_frequency"`
}

type TransactionLimitChannel struct {
	Channel  string                  `json:"channel" bson:"channel"`
	Services []TransactionLimitEntry `json:"services" bson:"services"`
}

type TransactionLimit struct {
	ID             string                    `json:"id" bson:"_id"`
	CustomerNumber string                    `json:"customer_number" bson:"customer_number"`
	ChannelLimits  []TransactionLimitChannel `json:"channel_limits" bson:"channel_limits"`
	UpdatedAt      time.Time                 `json:"updated_at" bson:"updated_at"`
}
