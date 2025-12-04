package department

import (
	"cbe-super-app-cps-action/internal/constants"
	department_dto "cbe-super-app-cps-action/internal/constants/dto/department"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
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

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	var actionMap interface{}
	b, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		fmt.Printf("failed to marshal CurrentAction: %v\n", err)
		return nil, fmt.Errorf("failed to marshal CurrentAction: %v", err)
	}
	fmt.Printf("JSON bytes: %s\n", string(b))
	err = json.Unmarshal(b, &actionMap)
	if err != nil {
		fmt.Printf("failed to unmarshal CurrentAction: %v\n", err)
		return nil, fmt.Errorf("failed to unmarshal to interface{}: %v", err)
	}

	actionData := department_core.Department_mapper(actionMap)
	if cpsAction.UniqueId != "" {
		objID, err := bson.ObjectIDFromHex(cpsAction.UniqueId)
		if err != nil {
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		actionData.ID = objID
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateDepartment):
		actionData.CreatedAt = time.Now()
		err := d.repo.Create(ctx, &actionData)
		if err != nil {
			d.logger.Errorf("Department Create action  failed", "error", err)
			return nil, err
		}
	case string(constants.RequestDeleteDepartment):
		err := d.repo.Delete(ctx, actionData.ID.Hex())
		if err != nil {
			d.logger.Errorf("Department Delete action  failed", "error", err)
			return nil, err
		}
	case string(constants.RequestEnableDisableDepartment):
		err := d.repo.EnableOrDisable(ctx, actionData.ID.Hex(), actionData.Enabled)
		if err != nil {
			d.logger.Errorf("Department Enable Disable action  failed", "error", err)
			return nil, err
		}
	case string(constants.RequestUpdateDepartment):
		err := d.repo.Update(ctx, actionData.ID.Hex(), &actionData)
		if err != nil {
			d.logger.Errorf("Department update action failed", "error", err)
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
		d.logger.Errorf("Create Department failed incomplete user data")
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
		return fmt.Errorf("%s", localization.ErrorDepartmentWithNameAlreadyExists.Code)
	}

	action := lib.CpsModelBuilder("", makerData, nil, new_department, string(constants.RequestCreateDepartment), constants.CREATE)

	err = d.cpsService.CreateCPSAction(ctx, &action)
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
	} else if err != nil {
		return err
	}

	if department.Enabled && enableDisable {
		return fmt.Errorf("%s", localization.ErrorDepartmentAlreadyEnabled.Code)
	}
	if !department.Enabled && !enableDisable {
		return fmt.Errorf("%s", localization.ErrorDepartmentAlreadyDisabled.Code)
	}

	new_department := *department
	new_department.Enabled = enableDisable

	action := lib.CpsModelBuilder(id, makerData, department, new_department, string(constants.RequestEnableDisableDepartment), constants.UPDATE)

	err = d.cpsService.CreateCPSAction(ctx, &action)
	if err != nil {
		return err
	}

	return nil
}

// GetAllDepartments implements service.DepartmentService.
func (d *DepartmentService) GetAllDepartments(ctx context.Context, filterParams *types.Filter) (types.PaginatedResponse[[]model.Department], error) {
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

	if department_request.Department != "" && department_request.Department != department.Department {
		existing_department, err := d.repo.FindByName(ctx, department_request.Department)
		code, _ := local_util.HandleMongoError(err)
		if code != localization.ErrorResourceNotFound.Code && existing_department != nil {
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
		return err
	}
	return nil
}
