package department

import (
	department_dto "cbe-super-app-cps-action/internal/constants/dto/department"
	"cbe-super-app-cps-action/internal/constants/localization"
	department_core "cbe-super-app-cps-action/internal/handlers/rest/http/department/core"
	"cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type DepartmentHandler struct {
	departmentService service.DepartmentService
	logger            utils.Logger
}

func NewDepartmentHandler(departmentService service.DepartmentService, logger utils.Logger) DepartmentHandler {
	return DepartmentHandler{
		departmentService: departmentService,
		logger:            logger,
	}
}

func (d *DepartmentHandler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	filterParams := common_utils.ExtractFilterParams(r)

	departments, err := d.departmentService.GetAllDepartments(r.Context(), filterParams)
	if err != nil {
		d.logger.Errorf("get all Departments failed", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessGetAllDepartments, departments)
}
func (d *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var departmentRequest department_dto.CreateDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&departmentRequest); err != nil {
		d.logger.Errorf("failed to decode department request", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	if response_code := department_core.ValidateDepartmentRequest(r, departmentRequest); response_code.Code != "" {
		d.logger.Errorf("invalid input", response_code)
		localization.SendErrorResponse(w, response_code, nil, nil)
		return
	}
	err := d.departmentService.CreateDepartment(r.Context(), departmentRequest)
	if err != nil {
		d.logger.Errorf("department create request failed", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessDepartmentCreateRequestCreated, nil)
}
func (d *DepartmentHandler) UpdateDepartmentRequest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	var departmentRequest department_dto.UpdateDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&departmentRequest); err != nil {
		d.logger.Errorf("failed to decode department request", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	if response_code := department_core.ValidateDepartmentRequest(r, departmentRequest); response_code.Code != "" {
		d.logger.Errorf("invalid input", response_code)
		localization.SendErrorResponse(w, response_code, nil, nil)
		return
	}

	err := d.departmentService.UpdateDepartment(r.Context(), id, departmentRequest)
	if err != nil {
		d.logger.Errorf("department update request failed", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessDepartmentUpdateRequestCreated, nil)
}
func (d *DepartmentHandler) GetDepartmentByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	department, err := d.departmentService.GetDepartmentByID(r.Context(), id)
	if err != nil {
		d.logger.Errorf("get department by id", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessGetDepartments, department)
}
func (d *DepartmentHandler) EnableDepartment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	err := d.departmentService.EnableDisableDepartment(r.Context(), id, true)
	if err != nil {
		d.logger.Errorf("enable department by id failed", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessDepartmentEnableRequestCreated, nil)
}
func (d *DepartmentHandler) DisableDepartment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		d.logger.Errorf("missing or invalid parameter 'id'")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	err := d.departmentService.EnableDisableDepartment(r.Context(), id, false)
	if err != nil {
		d.logger.Errorf("disable department by id failed", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessDepartmentDisableRequestCreated, nil)
}
