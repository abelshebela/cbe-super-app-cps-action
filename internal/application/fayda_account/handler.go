package faydaaccount

import (
	"context"
	"encoding/json"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/entity"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/fayda_account/service"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ApplicationService interface {
	InitiateDisableFaydaAccount(ctx context.Context, req entities.CPSAction) (*CPSAction, error)
	GetAllFaydaAccounts(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error)
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

func (f FaydaHandler) InitiateDisableFaydaAccount(ctx context.Context, req entities.CPSAction) (*CPSAction, error) {
	// Unmarshal CurrentAction to ActionData
	var actionData ActionData
	jsonBytes, err := json.Marshal(req.CurrentAction)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(jsonBytes, &actionData); err != nil {
		return nil, err
	}

	cpsReq := entities.CPSAction{
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



func (f FaydaHandler) GetAllFaydaAccounts(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*member.User], error) {
	return f.FaydaDomain.GetAllFaydaAccounts(ctx, filterParams)
}
