package faydaaccount

import (
	"context"
	"encoding/json"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/service"

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
	// Unmarshal CurrentAction to ActionData
	var actionData ActionData
	jsonBytes, err := json.Marshal(req.CurrentAction)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(jsonBytes, &actionData); err != nil {
		return nil, err
	}

	cpsReq := entity.CPSAction{
		ID:               req.ID,
		ActionCode:       req.ActionCode,
		MakerID:          req.MakerID,
		MakerName:        req.MakerName,
		MakerPhoneNumber: req.MakerPhoneNumber,
		Department:       req.Department,
		CurrentAction: entity.ActionData{
			UseCode:     actionData.UseCode,
			FullName:    actionData.FullName,
			PhoneNumber: actionData.PhoneNumber,
		},
	}
	cpsReq.ActionCode = utils.RandomGenerator(20)

	res, err := f.FaydaDomain.InitiateDisableFaydaAccount(ctx, cpsReq)
	if err != nil {
		return nil, err
	}

	// Unmarshal result CurrentAction to ActionData
	var resultActionData ActionData
	resultBytes, err := json.Marshal(res.CurrentAction)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resultBytes, &resultActionData); err != nil {
		return nil, err
	}

	return &CPSAction{
		ID:               res.ID,
		ActionCode:       res.ActionCode,
		MakerID:          res.MakerID,
		MakerName:        res.MakerName,
		MakerPhoneNumber: res.MakerPhoneNumber,
		ActionType:       ActionType(res.ActionType),
		CurrentAction: ActionData{
			UseCode:     resultActionData.UseCode,
			FullName:    resultActionData.FullName,
			PhoneNumber: resultActionData.PhoneNumber,
		},
		ActionStatus:    ActionStatus(res.ActionStatus),
		MakerActionTime: res.MakerActionTime,
		Department:      res.Department,
		RequestAction:   RequestAction(res.RequestAction),
	}, nil
}

func (f FaydaHandler) AuthorizeFaydaAccountDisable(ctx context.Context, req CPSAction) (*CPSAction, error) {
	var actionData ActionData
	jsonBytes, err := json.Marshal(req.CurrentAction)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(jsonBytes, &actionData); err != nil {
		return nil, err
	}

	cpsReq := entity.CPSAction{
		ID:                 req.ID,
		ActionCode:         req.ActionCode,
		CheckerID:          req.CheckerID,
		CheckerName:        req.CheckerName,
		CheckerPhoneNumber: req.CheckerPhoneNumber,
		Department:         req.Department,
	}

	res, err := f.FaydaDomain.AuthorizeFaydaAccountDisable(ctx, cpsReq)
	if err != nil {
		return nil, err
	}

	var resultActionData ActionData
	resultBytes, err := json.Marshal(res.CurrentAction)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resultBytes, &resultActionData); err != nil {
		return nil, err
	}

	return &CPSAction{
		ID:                 res.ID,
		ActionCode:         res.ActionCode,
		CheckerID:          res.CheckerID,
		CheckerName:        res.CheckerName,
		CheckerPhoneNumber: res.CheckerPhoneNumber,
		Department:         res.Department,
		ActionStatus:       ActionStatus(res.ActionStatus),
		ActionType:         ActionType(res.ActionType),
		MakerID:            res.MakerID,
		MakerName:          res.MakerName,
		MakerPhoneNumber:   res.MakerPhoneNumber,
		CurrentAction: ActionData{
			UseCode:     resultActionData.UseCode,
			FullName:    resultActionData.FullName,
			PhoneNumber: resultActionData.PhoneNumber,
		},
		RequestAction:     RequestAction(res.RequestAction),
		PreviosAction:     res.PreviosAction,
		MakerActionTime:   res.MakerActionTime,
		CheckerActionTime: res.CheckerActionTime,
	}, nil
}

func (f FaydaHandler) RejectFaydaAccountDisable(ctx context.Context, req CPSAction) (*CPSAction, error) {

	cpsReq := entity.CPSAction{
		ID:                 req.ID,
		ActionCode:         req.ActionCode,
		CheckerID:          req.CheckerID,
		CheckerName:        req.CheckerName,
		CheckerPhoneNumber: req.CheckerPhoneNumber,
		RejectionReason:    req.RejectionReason,

		Department: req.Department,
	}

	res, err := f.FaydaDomain.RejectFaydaAccountDisable(ctx, cpsReq)
	if err != nil {
		return nil, err
	}

	var resultActionData ActionData
	resultBytes, err := json.Marshal(res.CurrentAction)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(resultBytes, &resultActionData); err != nil {
		return nil, err
	}

	return &CPSAction{
		ID:                 res.ID,
		ActionCode:         res.ActionCode,
		CheckerID:          res.CheckerID,
		CheckerName:        res.CheckerName,
		CheckerPhoneNumber: res.CheckerPhoneNumber,
		RejectionReason:    res.RejectionReason,
		Department:         res.Department,
		ActionStatus:       ActionStatus(res.ActionStatus),
		ActionType:         ActionType(res.ActionType),
		MakerID:            res.MakerID,
		MakerName:          res.MakerName,
		MakerPhoneNumber:   res.MakerPhoneNumber,
		CurrentAction: ActionData{
			UseCode:     resultActionData.UseCode,
			FullName:    resultActionData.FullName,
			PhoneNumber: resultActionData.PhoneNumber,
		},
		RequestAction:     RequestAction(res.RequestAction),
		MakerActionTime:   res.MakerActionTime,
		CheckerActionTime: res.CheckerActionTime,
	}, nil
}
