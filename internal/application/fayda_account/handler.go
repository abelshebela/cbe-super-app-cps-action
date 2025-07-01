package faydaaccount

import (
	"context"

	"cbe-super-app-cps-action/internal/domain/fayda_account/entity"
	"cbe-super-app-cps-action/internal/domain/fayda_account/service"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	InitiateDisableFaydaAccount(ctx context.Context, req CPSAction) (*CPSAction, error)
	AuthorizeFaydaAccountDisable(ctx context.Context, req CPSAction) (*CPSAction, error)
	RejectFaydaAccountDisable(ctx context.Context, req CPSAction) (*CPSAction, error)
}

type FaydaHandler struct {
	FaydaDomain *service.FaydaAccountDomain
	logger      utils.Logger
}

func InitFaydaHandler(faydaDomain *service.FaydaAccountDomain, logger utils.Logger) ApplicationService {
	return FaydaHandler{
		FaydaDomain: faydaDomain,
		logger:      logger,
	}
}

func (f FaydaHandler) InitiateDisableFaydaAccount(ctx context.Context, req CPSAction) (*CPSAction, error) {
	if err := req.ActionData.Validate(); err != nil {
		f.logger.Errorf("validation errors", err)
		return nil, err
	}

	cpsReq := entity.CPSAction{
		ID:         req.ID,
		ActionCode: req.ActionCode,
		MakerUser: entity.User{
			UserCode:    req.MakerUser.UserCode,
			FullName:    req.MakerUser.FullName,
			PhoneNumber: req.MakerUser.PhoneNumber,
		},
		Department: req.Department,
		ActionData: entity.ActionData{
			UseCode:     req.ActionData.UseCode,
			FullName:    req.ActionData.FullName,
			PhoneNumber: req.ActionData.PhoneNumber,
		},
	}
	cpsReq.ActionCode = utils.RandomGenerator(20)

	res, err := f.FaydaDomain.InitiateDisableFaydaAccount(ctx, cpsReq)
	if err != nil {
		return nil, err
	}

	return &CPSAction{
		ID:         res.ID,
		ActionCode: res.ActionCode,
		MakerUser: User{
			UserCode:    res.MakerUser.UserCode,
			FullName:    res.MakerUser.FullName,
			PhoneNumber: res.MakerUser.PhoneNumber,
		},
		ActionType: ActionType(res.ActionType),
		ActionData: ActionData{
			UseCode:     res.ActionData.UseCode,
			FullName:    res.ActionData.FullName,
			PhoneNumber: res.ActionData.PhoneNumber,
		},
		Status:          ActionStatus(res.Status),
		MakerActionTime: res.MakerActionTime,
		Department:      res.Department,
		RequestAction:   RequestAction(res.RequestAction),
	}, nil
}

func (f FaydaHandler) AuthorizeFaydaAccountDisable(ctx context.Context, req CPSAction) (*CPSAction, error) {
	if err := req.Validate("APPROVE"); err != nil {
		f.logger.Errorf("validation error", err)
		return nil, err
	}

	cpsReq := entity.CPSAction{
		ID:         req.ID,
		ActionCode: req.ActionCode,
		CheckerUser: entity.User{
			UserCode:    req.CheckerUser.UserCode,
			FullName:    req.CheckerUser.FullName,
			PhoneNumber: req.CheckerUser.PhoneNumber,
		},
		Department: req.Department,
		ActionData: entity.ActionData{
			UseCode:     req.ActionData.UseCode,
			FullName:    req.ActionData.FullName,
			PhoneNumber: req.ActionData.PhoneNumber,
		},
	}

	res, err := f.FaydaDomain.AuthorizeFaydaAccountDisable(ctx, cpsReq)
	if err != nil {
		return nil, err
	}

	return &CPSAction{
		ID:         res.ID,
		ActionCode: res.ActionCode,
		CheckerUser: User{
			UserCode:    res.CheckerUser.UserCode,
			FullName:    res.CheckerUser.FullName,
			PhoneNumber: res.CheckerUser.PhoneNumber,
		},
		Department: res.Department,
		Status:     ActionStatus(res.Status),
		ActionType: ActionType(res.ActionType),
		MakerUser: User{
			UserCode:    res.MakerUser.UserCode,
			FullName:    res.MakerUser.FullName,
			PhoneNumber: res.MakerUser.PhoneNumber,
		},
		ActionData: ActionData{
			UseCode:     res.ActionData.UseCode,
			FullName:    res.ActionData.FullName,
			PhoneNumber: res.ActionData.PhoneNumber,
		},
		RequestAction:     RequestAction(res.RequestAction),
		PreviousData:      res.PreviousData,
		CurrentData:       res.CurrentData,
		MakerActionTime:   res.MakerActionTime,
		CheckerActionTime: res.CheckerActionTime,
	}, nil
}

func (f FaydaHandler) RejectFaydaAccountDisable(ctx context.Context, req CPSAction) (*CPSAction, error) {
	if err := req.Validate("REJECT"); err != nil {
		f.logger.Errorf("validation error", err)
		return nil, err
	}

	cpsReq := entity.CPSAction{
		ID:         req.ID,
		ActionCode: req.ActionCode,
		CheckerUser: entity.User{
			UserCode:    req.CheckerUser.UserCode,
			FullName:    req.CheckerUser.FullName,
			PhoneNumber: req.CheckerUser.PhoneNumber,
		},
		RejectedReason: req.RejectedReason,
		ActionData: entity.ActionData{
			UseCode:     req.ActionData.UseCode,
			FullName:    req.ActionData.FullName,
			PhoneNumber: req.ActionData.PhoneNumber,
		},
		Department: req.Department,
	}

	res, err := f.FaydaDomain.RejectFaydaAccountDisable(ctx, cpsReq)
	if err != nil {
		return nil, err
	}

	return &CPSAction{
		ID:         res.ID,
		ActionCode: res.ActionCode,
		CheckerUser: User{
			UserCode:    res.CheckerUser.UserCode,
			FullName:    res.CheckerUser.FullName,
			PhoneNumber: res.CheckerUser.PhoneNumber,
		},
		RejectedReason: res.RejectedReason,
		Department:     res.Department,
		Status:         ActionStatus(res.Status),
		ActionType:     ActionType(res.ActionType),
		MakerUser: User{
			UserCode:    res.MakerUser.UserCode,
			FullName:    res.MakerUser.FullName,
			PhoneNumber: res.MakerUser.PhoneNumber,
		},
		RequestAction:     RequestAction(res.RequestAction),
		MakerActionTime:   res.MakerActionTime,
		CheckerActionTime: res.CheckerActionTime,
	}, nil
}
