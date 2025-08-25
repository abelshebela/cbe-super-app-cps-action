package cpsaction

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/factory"

	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance"
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

func NewCPSActionService(repo storage.CPSActionRepository, persistence persistance.Persistence, logger utils.Logger) service.CPSActionService {
	serviceFactory := factory.NewServiceFactory(persistence, logger)

	services := serviceFactory.CreateServiceContainer()

	dispatcherPersist := NewDispatcher(services)
	return &cpsActionService{
		repo:       repo,
		logger:     logger,
		dispatcher: *dispatcherPersist,
	}
}

func (ca *cpsActionService) CreateCPSAction(ctx context.Context, cpsAction *model.CPSAction) error {

	existing, err := ca.GetCPSActionByUniqueID(ctx, cpsAction.RequestAction, cpsAction.Department)
	if err != nil && err.Error() != localization.ErrorActionNotFound.Code {
		return err
	}
	if existing != nil {
		return errors.New(localization.ErrorPendingCpsActionExists.Code)
	}

	return ca.repo.Save(ctx, cpsAction)
}

func (ca *cpsActionService) ApproveCPSAction(ctx context.Context, action *model.CPSAction) error {

	if err := ca.repo.Save(ctx, action); err != nil {
		return err
	}

	approve, err := ca.dispatcher.Authorize(ctx, action)
	if err != nil && approve == nil {
		ca.RollBack(ctx, action.ActionCode)
		return err
	}

	return nil
}
func (ca *cpsActionService) RejectCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error {
	data, err := ca.GetCPSActionByActionCode(ctx, action_code, action.Department)
	if err != nil && data == nil {
		return errors.New(localization.ErrorActionNotFound.Code)
	}

	return ca.repo.Update(ctx, action_code, *action)
}
func (ca *cpsActionService) GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {

	return ca.repo.FindAllWithPagination(ctx, *filterParams, department)
}
func (ca *cpsActionService) GetCPSActionByID(ctx context.Context, id, department string) (*model.CPSAction, error) {

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		ca.logger.Errorf("their is error when try to parse the string to bson object in service")
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return ca.repo.FindOne(ctx, model.CPSAction{ID: objID})
}
func (ca *cpsActionService) GetCPSActionByUniqueID(ctx context.Context, requestAction, department string) (*model.CPSAction, error) {

	return ca.repo.FindOne(ctx, model.CPSAction{Department: department, ActionStatus: string(constants.Pending), RequestAction: requestAction})
}
func (ca *cpsActionService) GetCPSActionByActionCode(ctx context.Context, uniqueID, department string) (*model.CPSAction, error) {
	return ca.repo.FindOne(ctx, model.CPSAction{ActionCode: uniqueID, Department: department})
}

func (ca *cpsActionService) RollBack(ctx context.Context, action_code string) error {
	return ca.repo.Update(ctx, action_code, model.CPSAction{ActionStatus: string(constants.Pending)})
}

// func (s *cpsActionService) CPSActionExists(ctx context.Context, user model.CheckCPSAction) (bool, error) {
// 	_,err  := s.repo.FindOne(ctx, model.CPSAction{Department: user.Department,ActionStatus: string(constants.Pending),RequestAAction: })
// 	if err != nil {
// 		if err.Error() == localization.ErrorActionNotFound.Code{
// 			return false,nil
// 		}
// 		return false,nil
// 	}

// 	return true,nil
// }
