package department

import (
	"errors"
	"gitlab.com/bersufekadgetachew/cbe-super-app-cps-action/internal/domain/department/entities"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"time"
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

func (s *Service) CreateDepartment(department string, portalCards []string, cpsAction entities.CPSAction) error {

	if ok, err := s.cpsActionRepo.CheckRequestExists(cpsAction); err != nil {
		return err
	} else if ok {
		return errors.New("You have a pending request for this action")
	}

	if exists, err := s.departmentRepo.CheckDepartmentExists(department); err != nil {
		return err
	} else if exists {
		return errors.New("Department already exists")
	}

	cpsAction.ActionCode = utils.Random(10, &utils.PreSufix{Prefix: "CPS_"})
	cpsAction.CurrentAction = map[string]interface{}{
		"department":      department,
		"department_code": utils.Random(10, &utils.PreSufix{Prefix: "DEP_"}),
		"portal_cards":    portalCards,
	}
	cpsAction.MakerActionTime = time.Now()

	if err := s.cpsActionRepo.CreateCPSAction(department, portalCards, cpsAction); err != nil {
		return err
	}

	s.logger.Infof("CPS action request created for department %s with action code %s", department, cpsAction.ActionCode)
	return nil

}

func (s *Service) ValidateActionRequest(actionCode string, userDept string) (*entities.CPSAction, error) {
	action, err := s.cpsActionRepo.FindByActionCode(actionCode)
	if err != nil {
		return nil, err
	}
	if action == nil {
		return nil, errors.New("action not found")
	}

	if action.Department != userDept {
		return nil, errors.New("You are not allowed to approve this request")
	}

	return action, nil
}
func (s *Service) ApproveActionRequest(actionCode string, user entities.CPSAction) error {
	if err := s.cpsActionRepo.ApproveActionRequest(actionCode, user); err != nil {
		return err
	}
	s.logger.Infof("Action request %s approved by user %s", actionCode, user)
	return nil
}

func (s *Service) UpdateDepartment(code string, department string, portalCards []string) error {
	if exists, err := s.departmentRepo.CheckDepartmentExists(department); err != nil {
		return err
	} else if !exists {
		return errors.New("Department does not exist")
	}

	return s.departmentRepo.UpdateDepartment(code, department, portalCards)
}
