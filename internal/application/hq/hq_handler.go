package hq

import (
	"context"

	// cpsaction "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/inbound/http/cps_action"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"
	cps_entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	domain_hq "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/hq"
	utils "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"

	sharedutils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationAbstracts interface {
	GetHQ(ctx context.Context, id string) (dto.HQ, error)
	GetHQDetail(ctx context.Context, filterParams *constant.Filter) (*utils.PaginatedResponse[[]*domain_hq.HQ], error)
	UpdateBlockTimeRequest(ctx context.Context, request dto.UpdateBlockTimeRequest, makerID, phone, fullName, dept string) (*cps_entities.CPSAction, error)
	UpdateArchiveTimeRequest(ctx context.Context, request dto.UpdateArchiveTimeRequest, makerID, phone, fullName, dept string) (*cps_entities.CPSAction, error)
	GetBlockTime(ctx context.Context) (dto.BlockTimeResponse, error)
	GetArchiveTime(ctx context.Context) (dto.ArchiveTimeResponse, error)
	GetPasswordExpiry(ctx context.Context) (dto.PasswordExpiryResponse, error)
	UpdatePasswordExpiryRequest(ctx context.Context, request dto.UpdatePasswordExpiryRequest, makerID, phone, fullName, dept string) (*cps_entities.CPSAction, error)
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
		ID:          hq.ID.Hex(),
		BlockTime:   hq.BlockTime,
		ArchiveTime: hq.ArchiveTime,
	}, nil
}

func (a *ApplicationStore) GetBlockTime(ctx context.Context) (dto.BlockTimeResponse, error) {
	resp, err := a.service.GetBlockTime(ctx)
	if err != nil {
		return dto.BlockTimeResponse{}, err
	}
	return dto.BlockTimeResponse{
		BlockTime:      resp.BlockTime,
		CreatedAtBlock: resp.CreatedAtBlock,
		UpdatedAtBlock: resp.UpdatedAtBlock,
	}, nil
}

func (a *ApplicationStore) GetArchiveTime(ctx context.Context) (dto.ArchiveTimeResponse, error) {
	resp, err := a.service.GetArchiveTime(ctx)
	if err != nil {
		return dto.ArchiveTimeResponse{}, err
	}
	return dto.ArchiveTimeResponse{
		ArchiveTime:      resp.ArchiveTime,
		CreatedAtArchive: resp.CreatedAtArchive,
		UpdatedAtArchive: resp.UpdatedAtArchive,
	}, nil
}

func (a *ApplicationStore) GetPasswordExpiry(ctx context.Context) (dto.PasswordExpiryResponse, error) {
	resp, err := a.service.GetPasswordExpiry(ctx)
	if err != nil {
		return dto.PasswordExpiryResponse{}, err
	}
	return dto.PasswordExpiryResponse{
		PasswordExpiry:          resp.PasswordExpiry,
		CreatedAtPasswordExpiry: resp.CreatedAtPasswordExpiry,
		UpdatedAtPasswordExpiry: resp.UpdatedAtPasswordExpiry,
	}, nil
}

func (a *ApplicationStore) UpdateBlockTimeRequest(ctx context.Context, request dto.UpdateBlockTimeRequest, makerID, phone, fullName, dept string) (*cps_entities.CPSAction, error) {
	return a.service.UpdateBlockTimeRequest(ctx, domain_hq.UpdateBlockTimeRequest{
		BlockTime:  request.BlockTime,
		MakerID:    makerID,
		MakerName:  fullName,
		MakerPhone: phone,
		Department: dept,
	})
}

func (a *ApplicationStore) UpdateArchiveTimeRequest(ctx context.Context, request dto.UpdateArchiveTimeRequest, makerID, phone, fullName, dept string) (*cps_entities.CPSAction, error) {
	return a.service.UpdateArchiveTimeRequest(ctx, domain_hq.UpdateArchiveTimeRequest{
		ArchiveTime: request.ArchiveTime,
		MakerID:     makerID,
		MakerName:   fullName,
		MakerPhone:  phone,
		Department:  dept,
	})
}

func (a *ApplicationStore) UpdatePasswordExpiryRequest(ctx context.Context, request dto.UpdatePasswordExpiryRequest, makerID, phone, fullName, dept string) (*cps_entities.CPSAction, error) {
	return a.service.UpdatePasswordExpiryRequest(ctx, domain_hq.UpdatePasswordExpiryRequest{
		PasswordExpiry: request.PasswordExpiry,
		MakerID:        makerID,
		MakerName:      fullName,
		MakerPhone:     phone,
		Department:     dept,
	})
}
