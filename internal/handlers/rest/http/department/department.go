package department

import (
	department_dto "cbe-super-app-cps-action/internal/constants/dto/department"
	"cbe-super-app-cps-action/internal/constants/localization"
	_ "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"
	"strings"

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

// GetDepartments godoc
//
//	@Summary		List departments
//	@Description	Retrieve departments with pagination and optional search
//	@Tags			Department
//	@Accept			json
//	@Produce		json
//	@Param			page		query		int																		false	"Page number"		default(1)
//	@Param			per_page	query		int																		false	"Items per page"	default(10)
//	@Param			search		query		string																	false	"Search term"
//	@Success		200			{object}	localization.StandardResponse{data=[]model.PaginatedDepartmentResponse}	"Departments retrieved successfully"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}									"Internal server error"
//	@Security		BearerAuth
//	@Router			/departments [get]
func (d *DepartmentHandler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	filterParams := common_utils.ExtractFilterParams(r)

	departments, err := d.departmentService.GetAllDepartments(r.Context(), filterParams)
	if err != nil {
		d.logger.Errorf("[GetAllDepartments] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessGetAllDepartments, departments)
}

// CreateDepartment godoc
//
//	@Summary		Create a new department
//	@Description	Create a new department with the provided information
//	@Tags			Department
//	@Accept			json
//	@Produce		json
//	@Param			department	body		department_dto.CreateDepartmentRequest	true	"Department information"
//	@Success		200			{object}	localization.StandardResponse{data=nil}	"Department created successfully"
//	@Failure		400			{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500			{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/departments [post]
func (d *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var departmentRequest department_dto.CreateDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&departmentRequest); err != nil {
		d.logger.Errorf("[CreateDepartment] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	if err := departmentRequest.Validate(); err != nil {
		d.logger.Errorf("[CreateDepartment] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	err := d.departmentService.CreateDepartment(r.Context(), departmentRequest)
	if err != nil {
		d.logger.Errorf("[CreateDepartment] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessDepartmentCreateRequestCreated, nil)
}

// UpdateDepartmentRequest handles HTTP requests to update an existing department.
//
//	@Summary		Update Department
//	@Description	Updates the details of an existing department by its ID.
//	@Tags			Department
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string									true	"Department ID"
//	@Param			body	body		department_dto.UpdateDepartmentRequest	true	"Update Department Request"
//	@Success		200		{object}	localization.StandardResponse{data=nil}	"Department updated successfully"
//	@Failure		400		{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		500		{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/departments/{id} [patch]
func (d *DepartmentHandler) UpdateDepartmentRequest(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		d.logger.Errorf("[UpdateDepartmentRequest] missing department ID")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	var departmentRequest department_dto.UpdateDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&departmentRequest); err != nil {
		d.logger.Errorf("[UpdateDepartmentRequest] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	if err := departmentRequest.Validate(); err != nil {
		d.logger.Errorf("[UpdateDepartmentRequest] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	err := d.departmentService.UpdateDepartment(r.Context(), id, departmentRequest)
	if err != nil {
		d.logger.Errorf("[UpdateDepartmentRequest] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessDepartmentUpdateRequestCreated, nil)
}

// GetDepartment godoc
//
//	@Summary		Get department by ID
//	@Description	Retrieve a department's details by ID
//	@Tags			Department
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string													true	"Department ID"
//	@Success		200	{object}	localization.StandardResponse{data=model.Department}	"Department retrieved successfully"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}					"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}					"Department not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/departments/{id} [get]
func (d *DepartmentHandler) GetDepartmentByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		d.logger.Errorf("[GetDepartmentByID] missing department ID")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	department, err := d.departmentService.GetDepartmentByID(r.Context(), id)
	if err != nil {
		d.logger.Errorf("[GetDepartmentByID] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessGetDepartments, department)
}

// EnableDepartment godoc
//
//	@Summary		Enable a department
//	@Description	Enable a department by ID
//	@Tags			Department
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Department ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Department enable request submitted"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Department not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/departments/enable/{id} [patch]
func (d *DepartmentHandler) EnableDepartment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		d.logger.Errorf("[EnableDepartment] missing department ID")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	err := d.departmentService.EnableDisableDepartment(r.Context(), id, true)
	if err != nil {
		d.logger.Errorf("[EnableDepartment] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessDepartmentEnableRequestCreated, nil)
}

// DisableDepartment godoc
//
//	@Summary		Disable a department
//	@Description	Disable a department by ID
//	@Tags			Department
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string									true	"Department ID"
//	@Success		200	{object}	localization.StandardResponse{data=nil}	"Department disable request submitted"
//	@Failure		400	{object}	localization.StandardResponse{data=nil}	"Bad request"
//	@Failure		404	{object}	localization.StandardResponse{data=nil}	"Department not found"
//	@Failure		500	{object}	localization.StandardResponse{data=nil}	"Internal server error"
//	@Security		BearerAuth
//	@Router			/departments/disable/{id} [patch]
func (d *DepartmentHandler) DisableDepartment(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		d.logger.Errorf("[DisableDepartment] missing department ID")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	err := d.departmentService.EnableDisableDepartment(r.Context(), id, false)
	if err != nil {
		d.logger.Errorf("[DisableDepartment] service: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessDepartmentDisableRequestCreated, nil)
}
