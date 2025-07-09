package hq

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	domain_hq "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/hq"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationAbstracts interface {
	GetHQ(ctx context.Context, id string) (dto.HQ, error)
	UpdateBlockTimeRequest(ctx context.Context, request dto.UpdateBlockTimeRequest, makerID, phone, fullName string) (string, error)
	UpdateArchiveTimeRequest(ctx context.Context, request dto.UpdateArchiveTimeRequest, makerID, phone, fullName string) (string, error)
	UpdateBlockTime(ctx context.Context, request dto.ApproveRejectRequest, checkerID, phone, fullName string) error
	UpdateArchiveTime(ctx context.Context, request dto.ApproveRejectRequest, checkerID, phone, fullName string) error
}

type ApplicationStore struct {
	service domain_hq.Service
	logger  utils.Logger
}

func NewApplication(service domain_hq.Service, logger utils.Logger) ApplicationAbstracts {
	return &ApplicationStore{service: service, logger: logger}
}

func (a *ApplicationStore) GetHQ(ctx context.Context, id string) (dto.HQ, error) {
	hq, err := a.service.GetHQ(ctx, id)
	if err != nil {
		return dto.HQ{}, err
	}
	return dto.HQ{
		ID:             hq.ID,
		Name:           hq.Name,
		BlockTime:      hq.BlockTime,
		ArchiveTime:    hq.ArchiveTime,
		CreatedAt:      hq.CreatedAt,
		LastModifiedAt: hq.LastModified,
	}, nil
}

func (a *ApplicationStore) UpdateBlockTimeRequest(ctx context.Context, request dto.UpdateBlockTimeRequest, makerID, phone, fullName string) (string, error) {
	return a.service.UpdateBlockTimeRequest(ctx, domain_hq.UpdateBlockTimeRequest{
		BlockTime:  request.BlockTime,
		MakerID:    makerID,
		MakerName:  fullName,
		MakerPhone: phone,
	})
}

func (a *ApplicationStore) UpdateArchiveTimeRequest(ctx context.Context, request dto.UpdateArchiveTimeRequest, makerID, phone, fullName string) (string, error) {
	return a.service.UpdateArchiveTimeRequest(ctx, domain_hq.UpdateArchiveTimeRequest{
		ArchiveTime: request.ArchiveTime,
		MakerID:     makerID,
		MakerName:   fullName,
		MakerPhone:  phone,
	})
}

func (a *ApplicationStore) UpdateBlockTime(ctx context.Context, request dto.ApproveRejectRequest, checkerID, phone, fullName string) error {
	return a.service.UpdateBlockTime(ctx, domain_hq.ApproveRejectRequest{
		ActionCode:   request.ActionID,
		CheckerID:    checkerID,
		CheckerName:  fullName,
		CheckerPhone: phone,
		Approved:     request.Approve,
	})
}

func (a *ApplicationStore) UpdateArchiveTime(ctx context.Context, request dto.ApproveRejectRequest, checkerID, phone, fullName string) error {
	return a.service.UpdateArchiveTime(ctx, domain_hq.ApproveRejectRequest{
		ActionCode:   request.ActionID,
		CheckerID:    checkerID,
		CheckerName:  fullName,
		CheckerPhone: phone,
		Approved:     request.Approve,
	})
}
