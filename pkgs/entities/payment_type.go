package entities

import "fmt"

// PaymentType represents the type of payment for service fees
type PaymentType string

const (
	PaymentTypePercentage PaymentType = "percentage"
	PaymentTypeFlatFee    PaymentType = "flat_fee"
)

// IsValid checks if the PaymentType is valid
func (pt PaymentType) IsValid() bool {
	switch pt {
	case PaymentTypePercentage, PaymentTypeFlatFee:
		return true
	default:
		return false
	}
}

// String returns the string representation of PaymentType
func (pt PaymentType) String() string {
	return string(pt)
}

// ParsePaymentType parses a string to PaymentType
func ParsePaymentType(s string) (PaymentType, error) {
	pt := PaymentType(s)
	if !pt.IsValid() {
		return "", fmt.Errorf("invalid payment type: %s, must be either 'percentage' or 'flat_fee'", s)
	}
	return pt, nil
}
