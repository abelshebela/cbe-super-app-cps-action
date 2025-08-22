package bpsuser

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type bpsUserService struct {
	cpsService service.CPSActionService
	repo   storage.BPSUserRepository
	logger utils.Logger
}

func NewBPSUserService(repo storage.BPSUserRepository, cpsService service.CPSActionService, logger utils.Logger) service.BPSUserService {
	return &bpsUserService{
		cpsService: cpsService,
		repo:   repo,
		logger: logger,
	}
}

// Authorize implements service.BPSUserService.
func (b *bpsUserService) Authorize(ctx context.Context, cpsAction *model.CPSAction)  error {

	updateData:=cpsAction.CurrentAction.(model.BPSUser)
	switch cpsAction.RequestAction {
	case string(constants.RequestEnableBPSUser):
		updateData.Enabled=true
	case string(constants.RequestDisableBPSUser):
		updateData.Enabled=false
	default:
		return errors.New(localization.ErrorActionNotFound.Code)
	}
	return b.repo.Update(ctx,&updateData)
}

// FetchUserByUserCode implements service.BPSUserService.
func (b *bpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*model.BPSUser, error) {
	return b.repo.GetByUserCode(ctx,userCode)
}

// GetAllBPSUsers implements service.BPSUserService.
func (b *bpsUserService) GetAllBPSUsers(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.BPSUser], error) {
	return b.repo.FindAllWithPagination(ctx,*filterParams)
}

// UpdateBpsUser implements service.BPSUserService.
func (b *bpsUserService) UpdateBpsUser(ctx context.Context, userCode string, status bool) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	user,err := b.repo.GetByUserCode(ctx,userCode)
	if err != nil {
		return err
	}

	if user.Enabled && status {
		return errors.New(localization.ErrorUserAlreadyEnabled.Code)
	}

	if !user.Enabled && !status {
		return errors.New(localization.ErrorUserAlreadyDisabled.Code)
	}

	updatedUser := *user
	updatedUser.Enabled = status

		var requestAction string
		if status{
			requestAction =string(constants.RequestEnableBPSUser)
		}else{
			requestAction =string(constants.RequestDisableBPSUser)

		}
		fmt.Println(requestAction)
		
		cpsActionData := lib.CpsModelBuilder(user.ID.Hex(),makerData,user,updatedUser,requestAction,constants.UPDATE)
		
		if err := b.cpsService.CreateCPSAction(ctx,&cpsActionData); err != nil {
			return err
		}
	return nil
}

