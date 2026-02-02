package permission

import (
	"cbe-super-app-cps-action/internal/constants/dto/permission"
	permission_int "cbe-super-app-cps-action/internal/constants/interfaces/permission"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	common_utils "cbe-super-app-cps-action/pkgs/utils"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type PermissionHandler struct {
	PermissionService service.PermissionService
	logger            utils.Logger
}

func InitPermissionHandler(svc service.PermissionService, logger utils.Logger) permission_int.PermissionHandler {
	return &PermissionHandler{
		PermissionService: svc,
		logger:            logger,
	}
}

// Create Permission Group
//
//	@Summary		Create Permission Group
//	@Description	Creates a new permission group
//	@Tags			Permission
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			body			body		permission.CreatePermissionGroupRequest	true	"Create Permission Group"
//	@Success		201				{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/permissions [post]
func (h *PermissionHandler) CreatePermissionGroup(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "permission", "PermissionHandler", "CreatePermissionGroup")
	defer span.End()
	var request permission.CreatePermissionGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		span.AddEvent("Failed to decode request", trace.WithAttributes(attribute.String("error", err.Error())))
		h.logger.Errorf("[CreatePermissionGroup] failed to decode request: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorFailedToDecodeRequest.Code)
		return
	}

	if err := request.Validate(); err != nil {
		span.AddEvent("Validation failed", trace.WithAttributes(attribute.String("error", err.Error())))
		h.logger.Warnf("[CreatePermissionGroup] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	userContext := common_utils.ExtractUserContext(r)
	if common_utils.IsIncomplete(userContext) {
		span.AddEvent("Incomplete user info", trace.WithAttributes(attribute.String("user_id", userContext.UserID)))
		localization.SendErrorByCodeResponse(w, localization.ErrorIncompleteUserInfo.Code)
		return
	}

	err := h.PermissionService.CreatePermissionGroup(ctx, request)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("user_id", userContext.UserID)))
		h.logger.Errorf("[CreatePermissionGroup] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.AddEvent("Permission group created", trace.WithAttributes(attribute.String("user_id", userContext.UserID)))
	h.logger.Infof("[CreatePermissionGroup] request sent successfully by user: %s", userContext.UserID)
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupRequestCreated, nil)
}

// Get Permission Groups
//
//	@Summary		Get Permission Groups
//	@Description	Retrieves a paginated list of permission groups. Searchable field: group_name.
//	@Tags			Permission
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page			query		int		false	"Page number"
//	@Param			per_page		query		int		false	"Items per page"
//	@Param			enabled			query		bool	false	"Filter by enabled status"
//	@Param			is_deleted		query		bool	false	"Filter by deleted status"
//	@Param			department_id	query		string	false	"Filter by department ID"
//	@Param			role			query		string	false	"Filter by role"
//	@Param			realm			query		string	false	"Filter by realm"
//	@Param			group_name		query		string	false	"Filter by group name"
//	@Param			search			query		string	false	"Search term (searches group_name)"
//	@Success		200				{object}	localization.StandardResponse{data=permission.PaginatedPermissionGroupResponse}
//	@Failure		400,500			{object}	localization.StandardResponse{data=nil}
//	@Router			/permissions [get]
func (h *PermissionHandler) GetPermissionGroups(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "permission", "PermissionHandler", "GetPermissionGroups")
	defer span.End()
	filterparams := common_utils.ExtractFilterParams(r)

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

	permissionGroups, err := h.PermissionService.GetPermissionGroups(ctx, filterparams)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		h.logger.Errorf("[GetPermissionGroups] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("Permission groups retrieved", trace.WithAttributes(attribute.Int("count", len(permissionGroups.Data))))
	h.logger.Infof("[GetPermissionGroups] retrieved %d permission groups", len(permissionGroups.Data))
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupsFetched, permissionGroups)
}

// //Get Permission Group by Name
//
//	//@Summary		Get Permission Group by Group_name
//	//@Description	Retrieves a permission  by group_name
//	//@Tags			Permission
//	//@Security		BearerAuth
//	//@Produce		json
//	//@Param			group_name			path		string	true	"group_name"
//	//@Success		200			{object}	localization.StandardResponse{data=permission.PermissionCategoryResponse}
//	//@Failure		400,404,500	{object}	localization.StandardResponse{data=nil}
//	//@Router			/permissions/{group_name} [get]
func (h *PermissionHandler) GetPermissionGroup(w http.ResponseWriter, r *http.Request) {
	_, span := common_utils.TraceLogger(r.Context(), "handler", "permission", "PermissionHandler", "GetPermissionGroup")
	defer span.End()
	groupName := chi.URLParam(r, "group_name")
	permissionGroup, err := h.PermissionService.GetPermissionGroup(groupName)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("group_name", groupName)))
		h.logger.Errorf("[GetPermissionGroup] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("Permission group retrieved", trace.WithAttributes(attribute.String("group_name", groupName)))
	h.logger.Infof("[GetPermissionGroup] permission group retrieved successfully for group_name: %s", groupName)
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupFetched, permissionGroup)
}

// Get Permission Group by ID
//
//	@Summary		Get Permission Group by ID
//	@Description	Retrieves a permission  by ID
//	@Tags			Permission
//	@Security		BearerAuth
//	@Produce		json
//	@Param			ID			path		string	true	"ID"
//	@Success		200			{object}	localization.StandardResponse{data=permission.PermissionCategoryResponse}
//	@Failure		400,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/permissions/by_id/{id} [get]
func (h *PermissionHandler) GetPermissionGroupById(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "permission", "PermissionHandler", "GetPermissionGroupById")
	defer span.End()
	id := chi.URLParam(r, "id")
	permissionGroup, err := h.PermissionService.GetPermissionGroupById(ctx, id)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		h.logger.Errorf("[GetPermissionGroupById] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	span.AddEvent("Permission group retrieved", trace.WithAttributes(attribute.String("id", id)))
	h.logger.Infof("[GetPermissionGroupById] permission group retrieved successfully for id: %s", id)
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupFetched, permissionGroup)
}

// Update Permission Group
//
//	@Summary		Update Permission Group
//	@Description	Updates an existing permission group
//	@Tags			Permission
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			group_name		path		string									true	"Group Name"
//	@Param			body			body		permission.UpdatePermissionGroupRequest	true	"Update Permission Group DTO"
//	@Success		200				{object}	localization.StandardResponse{data=nil}
//	@Failure		400,404,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/permissions/{group_name} [put]
func (h *PermissionHandler) UpdatePermissionGroup(w http.ResponseWriter, r *http.Request) {
	// Extract old group name from path first, then validate
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "permission", "PermissionHandler", "UpdatePermissionGroup")
	defer span.End()
	id := chi.URLParam(r, "id")
	if id == "" {
		span.AddEvent("Missing group name", trace.WithAttributes(attribute.String("error", "group name required")))
		localization.SendErrorResponse(w, localization.ErrorGroupNameRequired, nil, nil)
		return
	}

	var request permission.UpdatePermissionGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		span.AddEvent("Failed to decode request", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		h.logger.Errorf("[UpdatePermissionGroup] failed to decode request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	request.Id = id

	if err := request.Validate(); err != nil {
		span.AddEvent("Validation failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		h.logger.Warnf("[UpdatePermissionGroup] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	userContext := common_utils.ExtractUserContext(r)
	if common_utils.IsIncomplete(userContext) {
		span.AddEvent("Incomplete user info", trace.WithAttributes(attribute.String("user_id", userContext.UserID)))
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err := h.PermissionService.UpdatePermissionGroup(ctx, request)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("user_id", userContext.UserID), attribute.String("id", id)))
		h.logger.Errorf("[UpdatePermissionGroup] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.AddEvent("Permission group updated", trace.WithAttributes(attribute.String("user_id", userContext.UserID), attribute.String("id", id)))
	h.logger.Infof("[UpdatePermissionGroup] request sent successfully by user: %s", userContext.UserID)
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupRequestUpdated, nil)
}

// Get All Permission Categories With Permissions
//
//	@Summary		Get All Permission Categories With Permissions
//	@Description	Retrieves all permission categories with their permissions
//	@Tags			Permission
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200		{object}	localization.StandardResponse{data=permission.PaginatedPermissionGroupResponse}
//	@Failure		400,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/permissions/categories [get]
func (h *PermissionHandler) GetAllPermissionCategoriesWithPermissions(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "permission", "PermissionHandler", "GetAllPermissionCategoriesWithPermissions")
	defer span.End()
	categories, err := h.PermissionService.GetAllPermissionCategoriesWithPermissions(ctx)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error())))
		h.logger.Errorf("[GetAllPermissionCategoriesWithPermissions] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	// Group by category name and map fields to expected JSON keys
	grouped := map[string][]map[string]interface{}{}
	for _, c := range categories {
		key := c.CategoryName
		item := map[string]interface{}{
			"_id":           c.ID.Hex(),
			"category_name": c.CategoryName,
			"access":        c.Access,
			"permissions":   c.Permissions,
			"enabled":       c.Enabled,
			"is_deleted":    c.IsDeleted,
			"created_at":    c.CreatedAt,
			"updated_at":    c.UpdatedAt,
			"__v":           0,
		}
		grouped[key] = append(grouped[key], item)
	}

	span.AddEvent("Permission categories retrieved", trace.WithAttributes(attribute.Int("count", len(grouped))))
	h.logger.Infof("[GetAllPermissionCategoriesWithPermissions] retrieved %d permission categories", len(grouped))
	localization.SendSuccessResponse(w, localization.SuccessPermissionCategoriesFetched, grouped)
}

func (h *PermissionHandler) GetPermissionCategoriesByDepartment(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "permission", "PermissionHandler", "GetPermissionCategoriesByDepartment")
	defer span.End()
	departmentID := chi.URLParam(r, "department_id")
	if departmentID == "" {
		span.AddEvent("Missing department ID", trace.WithAttributes(attribute.String("error", "department ID required")))
		localization.SendErrorByCodeResponse(w, localization.ErrorDepartmentIDRequired.Code)
		return
	}
	permissionCategory, err := h.PermissionService.GetPermissionCategoriesByDepartment(ctx, departmentID)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("department_id", departmentID)))
		h.logger.Errorf("[GetPermissionCategoriesByDepartment] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.AddEvent("Permission categories retrieved", trace.WithAttributes(attribute.String("department_id", departmentID)))
	h.logger.Infof("[GetPermissionCategoriesByDepartment] retrieved permission categories for department_id: %s", departmentID)
	localization.SendSuccessResponse(w, localization.SuccessPermissionCategoriesFetched, permissionCategory)
}

func (h *PermissionHandler) GetPermissionGroupsByDepartment(w http.ResponseWriter, r *http.Request) {
	ctx, span := common_utils.TraceLogger(r.Context(), "handler", "permission", "PermissionHandler", "GetPermissionGroupsByDepartment")
	defer span.End()
	departmentID := chi.URLParam(r, "department_id")
	if departmentID == "" {
		span.AddEvent("Missing department ID", trace.WithAttributes(attribute.String("error", "department ID required")))
		localization.SendErrorByCodeResponse(w, localization.ErrorDepartmentIDRequired.Code)
		return
	}
	filterParam := common_utils.ExtractFilterParams(r)

	permissionGroups, err := h.PermissionService.GetPermissionGroupsByDepartment(ctx, departmentID, filterParam)
	if err != nil {
		span.AddEvent("Service error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("department_id", departmentID)))
		h.logger.Errorf("[GetPermissionGroupsByDepartment] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	span.AddEvent("Permission groups retrieved", trace.WithAttributes(attribute.Int("count", len(permissionGroups.Data)), attribute.String("department_id", departmentID)))
	h.logger.Infof("[GetPermissionGroupsByDepartment] retrieved %d permission groups for department_id: %s", len(permissionGroups.Data), departmentID)
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupsFetched, permissionGroups)
}
