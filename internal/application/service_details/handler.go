package service_details_app

import (
	"context"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/application/dto"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/service"
)

type ApplicationAbstracts interface {
	GetAllServiceDetails(ctx context.Context) ([]*dto.ServiceDetailsResponse, error)
	GetServiceDetailsByID(ctx context.Context, id string) (*dto.ServiceDetailsResponse, error)
	UpdateServiceDetailsRequest(ctx context.Context, id string, update *dto.UpdateServiceDetailsRequest) (*dto.UpdateServiceDetailsResponse, error)
	UpdateServiceDetails(ctx context.Context, req *dto.ApproveServiceDetailsRequest) error

	InitiateServiceFeeUpdate(ctx context.Context, cpsAction *dto.CPSAction) (*dto.UpdateServiceDetailsResponse, error)
	ApproveServiceFeeUpdate(ctx context.Context, cpsAction *dto.CPSAction) (*dto.UpdateServiceDetailsResponse, error)
	RejectServiceFeeUpdate(ctx context.Context, cpsAction *dto.CPSAction) (*dto.UpdateServiceDetailsResponse, error)
	UpdateServiceCap(ctx context.Context, id string, cap *service.Cap) (*dto.ServiceDetailsResponse, error) // <-- FIXED
	ApproveServiceDetails(ctx context.Context, req *dto.ApproveServiceDetailsRequest) error                 // <-- Add this

}

type ApplicationStore struct {
	service service.ServiceInterface
	logger  utils.Logger
}

func NewApplication(service service.ServiceInterface, logger utils.Logger) ApplicationAbstracts {
	return &ApplicationStore{
		service: service,
		logger:  logger,
	}
}

func (a *ApplicationStore) GetAllServiceDetails(ctx context.Context) ([]*dto.ServiceDetailsResponse, error) {
	services, err := a.service.GetAllServiceDetails(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.ServiceDetailsResponse, len(services))
	for i, svc := range services {
		responses[i] = &dto.ServiceDetailsResponse{
			ID:                 svc.ID,
			ServiceID:          svc.ServiceCode,
			ServiceCode:        svc.ServiceCode,
			ServiceName:        svc.ServiceName,
			ServiceType:        svc.ServiceType,
			Key:                svc.Key,
			Cap:                svc.Cap,
			CBEProductCodes:    svc.CBEProductCodes,
			CBEIFBProductCodes: svc.CBEIFBProductCodes,
			AboveAmount:        svc.AboveAmount,
			AboveServiceFee:    svc.AboveServiceFee,
			PaymentType:        svc.PaymentType,
			Tiers:              svc.Tiers,
			CBEGLEntry:         svc.CBEGLEntry,
			CBEIFBGLEntry:      svc.CBEIFBGLEntry,
			Enabled:            svc.Enabled,
			IsDeleted:          svc.IsDeleted,
			CreatedAt:          svc.CreatedAt,
			LastModifiedAt:     svc.LastModifiedAt,
			DeletedAt:          svc.DeletedAt,
		}
	}
	return responses, nil
}

func (a *ApplicationStore) GetServiceDetailsByID(ctx context.Context, id string) (*dto.ServiceDetailsResponse, error) {
	service, err := a.service.GetServiceDetailsByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.ServiceDetailsResponse{
		ID:                 service.ID,
		ServiceID:          service.ServiceCode,
		ServiceCode:        service.ServiceCode,
		ServiceName:        service.ServiceName,
		ServiceType:        service.ServiceType,
		Key:                service.Key,
		Cap:                service.Cap,
		CBEProductCodes:    service.CBEProductCodes,
		CBEIFBProductCodes: service.CBEIFBProductCodes,
		AboveAmount:        service.AboveAmount,
		AboveServiceFee:    service.AboveServiceFee,
		PaymentType:        service.PaymentType,
		Tiers:              service.Tiers,
		CBEGLEntry:         service.CBEGLEntry,
		CBEIFBGLEntry:      service.CBEIFBGLEntry,
		Enabled:            service.Enabled,
		IsDeleted:          service.IsDeleted,
		CreatedAt:          service.CreatedAt,
		LastModifiedAt:     service.LastModifiedAt,
		DeletedAt:          service.DeletedAt,
	}, nil
}

func (a *ApplicationStore) UpdateServiceDetailsRequest(ctx context.Context, id string, update *dto.UpdateServiceDetailsRequest) (*dto.UpdateServiceDetailsResponse, error) {
	// Convert DTO to domain entity
	serviceUpdate := &service.Service{
		ID:                 update.ID,
		ServiceCode:        update.ServiceCode,
		ServiceName:        update.ServiceName,
		ServiceType:        update.ServiceType,
		Key:                update.Key,
		Cap:                update.Cap,
		CBEProductCodes:    update.CBEProductCodes,
		CBEIFBProductCodes: update.CBEIFBProductCodes,
		AboveAmount:        update.AboveAmount,
		AboveServiceFee:    update.AboveServiceFee,
		PaymentType:        update.PaymentType,
		Tiers:              update.Tiers,
		CBEGLEntry:         update.CBEGLEntry,
		CBEIFBGLEntry:      update.CBEIFBGLEntry,
		Enabled:            update.Enabled,
	}

	actionID, err := a.service.UpdateServiceDetailsRequest(ctx, id, serviceUpdate, update.MakerID)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateServiceDetailsResponse{
		ActionID: actionID,
	}, nil
}

func (a *ApplicationStore) UpdateServiceDetails(ctx context.Context, req *dto.ApproveServiceDetailsRequest) error {
	return a.service.UpdateServiceDetails(ctx, req.ActionID, req.Approve, req.CheckerID, req.RejectionReason)
}

func (a *ApplicationStore) InitiateServiceFeeUpdate(ctx context.Context, req *dto.CPSAction) (*dto.UpdateServiceDetailsResponse, error) {
	//validation to be implimented

	cpsReq := service.CPSAction{
		ID:               req.ID,
		ActionCode:       req.ActionCode,
		MakerID:          req.MakerUser.UserCode,
		MakerName:        req.MakerUser.FullName,
		MakerPhoneNumber: req.MakerUser.PhoneNumber,
		Department:       req.Department,
		CurrentAction:    req.ActionData,
	}
	cpsReq.ActionCode = utils.RandomGenerator(20)

	res, err := a.service.InitiateServiceFeeUpdate(ctx, cpsReq)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateServiceDetailsResponse{
		ActionID: res.ActionID,
	}, nil
}

func (a *ApplicationStore) ApproveServiceFeeUpdate(ctx context.Context, req *dto.CPSAction) (*dto.UpdateServiceDetailsResponse, error) {
	//validation going to impliment
	cpsReq := service.CPSAction{
		ID:               req.ID,
		ActionCode:       req.ActionCode,
		MakerID:          req.MakerUser.UserCode,
		MakerName:        req.MakerUser.FullName,
		MakerPhoneNumber: req.MakerUser.PhoneNumber,
		Department:       req.Department,
		CurrentAction:    req.ActionData,
	}
	cpsReq.ActionCode = utils.RandomGenerator(20)

	res, err := a.service.InitiateServiceFeeUpdate(ctx, cpsReq)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateServiceDetailsResponse{
		ActionID: res.ActionID,
	}, nil
}

func (a *ApplicationStore) RejectServiceFeeUpdate(ctx context.Context, req *dto.CPSAction) (*dto.UpdateServiceDetailsResponse, error) {
	//validation going to impliment
	cpsReq := service.CPSAction{
		ID:               req.ID,
		ActionCode:       req.ActionCode,
		MakerID:          req.MakerUser.UserCode,
		MakerName:        req.MakerUser.FullName,
		MakerPhoneNumber: req.MakerUser.PhoneNumber,
		RejectionReason:  req.RejectedReason,
		Department:       req.Department,
		CurrentAction:    req.ActionData,
	}
	cpsReq.ActionCode = utils.RandomGenerator(20)

	res, err := a.service.InitiateServiceFeeUpdate(ctx, cpsReq)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateServiceDetailsResponse{
		ActionID: res.ActionID,
	}, nil
}
func (a *ApplicationStore) UpdateServiceCap(ctx context.Context, id string, cap *service.Cap) (*dto.ServiceDetailsResponse, error) {
	svc, err := a.service.UpdateServiceCap(ctx, id, *cap)
	if err != nil {
		return nil, err
	}
	return &dto.ServiceDetailsResponse{
		ID:          svc.ID,
		ServiceID:   svc.ServiceCode,
		ServiceCode: svc.ServiceCode,
		ServiceName: svc.ServiceName,
		ServiceType: svc.ServiceType,
		Key:         svc.Key,
		Cap:         svc.Cap,
	}, nil
}
func (a *ApplicationStore) ApproveServiceDetails(ctx context.Context, req *dto.ApproveServiceDetailsRequest) error {
	return a.service.UpdateServiceDetails(ctx, req.ActionID, req.Approve, req.CheckerID, req.RejectionReason)
}
