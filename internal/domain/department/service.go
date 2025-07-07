// Package department provides services and business logic for managing departments and related CPS actions.
package department

import (
	"context"
	"fmt"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"
	"time"

	err_msg "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

const (
	CPSPrefix  = "CPS_"
	DeptPrefix = "DEP_"
)

type Service struct {
	cpsActionRepo  CPSActionRepository
	departmentRepo DepartmentRepository
	logger         utils.Logger
}

func InitDepartmentDomain(cpsActionRepo CPSActionRepository, departmentRepo DepartmentRepository, logger utils.Logger) *Service {
	return &Service{
		cpsActionRepo:  cpsActionRepo,
		departmentRepo: departmentRepo,
		logger:         logger,
	}
}

func (s *Service) CheckRequestExists(ctx context.Context, cpsAction entities.CPSAction) (bool, error) {
	return s.cpsActionRepo.CheckRequestExists(ctx, cpsAction)
}

func (s *Service) CheckDepartmentExists(ctx context.Context, department string) (bool, error) {
	return s.departmentRepo.CheckDepartmentExists(ctx, department)
}

func (s *Service) CreateDepartment(ctx context.Context, department string, portalCards []string, cpsAction entities.CPSAction) error {

	cpsAction.ActionCode = utils.Random(10, &utils.PreSufix{Prefix: CPSPrefix})
	cpsAction.CurrentAction = map[string]any{
		"department":      department,
		"department_code": utils.Random(10, &utils.PreSufix{Prefix: DeptPrefix}),
		"portal_cards":    portalCards,
	}
	cpsAction.MakerActionTime = time.Now()

	if err := s.cpsActionRepo.CreateCPSAction(ctx, department, portalCards, cpsAction); err != nil {
		return err
	}

	s.logger.Infof("CPS action request created for department %s with action code %s", department, cpsAction.ActionCode)
	return nil

}

func (s *Service) ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*entities.CPSAction, error) {
	action, err := s.cpsActionRepo.FindByActionCode(ctx, actionCode)
	if err != nil {
		return nil, err
	}

	if action.Department != userDept {
		return nil, fmt.Errorf(err_msg.ActionNotAllowed)
	}

	return action, nil
}
func (s *Service) ApproveActionRequest(ctx context.Context, actionCode string, user entities.CPSAction) error {
	if err := s.cpsActionRepo.ApproveActionRequest(ctx, actionCode, user); err != nil {
		return err
	}
	s.logger.Infof("Action request %s approved by user from department %s", actionCode, user.Department)

	return nil
}

func (s *Service) UpdateDepartment(ctx context.Context, code string, department string, portalCards []string) error {
	return s.departmentRepo.UpdateDepartment(ctx, code, department, portalCards)
}
