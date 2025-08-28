package department

import (
	"cbe-super-app-cps-action/internal/constants"
	department_dto "cbe-super-app-cps-action/internal/constants/dto/department"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"fmt"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type DepartmentService struct {
	portal_card      storage.PortalCardRepository
	permission_group storage.PermissionGroupRepository
	repo             storage.DepartmentRepository
	cpsService       service.CPSActionService
	logger           utils.Logger
}

func NewDepartmentService(repo storage.DepartmentRepository, cpsService service.CPSActionService, portal_card storage.PortalCardRepository, permission_group storage.PermissionGroupRepository, logger utils.Logger) service.DepartmentService {
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
	var actionData model.Department
	raw, _ := bson.Marshal(cpsAction.CurrentAction)
	if err := bson.Unmarshal(raw, &actionData); err != nil {
		d.logger.Errorf("failed to unmarshal action data for authorization, action_code: %s", (cpsAction.ActionCode))
		return nil, fmt.Errorf("%s", localization.MsgDepartmentInvalidRequestAction)
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateDepartment):
		actionData.CreatedAt = time.Now()
		err := d.repo.Create(ctx, &actionData)
		if err != nil {
			d.logger.Errorf("Bank Create action  failed", "error", err)
			return nil, err
		}
	case string(constants.RequestDeleteDepartment):
		err := d.repo.Delete(ctx, actionData.ID)
		if err != nil {
			d.logger.Errorf("Bank Delete action  failed", "error", err)
			return nil, err
		}
	case string(constants.RequestEnableDisableDepartment):
		err := d.repo.EnableOrDisable(ctx, actionData.ID, actionData.Enabled)
		if err != nil {
			d.logger.Errorf("Bank Enable Disable action  failed", "error", err)
			return nil, err
		}
	case string(constants.RequestUpdateDepartment):
		err := d.repo.Update(ctx, actionData.ID, &actionData)
		if err != nil {
			d.logger.Errorf("Bank update action failed", "error", err)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%s", localization.MsgDepartmentInvalidRequestAction)
	}
	return cpsAction, nil
}

// CreateDepartment implements service.DepartmentService.
func (d *DepartmentService) CreateDepartment(ctx context.Context, department department_dto.CreateDepartmentRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		d.logger.Errorf("Create Bank failed incomplete user data")
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	new_department := model.Department{
		Department:       department.Department,
		PortalCards:      department.PortalCards,
		PermissionGroups: department.PermissionGroups,
	}
	if all_valid, err := d.permission_group.ValidatePermissionGroupByID(ctx, department.PermissionGroups); !all_valid || err != nil {
		return fmt.Errorf("%s", localization.ErrorInvalidDepartmentPermissionGroup.Code)
	}

	if all_valid, err := d.portal_card.ValidatePortalCardByID(ctx, department.PortalCards); !all_valid || err != nil {
		return fmt.Errorf("%s", localization.ErrorInvalidDepartmentPortalCard.Code)
	}

	action := lib.CpsModelBuilder("", makerData, nil, new_department, string(constants.RequestCreateDepartment), constants.CREATE)

	err := d.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}
	return nil
}

// EnableDisableDepartment implements service.DepartmentService.
func (d *DepartmentService) EnableDisableDepartment(ctx context.Context, id string, enableDisable bool) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return fmt.Errorf(constants.IncompleteUserInfo)
	}
	department, err := d.repo.FindByID(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		return fmt.Errorf("%s", code)
	}

	if department.Enabled && enableDisable {
		return fmt.Errorf("%s", localization.ErrorDepartmentAlreadyEnabled.Code)
	}
	if !department.Enabled && !enableDisable {
		return fmt.Errorf("%s", localization.ErrorDepartmentAlreadyEnabled.Code)
	}

	new_department := *department
	new_department.Enabled = enableDisable

	action := lib.CpsModelBuilder(id, makerData, department, new_department, string(constants.RequestUpdateBank), constants.UPDATE)

	err = d.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}

	return nil
}

// GetAllDepartments implements service.DepartmentService.
func (d *DepartmentService) GetAllDepartments(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Department], error) {
	return d.repo.FindAllWithPagination(ctx, *filterParams)
}

// GetDepartmentByID implements service.DepartmentService.
func (d *DepartmentService) GetDepartmentByID(ctx context.Context, id string) (*model.Department, error) {
	department, err := d.repo.FindByID(ctx, id)
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		return nil, fmt.Errorf("%s", code)
	}
	return department, nil

}

// UpdateDepartment implements service.DepartmentService.
func (d *DepartmentService) UpdateDepartment(ctx context.Context, id string, department_request department_dto.UpdateDepartmentRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return fmt.Errorf(constants.IncompleteUserInfo)
	}

	department, err := d.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	code, _ := local_util.HandleMongoError(err)
	if code == localization.ErrorResourceNotFound.Code {
		return fmt.Errorf("%s", code)
	}

	updatedDepartment := department

	if department.Department != "" {
		updatedDepartment.Department = department_request.Department
	}

	if department.PortalCards != nil {
		updatedDepartment.PortalCards = department_request.PortalCards
	}

	if department.PermissionGroups != nil {
		updatedDepartment.PermissionGroups = department_request.PermissionGroups
	}

	if all_valid, err := d.permission_group.ValidatePermissionGroupByID(ctx, department.PermissionGroups); !all_valid || err != nil {
		return fmt.Errorf("%s", localization.ErrorInvalidDepartmentPermissionGroup.Code)
	}

	if all_valid, err := d.portal_card.ValidatePortalCardByID(ctx, department_request.PortalCards); !all_valid || err != nil {
		return fmt.Errorf("%s", localization.ErrorInvalidDepartmentPortalCard.Code)
	}

	action := lib.CpsModelBuilder(id, makerData, department, updatedDepartment, string(constants.RequestUpdateBank), constants.UPDATE)

	err = d.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}
	return nil
}
