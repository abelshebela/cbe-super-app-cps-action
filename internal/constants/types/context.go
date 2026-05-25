package types

import (
	"cbe-super-app-cps-action/internal/constants"
	"context"
)

type B struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Merchant struct {
	ID       string `json:"id"`
	Branches []B    `json:"branches"`
}

// ContextMetadata is a pointer-based metadata object that can be passed through context
// to allow the service layer to communicate information back to the handler.
type ContextMetadata struct {
	IsMakerOnly   bool
	CPSActionCode string
	Merchant      *Merchant
}

// GetMetadata retrieves the ContextMetadata from the context.
func GetMetadata(ctx context.Context) *ContextMetadata {
	if md, ok := ctx.Value(constants.ContextKeyMetadata).(*ContextMetadata); ok {
		return md
	}
	return nil
}

// SetIsMakerOnly sets the IsMakerOnly field in the ContextMetadata if present in the context.
func SetIsMakerOnly(ctx context.Context, value bool) {
	if md := GetMetadata(ctx); md != nil {
		md.IsMakerOnly = value
	}
}

func SetCPSActionCode(ctx context.Context, value string) {
	if md := GetMetadata(ctx); md != nil {
		md.CPSActionCode = value
	}
}

func SetMerchant(ctx context.Context, merchant *Merchant) {
	if md := GetMetadata(ctx); md != nil {
		md.Merchant = merchant
	}
}
