package hq

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	domain_hq "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/hq"
	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	sharedutils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationAbstracts interface {
	GetHQ(ctx context.Context, id string) (dto.HQ, error)
	GetHQDetail(ctx context.Context, filterParams *constant.Filter) (*utils.PaginatedResponse[[]*domain_hq.HQ], error)
	UpdateBlockTimeRequest(ctx context.Context, request dto.UpdateBlockTimeRequest, makerID, phone, fullName string) (string, error)
	UpdateArchiveTimeRequest(ctx context.Context, request dto.UpdateArchiveTimeRequest, makerID, phone, fullName string) (string, error)
	UpdateBlockTime(ctx context.Context, request dto.ApproveRejectRequest, checkerID, phone, fullName string) error
	UpdateArchiveTime(ctx context.Context, request dto.ApproveRejectRequest, checkerID, phone, fullName string) error
}

type ApplicationStore struct {
	service domain_hq.Service
	logger  sharedutils.Logger
}

func NewApplication(service domain_hq.Service, logger sharedutils.Logger) ApplicationAbstracts {
	return &ApplicationStore{service: service, logger: logger}
}

func (a *ApplicationStore) GetHQDetail(ctx context.Context, filterParams *constant.Filter) (*utils.PaginatedResponse[[]*domain_hq.HQ], error) {
	hqData, err := a.service.GetAllHQ(ctx, filterParams)
	if err != nil {
		return nil, err
	}
	return hqData, nil
}
func (a *ApplicationStore) GetHQ(ctx context.Context, id string) (dto.HQ, error) {
	hq, err := a.service.GetHQ(ctx, id)
	if err != nil {
		return dto.HQ{}, err
	}
	return dto.HQ{
		ID:             hq.ID.Hex(),
		Name:           hq.Name,
		BlockTime:      hq.BlockTime,
		ArchiveTime:    hq.ArchiveTime,
		CreatedAt:      hq.CreatedAt,
		LastModifiedAt: hq.LastModified,
	}, nil
}

func (a *ApplicationStore) UpdateBlockTimeRequest(ctx context.Context, request dto.UpdateBlockTimeRequest, makerID, phone, fullName string) (string, error) {
	return a.service.UpdateBlockTimeRequest(ctx, domain_hq.UpdateBlockTimeRequest{
		ID:         request.ID,
		BlockTime:  request.BlockTime,
		MakerID:    makerID,
		MakerName:  fullName,
		MakerPhone: phone,
	})
}

func (a *ApplicationStore) UpdateArchiveTimeRequest(ctx context.Context, request dto.UpdateArchiveTimeRequest, makerID, phone, fullName string) (string, error) {
	return a.service.UpdateArchiveTimeRequest(ctx, domain_hq.UpdateArchiveTimeRequest{
		ID:          request.ID,
		ArchiveTime: request.ArchiveTime,
		MakerID:     makerID,
		MakerName:   fullName,
		MakerPhone:  phone,
	})
}

func (a *ApplicationStore) UpdateBlockTime(ctx context.Context, request dto.ApproveRejectRequest, checkerID, phone, fullName string) error {
	return a.service.UpdateBlockTime(ctx, domain_hq.ApproveRejectRequest{
		ActionCode:   request.ActionCode,
		CheckerID:    checkerID,
		CheckerName:  fullName,
		CheckerPhone: phone,
		Decision:     request.Decison,
		RejectedReason: func() string {
			if request.Decison == utils.DecisionDenied {
				return request.RejectedReason
			}
			return ""
		}(),
	})
}

func (a *ApplicationStore) UpdateArchiveTime(ctx context.Context, request dto.ApproveRejectRequest, checkerID, phone, fullName string) error {
	return a.service.UpdateArchiveTime(ctx, domain_hq.ApproveRejectRequest{
		ActionCode:   request.ActionCode,
		CheckerID:    checkerID,
		CheckerName:  fullName,
		CheckerPhone: phone,
		Decision:     request.Decison,
		RejectedReason: func() string {
			if request.Decison == utils.DecisionDenied {
				return request.RejectedReason
			}
			return ""
		}(),
	})
}
