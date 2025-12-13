package cpsaction

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"

	actionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	"cbe-super-app-cps-action/internal/storage"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type cpsActionService struct {
	repo       storage.CPSActionRepository
	logger     utils.Logger
	dispatcher Dispatcher
}

func NewCPSActionService(repo storage.CPSActionRepository, logger utils.Logger, dispatcher Dispatcher) service.CPSActionService {
	return &cpsActionService{
		repo:       repo,
		logger:     logger,
		dispatcher: dispatcher,
	}
}

func (ca *cpsActionService) CreateCPSAction(ctx context.Context, cpsAction *model.CPSAction) error {
	existing, err := ca.GetCPSActionByUniqueID(ctx, cpsAction.RequestAction, cpsAction.Department)

	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		return err
	}

	if existing != nil {
		return errors.New(localization.ErrorPendingCpsActionExists.Code)
	}
	return ca.repo.Save(ctx, cpsAction)
}

func (ca *cpsActionService) ApproveCPSAction(ctx context.Context, action *model.CPSAction) error {
	data, err := ca.repo.Update(ctx, action.ActionCode, *action)
	if err != nil || data == nil {
		return err
	}
	approve, err := ca.dispatcher.Authorize(ctx, data)
	if err != nil && approve == nil {
		RollErr := ca.RollBack(ctx, action)
		if RollErr != nil {
			return RollErr
		}
		return err
	}
	return nil

}
func (ca *cpsActionService) RejectCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error {

	_, err := ca.repo.Update(ctx, action_code, *action)
	if err != nil {
		return err
	}
	return nil
}
func (ca *cpsActionService) GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	return ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, department)
}
func (ca *cpsActionService) GetCPSActionByID(ctx context.Context, id, department string) (*model.CPSAction, error) {
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		ca.logger.Errorf("their is error when try to parse the string to bson object in service")
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return ca.repo.SanitizedFindOne(ctx, bson.M{"_id": objID, "department": department})
}
func (ca *cpsActionService) GetCPSActionByUniqueID(ctx context.Context, requestAction, department string) (*model.CPSAction, error) {
	filter := bson.M{
		"department":     department,
		"action_status":  string(constants.Pending),
		"request_action": requestAction,
	}

	return ca.repo.SanitizedFindOne(context.Background(), filter)
}
func (ca *cpsActionService) GetCPSActionByActionCode(ctx context.Context, uniqueID, department string) (*model.CPSAction, error) {
	return ca.repo.SanitizedFindOne(ctx, bson.M{"action_code": uniqueID})
}

func (ca *cpsActionService) RollBack(ctx context.Context, action *model.CPSAction) error {
	return ca.repo.UpdateCustome(ctx, bson.M{"action_code": action.ActionCode}, bson.M{"action_status": string(constants.Pending), "checker_id": "", "checker_name": "", "checker_phone_number": ""})
}

func (ca *cpsActionService) GetActionCountsByDepartemnt(ctx context.Context, department string) (*actionDto.CPSActionCountResponse, error) {
	return ca.repo.GetCountByDepartment(ctx, department)
}
