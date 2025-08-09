package unlink

import (
	"context"
)

// type Repository interface {
// 	UnlinkDevice(ctx context.Context, userCode string, cpsAction entities.CPSAction) (string, error)
// 	ApproveOrDecline(ctx context.Context, userCode, decision, reason string, cpsAction entities.CPSAction) error
// }

type UnlinkAccount interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*any, error)
	UnlinkUserCif(ctx context.Context, userCode string) error
	Authorize(ctx context.Context, cpsAction any) (any, error)
}
