package transaction_dto

import (
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"time"
)

// TransactionLimitResponse is the API response shape for a transaction limit document.
type TransactionLimitResponse struct {
	ID             string                          `json:"id"`
	CustomerNumber string                          `json:"customer_number"`
	ChannelLimits  []imodel.TransactionLimitChannel `json:"channel_limits"`
	UpdatedAt      time.Time                       `json:"updated_at"`
}
