package department

import (
	"cbe-super-app-cps-action/internal/constants"
	department_dto "cbe-super-app-cps-action/internal/constants/dto/department"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	department_core "cbe-super-app-cps-action/internal/service/department/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DepartmentService struct {
	portal_card      storage.PortalCardRepository
	permission_group storage.PermissionRepository
	repo             storage.DepartmentRepository
	cpsService       service.CPSActionService
	logger           utils.Logger
}

func NewDepartmentService(repo storage.DepartmentRepository, cpsService service.CPSActionService, portal_card storage.PortalCardRepository, permission_group storage.PermissionRepository, logger utils.Logger) service.DepartmentService {
	return &DepartmentService{
		repo:             repo,
		cpsService:       cpsService,
		portal_card:      portal_card,
		permission_group: permission_group,
		logger:           logger,
	}

}

// Authorize implements service.DepartmentService.
func (d *DepartmentService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Department", "Authorize")
	defer span.End()

	d.logger.Infof("[Authorize] authorizing department action: %s", cpsAction.RequestAction)
	var actionMap interface{}
	b, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		d.logger.Errorf("[Authorize] failed to marshal CurrentAction: %v", err)
		span.AddEvent("Failed to marshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, fmt.Errorf("failed to marshal CurrentAction: %v", err)
	}
	err = json.Unmarshal(b, &actionMap)
	if err != nil {
		d.logger.Errorf("[Authorize] failed to unmarshal CurrentAction: %v", err)
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		return nil, fmt.Errorf("failed to unmarshal to interface{}: %v", err)
	}

	actionData := department_core.Department_mapper(actionMap)
	if cpsAction.UniqueId != "" {
		objID, err := bson.ObjectIDFromHex(cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("Failed to parse ObjectID", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		actionData.ID = objID
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateDepartment):
		actionData.CreatedAt = time.Now()
		err := d.repo.Create(ctx, &actionData)
		if err != nil {
			d.logger.Errorf("[Authorize] department create action failed: %v", err)
			span.AddEvent("Department create action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", actionData.ID.Hex()),
			))
			return nil, err
		}
		d.logger.Infof("[Authorize] department created successfully with id: %s", actionData.ID.Hex())
	case string(constants.RequestDeleteDepartment):
		err := d.repo.Delete(ctx, actionData.ID.Hex())
		if err != nil {
			d.logger.Errorf("[Authorize] department delete action failed: %v", err)
			span.AddEvent("Department delete action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", actionData.ID.Hex()),
			))
			return nil, err
		}
		d.logger.Infof("[Authorize] department deleted successfully with id: %s", actionData.ID.Hex())
	case string(constants.RequestEnableDisableDepartment):
		err := d.repo.EnableOrDisable(ctx, actionData.ID.Hex(), actionData.Enabled)
		if err != nil {
			d.logger.Errorf("[Authorize] department enable/disable action failed: %v", err)
			span.AddEvent("Department enable/disable action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", actionData.ID.Hex()),
			))
			return nil, err
		}
		d.logger.Infof("[Authorize] department enable/disable action completed successfully for id: %s, enabled: %v", actionData.ID.Hex(), actionData.Enabled)
	case string(constants.RequestUpdateDepartment):
		err := d.repo.Update(ctx, actionData.ID.Hex(), &actionData)
		if err != nil {
			d.logger.Errorf("[Authorize] department update action failed: %v", err)
			span.AddEvent("Department update action failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", actionData.ID.Hex()),
			))
			return nil, err
		}
		d.logger.Infof("[Authorize] department updated successfully with id: %s", actionData.ID.Hex())
	default:
		d.logger.Errorf("[Authorize] unsupported action: %s", cpsAction.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.MsgDepartmentInvalidRequestAction),
			attribute.String("request_action", string(cpsAction.RequestAction)),
		))
		return nil, fmt.Errorf("%s", localization.MsgDepartmentInvalidRequestAction)
	}
	d.logger.Infof("[Authorize] department action authorized successfully: %s", cpsAction.RequestAction)
	return cpsAction, nil
}

// CreateDepartment implements service.DepartmentService.
func (d *DepartmentService) CreateDepartment(ctx context.Context, department department_dto.CreateDepartmentRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateDepartment", "Department", "CreateDepartment")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		d.logger.Errorf("Create Department failed incomplete user data")
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	new_department := model.Department{
		Department:  department.Department,
		PortalCards: department.PortalCards,
		Enabled:     true,
	}

	new_department.DepartmentCode = utils.RandomGenerator(20)
	existing_department, err := d.repo.FindByName(ctx, department.Department)
	code, _ := local_util.HandleMongoError(err)
	if code != localization.ErrorResourceNotFound.Code && existing_department != nil {
		span.AddEvent("Department with name already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorDepartmentWithNameAlreadyExists.Code),
			attribute.String("department", department.Department),
		))
		return fmt.Errorf("%s", localization.ErrorDepartmentWithNameAlreadyExists.Code)
	}

	action := lib.CpsModelBuilder("", makerData, nil, new_department, string(constants.RequestCreateDepartment), constants.CREATE)

	err = d.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}
	return nil
}

// EnableDisableDepartment implements service.DepartmentService.
func (d *DepartmentService) EnableDisableDepartment(ctx context.Context, id string, enableDisable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableDisableDepartment", "Department", "EnableDisableDepartment")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}
	department, err := d.repo.FindByID(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		d.logger.Errorf("[EnableDisableDepartment] department not found: %s", id)
		span.AddEvent("Department not found", trace.WithAttributes(
			attribute.String("error", code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", code)
	} else if err != nil {
		d.logger.Errorf("[EnableDisableDepartment] failed to find department: %v", err)
		span.AddEvent("Failed to find department", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	if department.Enabled && enableDisable {
		span.AddEvent("Department already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorDepartmentAlreadyEnabled.Code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", localization.ErrorDepartmentAlreadyEnabled.Code)
	}
	if !department.Enabled && !enableDisable {
		span.AddEvent("Department already disabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorDepartmentAlreadyDisabled.Code),
			attribute.String("id", id),
		))
		return fmt.Errorf("%s", localization.ErrorDepartmentAlreadyDisabled.Code)
	}

	new_department := *department
	new_department.Enabled = enableDisable

	action := lib.CpsModelBuilder(id, makerData, department, new_department, string(constants.RequestEnableDisableDepartment), constants.UPDATE)

	err = d.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		d.logger.Errorf("[EnableDisableDepartment] failed to create CPS action: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	d.logger.Infof("[EnableDisableDepartment] enable/disable request created successfully for id: %s, enabled: %v", id, enableDisable)
	return nil
}

// GetAllDepartments implements service.DepartmentService.
func (d *DepartmentService) GetAllDepartments(ctx context.Context, filterParams *types.Filter) (types.PaginatedResponse[[]model.Department], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllDepartments", "Department", "GetAllDepartments")
	defer span.End()

	result, err := d.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		d.logger.Errorf("[GetAllDepartments] failed to fetch departments: %v", err)
		span.AddEvent("Failed to fetch departments", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return result, err
	}
	d.logger.Infof("[GetAllDepartments] retrieved %d departments", len(result.Data))
	return result, nil
}

// GetDepartmentByID implements service.DepartmentService.
func (d *DepartmentService) GetDepartmentByID(ctx context.Context, id string) (*model.Department, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetDepartmentByID", "Department", "GetDepartmentByID")
	defer span.End()

	department, err := d.repo.FindByID(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		d.logger.Errorf("[GetDepartmentByID] department not found: %s", id)
		span.AddEvent("Department not found", trace.WithAttributes(
			attribute.String("error", code),
			attribute.String("id", id),
		))
		return nil, fmt.Errorf("%s", code)
	}
	if err != nil {
		d.logger.Errorf("[GetDepartmentByID] failed to fetch department: %v", err)
		span.AddEvent("Failed to fetch department", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	d.logger.Infof("[GetDepartmentByID] department retrieved successfully for id: %s", id)
	return department, nil

}

// UpdateDepartment implements service.DepartmentService.
func (d *DepartmentService) UpdateDepartment(ctx context.Context, id string, department_request department_dto.UpdateDepartmentRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateDepartment", "Department", "UpdateDepartment")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", constants.IncompleteUserInfo),
			attribute.String("id", id),
		))
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	department, err := d.repo.FindByID(ctx, id)
	if err != nil {
		d.logger.Errorf("[UpdateDepartment] failed to find department: %v", err)
		span.AddEvent("Failed to find department", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	if department_request.Department != "" && department_request.Department != department.Department {
		existing_department, err := d.repo.FindByName(ctx, department_request.Department)
		code, _ := local_util.HandleMongoError(err)
		if code != localization.ErrorResourceNotFound.Code && existing_department != nil {
			span.AddEvent("Department with name already exists", trace.WithAttributes(
				attribute.String("error", localization.ErrorDepartmentWithNameAlreadyExists.Code),
				attribute.String("id", id),
				attribute.String("department", department_request.Department),
			))
			return fmt.Errorf("%s", localization.ErrorDepartmentWithNameAlreadyExists.Code)
		}
	}

	updatedDepartment := *department

	if department_request.Department != "" {
		updatedDepartment.Department = department_request.Department
	}

	if len(department_request.PortalCards) > 0 {
		updatedDepartment.PortalCards = department_request.PortalCards
	}
	action := lib.CpsModelBuilder(id, makerData, department, updatedDepartment, string(constants.RequestUpdateDepartment), constants.UPDATE)

	err = d.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		d.logger.Errorf("[UpdateDepartment] failed to create CPS action: %v", err)
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	d.logger.Infof("[UpdateDepartment] department update request created successfully for id: %s", id)
	return nil
}
