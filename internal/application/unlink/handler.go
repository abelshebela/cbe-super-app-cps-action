package unlink

import (
    domain "gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/unlink"
)

type ApplicationService interface {
    UnlinkDevice(userCode, makerUser string, branchCode []string, homeBranch string) error
    ApproveOrDecline(userCode, decision, reason, checkerUser string) error
}

type UnlinkHandler struct {
    service *domain.Service
}

func NewUnlinkHandler(service *domain.Service) ApplicationService {
    return &UnlinkHandler{service: service}
}

func (h *UnlinkHandler) UnlinkDevice(userCode, makerUser string, branchCode []string, homeBranch string) error {
    err := h.service.UnlinkDevice(userCode, makerUser, branchCode, homeBranch)
    
    if err != nil {
        return err
    }

    return nil
}


func (h *UnlinkHandler) ApproveOrDecline(userCode, decision, reason, checkerUser string) error {
    err := h.service.ApproveOrDecline(userCode, decision, reason, checkerUser)
    if err != nil {
        return err
    }

    return nil
}

