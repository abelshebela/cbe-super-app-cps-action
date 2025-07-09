package unlink

import "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink/entities"

type Repository interface {
	UnlinkDevice(userCode string, cpsAction entities.CPSAction) error
	ApproveOrDecline(userCode, decision, reason string, cpsAction entities.CPSAction) error
}
