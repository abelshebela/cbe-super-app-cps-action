package department

type CPSActionRepository interface {
	CheckRequestExists(action CPSAction) (bool, error)
	CreateCPSAction(department string, portalCards []string, action CPSAction) error
	FindByActionCode(code string) (*entities.CPSAction, error)
	UpdateActionStatus(actionCode string, status string) error
}

type DepartmentRepository interface {
	CheckDepartmentExists(department string) (bool, error)
	CreateDepartment(dept entities.Department) error
	UpdateDepartment(code string, department string, portalCards []string) error

}
