package cpsaction

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	actionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/persistance"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type cpsActionService struct {
	repo       storage.CPSActionRepository
	logger     utils.Logger
	dispatcher Dispatcher
}

func NewCPSActionService(repo storage.CPSActionRepository, persistence persistance.Persistence, logger utils.Logger, dispatcher Dispatcher) service.CPSActionService {
	return &cpsActionService{
		repo:       repo,
		logger:     logger,
		dispatcher: dispatcher,
	}
}

func (ca *cpsActionService) CreateCPSAction(ctx context.Context, cpsAction *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateCPSAction", "CPSAction", "CreateCPSAction")
	defer span.End()

	existing, err := ca.GetCPSActionByUniqueID(ctx, cpsAction.RequestAction, cpsAction.Department)

	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		span.AddEvent("failed to get cps action by unique id", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if existing != nil {
		span.AddEvent("pending cps action exists", trace.WithAttributes(attribute.String("error", "pending cps action exists")))
		return errors.New(localization.ErrorPendingCpsActionExists.Code)
	}
	err = ca.repo.Save(ctx, cpsAction)
	if err != nil {
		span.AddEvent("failed to save cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}

func (ca *cpsActionService) ApproveCPSAction(ctx context.Context, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "ApproveCPSAction", "CPSAction", "ApproveCPSAction")
	defer span.End()

	data, err := ca.repo.Update(ctx, action.ActionCode, *action)
	if err != nil || data == nil {
		span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	approve, err := ca.dispatcher.Authorize(ctx, data)
	if err != nil && approve == nil {
		span.AddEvent("failed to authorize cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		RollErr := ca.RollBack(ctx, action)
		if RollErr != nil {
			span.AddEvent("failed to roll back cps action", trace.WithAttributes(attribute.String("error", RollErr.Error())))
			return RollErr
		}
		return err
	}
	return nil

}
func (ca *cpsActionService) RejectCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "RejectCPSAction", "CPSAction", "RejectCPSAction")
	defer span.End()
	_, err := ca.repo.Update(ctx, action_code, *action)
	if err != nil {
		span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}
func (ca *cpsActionService) GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionsByDepartment", "CPSAction", "GetCPSActionsByDepartment")
	defer span.End()
	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, department)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}
func (ca *cpsActionService) GetCPSActionByID(ctx context.Context, id, department string) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionByID", "CPSAction", "GetCPSActionByID")
	defer span.End()
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		span.AddEvent("failed to parse the string to bson object", trace.WithAttributes(attribute.String("error", err.Error())))
		ca.logger.Errorf("their is error when try to parse the string to bson object in service")
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	action, err := ca.repo.SanitizedFindOne(ctx, bson.M{"_id": objID, "department": department})
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}
func (ca *cpsActionService) GetCPSActionByUniqueID(ctx context.Context, requestAction, department string) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionByUniqueID", "CPSAction", "GetCPSActionByUniqueID")
	defer span.End()
	filter := bson.M{
		"department":     department,
		"action_status":  string(constants.Pending),
		"request_action": requestAction,
	}

	action, err := ca.repo.SanitizedFindOne(context.Background(), filter)
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}
func (ca *cpsActionService) GetCPSActionByActionCode(ctx context.Context, uniqueID, department string) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionByActionCode", "CPSAction", "GetCPSActionByActionCode")
	defer span.End()
	action, err := ca.repo.SanitizedFindOne(ctx, bson.M{"action_code": uniqueID})
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}

func (ca *cpsActionService) RollBack(ctx context.Context, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "RollBack", "CPSAction", "RollBack")
	defer span.End()
	err := ca.repo.UpdateCustome(ctx, bson.M{"action_code": action.ActionCode}, bson.M{"action_status": string(constants.Pending), "checker_id": "", "checker_name": "", "checker_phone_number": ""})
	if err != nil {
		span.AddEvent("failed to update custom", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}

func (ca *cpsActionService) GetActionCountsByDepartemnt(ctx context.Context, department string) (*actionDto.CPSActionCountResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetActionCountsByDepartemnt", "CPSAction", "GetActionCountsByDepartemnt")
	defer span.End()
	count, err := ca.repo.GetCountByDepartment(ctx, department)
	if err != nil {
		span.AddEvent("failed to get count by department", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return count, nil
}
