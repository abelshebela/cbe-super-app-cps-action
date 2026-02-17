package department

import (
	department_dto "cbe-super-app-cps-action/internal/constants/dto/department"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"
	"strings"

	_ "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
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

// GetAllDepartments
//
//	@Summary		Get All Departments
//	@Description	Retrieve all departments with pagination, filtering, and search. Searchable fields: department_code, department, created_at, last_modified_at.
//	@Tags			Department
//	@Accept			json
//	@Produce		json
//	@Param			page			query		int														false	"Page number"
//	@Param			per_page		query		int														false	"Items per page"
//	@Param			search			query		string													false	"Search term (searches department_code, department, created_at, last_modified_at)"
//	@Param			department_code	query		string													false	"Filter by department code"
//	@Param			department		query		string													false	"Filter by department name"
//	@Param			enabled			query		bool													false	"Filter by enabled status"
//	@Param			is_deleted		query		bool													false	"Filter by deleted status"
//	@Param			created_at		query		string													false	"Filter by created at"
//	@Param			last_modified	query		string													false	"Filter by last modified"
//	@Success		200				{object}	localization.StandardResponse{data=[]model.Department}	"Departments retrieved successfully"
//	@Failure		400				{object}	localization.StandardResponse{data=nil}					"Bad request"
//	@Failure		500				{object}	localization.StandardResponse{data=nil}					"Internal server error"
//	@Security		BearerAuth
//	@Router			/departments [get]
func (d *DepartmentHandler) GetAllDepartments(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "getAllDepartments", "handler", "department")
	defer span.End()
	filterParams := common_utils.ExtractFilterParams(r)

	search := r.URL.Query().Get("search")
	filter := r.URL.Query().Get("filter")

	if err := common_utils.NoSpecialChars(search); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	if err := common_utils.NoSpecialChars(filter); err != nil {
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	departments, err := d.departmentService.GetAllDepartments(ctx, filterParams)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("[GetAllDepartments] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.Int("department.count", len(departments.Data)))
	d.logger.Infof("[GetAllDepartments] retrieved %d departments", len(departments.Data))
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
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "createDepartment", "handler", "department")
	defer span.End()
	var departmentRequest department_dto.CreateDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&departmentRequest); err != nil {
		span.RecordError(err)
		d.logger.Errorf("[CreateDepartment] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	if err := departmentRequest.Validate(); err != nil {
		span.RecordError(err)
		d.logger.Errorf("[CreateDepartment] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("department.name", departmentRequest.Department))

	err := d.departmentService.CreateDepartment(ctx, departmentRequest)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("[CreateDepartment] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	d.logger.Infof("[CreateDepartment] request sent successfully for department: %s", departmentRequest.Department)
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
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "updateDepartment", "handler", "department")
	defer span.End()
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		d.logger.Errorf("[UpdateDepartmentRequest] missing department ID")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	var departmentRequest department_dto.UpdateDepartmentRequest

	if err := json.NewDecoder(r.Body).Decode(&departmentRequest); err != nil {
		span.RecordError(err)
		d.logger.Errorf("[UpdateDepartmentRequest] decode: %v", err)
		localization.SendBadRequestResponse(w, localization.ErrorInvalidRequestBody.Code)
		return
	}

	if err := departmentRequest.Validate(); err != nil {
		span.RecordError(err)
		d.logger.Errorf("[UpdateDepartmentRequest] validation: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	span.SetAttributes(attribute.String("department.id", id))

	err := d.departmentService.UpdateDepartment(ctx, id, departmentRequest)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("[UpdateDepartmentRequest] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	d.logger.Infof("[UpdateDepartmentRequest] request sent successfully for id: %s", id)
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
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "getDepartmentById", "handler", "department")
	defer span.End()
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		d.logger.Errorf("[GetDepartmentByID] missing department ID")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("department.id", id))

	department, err := d.departmentService.GetDepartmentByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("[GetDepartmentByID] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	d.logger.Infof("[GetDepartmentByID] department retrieved successfully for id: %s", id)
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
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "enableDepartment", "handler", "department")
	defer span.End()
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		d.logger.Errorf("[EnableDepartment] missing department ID")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("department.id", id))

	err := d.departmentService.EnableDisableDepartment(ctx, id, true)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("[EnableDepartment] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	d.logger.Infof("[EnableDepartment] request sent successfully for id: %s", id)
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
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "disableDepartment", "handler", "department")
	defer span.End()
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		d.logger.Errorf("[DisableDepartment] missing department ID")
		localization.SendErrorResponse(w, localization.ErrorRequiredFieldMissing, nil, nil)
		return
	}

	span.SetAttributes(attribute.String("department.id", id))

	err := d.departmentService.EnableDisableDepartment(ctx, id, false)
	if err != nil {
		span.RecordError(err)
		d.logger.Errorf("[DisableDepartment] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	d.logger.Infof("[DisableDepartment] request sent successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessDepartmentDisableRequestCreated, nil)
}
