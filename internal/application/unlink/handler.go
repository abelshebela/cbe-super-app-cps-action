// Package unlink implements the use case or business logic for the unlink functionality.
package unlink

import (
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/unlink/entities"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	UnlinkDevice(userCode string, cpsAction entities.CPSAction) error
	ApproveOrDecline(userCode, decision, reason string, cpsAction entities.CPSAction) error
}

type UnlinkHandler struct {
	service *domain.Service
}

func NewUnlinkHandler(service *domain.Service) ApplicationService {
	return &UnlinkHandler{service: service}
}

func (h *UnlinkHandler) UnlinkDevice(userCode string, cpsAction entities.CPSAction) error {
	cpsAction.ActionCode = utils.RandomGenerator(20)
	err := h.service.UnlinkDevice(userCode, cpsAction)

	if err != nil {
		return err
	}

	return nil
}

func (h *UnlinkHandler) ApproveOrDecline(userCode, decision, reason string, cpsAction entities.CPSAction) error {
	err := h.service.ApproveOrDecline(userCode, decision, reason, cpsAction)
	if err != nil {
		return err
	}

	return nil
}
