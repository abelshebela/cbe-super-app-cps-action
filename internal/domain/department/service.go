// Package department provides services and business logic for managing departments and related CPS actions.
package department

import (
	"context"
	"fmt"

	"log"

	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/permission"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	err_msg "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"go.mongodb.org/mongo-driver/v2/bson"

	cpsactions "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	CPSPrefix  = "CPS_"
	DeptPrefix = "DEP_"
)

type UpdateDepartmentRequest struct {
	Department       string   `json:"department"`
	PortalCards      []string `json:"portal_cards"`
	PermissionGroups []string `json:"permission_groups"`
}

type Service struct {
	cpsActionRepo     CPSActionRepository
	departmentRepo    DepartmentRepository
	permissionService permission.PermissionDomainService
	logger            utils.Logger
}

func InitDepartmentDomain(
	cpsActionRepo CPSActionRepository,
	departmentRepo DepartmentRepository,
	permissionService permission.PermissionDomainService,
	logger utils.Logger,
) *Service {
	return &Service{
		cpsActionRepo:     cpsActionRepo,
		departmentRepo:    departmentRepo,
		permissionService: permissionService,
		logger:            logger,
	}
}

func (s *Service) CheckRequestExists(ctx context.Context, cpsAction cpsactions.CPSAction) (*cpsactions.CPSAction, error) {
	action, err := s.cpsActionRepo.CheckRequestExists(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (s *Service) CheckDepartmentExists(ctx context.Context, department string) (bool, error) {
	return s.departmentRepo.CheckDepartmentExists(ctx, department)
}

func (s *Service) CheckDepartmentExistsByID(ctx context.Context, id string) (bool, error) {
	return s.departmentRepo.CheckDepartmentExistsByID(ctx, id)
}

// Change the return type of CreateCPSAction to (string, error)
func (s *Service) CreateCPSAction(ctx context.Context, department string, portalCards []string, permissionGroups []string, cpsAction cpsactions.CPSAction) (string, error) {
	groupIDs := permissionGroups
	_, err := s.permissionService.ValidatePermissionGroups(groupIDs)
	if err != nil {
		return "", fmt.Errorf("invalid permission groups: %w", err)
	}

	cpsAction.ActionCode = utils.RandomGenerator(20)
	cpsAction.CurrentAction = map[string]any{
		"department":        department,
		"portal_cards":      portalCards,
		"permission_groups": permissionGroups,
	}
	cpsAction.MakerActionTime = time.Now()

	if err := s.cpsActionRepo.CreateCPSAction(ctx, department, portalCards, permissionGroups, cpsAction); err != nil {
		return "", err
	}

	s.logger.Infof("CPS action request created for department %s with action code %s", department, cpsAction.ActionCode)
	return cpsAction.ActionCode, nil
}
func (s *Service) CreateDepartment(ctx context.Context, department string, portalCards []string, cpsAction cpsactions.CPSAction) error {
	dept := entities.Department{
		DepartmentCode: utils.RandomGenerator(20),
		Department:     department,
		PortalCards:    portalCards,
		CreatedAt:      time.Now(),
		LastModified:   time.Now(),
	}
	if err := s.departmentRepo.CreateDepartment(ctx, dept); err != nil {
		return err
	}

	return nil
}

func (s *Service) ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*cpsactions.CPSAction, error) {
	action, err := s.cpsActionRepo.FindByActionCode(ctx, actionCode)
	if err != nil {
		log.Println("error", err)
		return nil, err
	}

	if action.Department != userDept {
		return nil, fmt.Errorf(err_msg.ActionNotAllowed)
	}

	return action, nil
}
func (s *Service) ApproveActionRequest(ctx context.Context, actionCode string, user cpsactions.CPSAction) error {

	if err := s.cpsActionRepo.ApproveActionRequest(ctx, actionCode, user); err != nil {
		return err
	}
	s.logger.Infof("Action request %s approved by user from department %s", actionCode, user.Department)

	return nil
}

func (s *Service) RejectActionRequest(ctx context.Context, actionCode string, user cpsactions.CPSAction) error {
	if err := s.cpsActionRepo.RejectActionRequest(ctx, actionCode, user); err != nil {
		return err
	}
	s.logger.Infof("Action request %s approved by user from department %s", actionCode, user.Department)

	return nil
}

func (s *Service) UpdateDepartment(ctx context.Context, id string, req UpdateDepartmentRequest) (*entities.Department, error) {
	for _, idStr := range req.PermissionGroups {
		if _, err := bson.ObjectIDFromHex(idStr); err != nil {
			return nil, fmt.Errorf("invalid permission group id: %s", idStr)
		}
	}
	if _, err := s.permissionService.ValidatePermissionGroups(req.PermissionGroups); err != nil {
		return nil, fmt.Errorf("invalid permission groups: %w", err)
	}
	data, err := s.departmentRepo.UpdateDepartment(ctx, id, req.Department, req.PortalCards, req.PermissionGroups)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// CreateDepartmentUpdateCPSAction creates a CPS action for department updates
func (s *Service) CreateDepartmentUpdateCPSAction(ctx context.Context, req cpsactions.CPSAction) (*cpsactions.CPSAction, error) {
	// Set up the CPS action for department update
	req.ActionCode = utils.RandomGenerator(20)
	req.MakerActionTime = time.Now()
	req.CreatedAt = time.Now()
	req.LastModifiedAt = time.Now()

	cpsAction, err := s.cpsActionRepo.CreateDepartmentUpdateCPSAction(ctx, req)
	if err != nil {
		return nil, err
	}

	s.logger.Infof("Department update CPS action created with action code %s", cpsAction.ActionCode)
	return cpsAction, nil
}

// General CPS Action Approve/Reject integration
func (s *Service) Authorize(ctx context.Context, action *cpsactions.CPSAction) (*cpsactions.CPSAction, error) {
	// fmt.Println("Authorizing action", action.ActionCode, "for department", action.Department)
	if err := s.cpsActionRepo.ApproveDepartmentUpdate(ctx, *action); err != nil {
		return nil, err
	}
	s.logger.Infof("Department action %s approved by user from department %s", action.ActionCode, action.Department)
	return action, nil
}

func (s *Service) Reject(ctx context.Context, action *cpsactions.CPSAction) (*cpsactions.CPSAction, error) {
	if err := s.cpsActionRepo.RejectDepartmentUpdate(ctx, *action); err != nil {
		return nil, err
	}
	s.logger.Infof("Department action %s rejected by user from department %s", action.ActionCode, action.Department)
	return action, nil
}

func (s *Service) GetAllDepartments(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Department], error) {
	departments, err := s.departmentRepo.GetAllDepartments(ctx, filterParams)
	if err != nil {
		return nil, err
	}
	return departments, nil
}

func (s *Service) GetDepartmentByID(ctx context.Context, id string) (*entities.Department, error) {
	return s.departmentRepo.GetDepartmentByID(ctx, id)
}
