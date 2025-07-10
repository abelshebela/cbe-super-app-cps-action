// Package department provides services and business logic for managing departments and related CPS actions.
package department

import (
	"context"
	"fmt"

	"log"

	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/department/entities"

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

func (s *Service) CheckRequestExists(ctx context.Context, cpsAction entities.CPSAction) (*entities.CPSAction, error) {
	action, err := s.cpsActionRepo.CheckRequestExists(ctx, cpsAction)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (s *Service) CheckDepartmentExists(ctx context.Context, department string) (bool, error) {
	return s.departmentRepo.CheckDepartmentExists(ctx, department)
}

func (s *Service) CreateCPSAction(ctx context.Context, department string, portalCards []string, cpsAction entities.CPSAction) (*entities.CPSAction, error) {

	cpsAction.ActionCode = utils.RandomGenerator(20)
	cpsAction.CurrentAction = map[string]any{
		"department":   department,
		"portal_cards": portalCards,
	}
	cpsAction.MakerActionTime = time.Now()

	if err := s.cpsActionRepo.CreateCPSAction(ctx, department, portalCards, cpsAction); err != nil {
		return nil, err
	}

	s.logger.Infof("CPS action request created for department %s with action code %s", department, cpsAction.ActionCode)
	return &cpsAction, nil

}
func (s *Service) CreateDepartment(ctx context.Context, department string, portalCards []string, cpsAction entities.CPSAction) error {
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

func (s *Service) ValidateActionRequest(ctx context.Context, actionCode string, userDept string) (*entities.CPSAction, error) {
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
func (s *Service) ApproveActionRequest(ctx context.Context, actionCode string, user entities.CPSAction) error {

	if err := s.cpsActionRepo.ApproveActionRequest(ctx, actionCode, user); err != nil {
		return err
	}
	s.logger.Infof("Action request %s approved by user from department %s", actionCode, user.Department)

	return nil
}

func (s *Service) RejectActionRequest(ctx context.Context, actionCode string, user entities.CPSAction) error {
	if err := s.cpsActionRepo.RejectActionRequest(ctx, actionCode, user); err != nil {
		return err
	}
	s.logger.Infof("Action request %s approved by user from department %s", actionCode, user.Department)

	return nil
}

func (s *Service) UpdateDepartment(ctx context.Context, code string, department string, portalCards []string) (*entities.Department, error) {
	data, err := s.departmentRepo.UpdateDepartment(ctx, code, department, portalCards)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// CreateDepartmentUpdateCPSAction creates a CPS action for department updates
func (s *Service) CreateDepartmentUpdateCPSAction(ctx context.Context, req entities.CPSAction) (*entities.CPSAction, error) {
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

// ApproveDepartmentUpdate approves a department update CPS action
func (s *Service) ApproveDepartmentUpdate(ctx context.Context, cpsAction entities.CPSAction) error {
	if err := s.cpsActionRepo.ApproveDepartmentUpdate(ctx, cpsAction); err != nil {
		return err
	}
	s.logger.Infof("Department update action %s approved by user from department %s", cpsAction.ActionCode, cpsAction.Department)
	return nil
}

// RejectDepartmentUpdate rejects a department update CPS action
func (s *Service) RejectDepartmentUpdate(ctx context.Context, cpsAction entities.CPSAction) error {
	if err := s.cpsActionRepo.RejectDepartmentUpdate(ctx, cpsAction); err != nil {
		return err
	}
	s.logger.Infof("Department update action %s rejected by user from department %s", cpsAction.ActionCode, cpsAction.Department)
	return nil
}
