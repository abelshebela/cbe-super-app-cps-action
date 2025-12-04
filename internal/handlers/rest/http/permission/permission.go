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
//	@Param			body			body		permission.CreatePermissionGroupRequest	true	"Create Permission Group DTO"
//	@Success		201				{object}	localization.StandardResponse{data=nil}
//	@Failure		400,401,422,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/permissions [post]
func (h *PermissionHandler) CreatePermissionGroup(w http.ResponseWriter, r *http.Request) {
	var request permission.CreatePermissionGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[CreatePermissionGroup] failed to decode request: %v", err)
		localization.SendErrorByCodeResponse(w, localization.ErrorFailedToDecodeRequest.Code)
		return
	}

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[CreatePermissionGroup] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	userContext := common_utils.ExtractUserContext(r)
	if common_utils.IsIncomplete(userContext) {
		localization.SendErrorByCodeResponse(w, localization.ErrorIncompleteUserInfo.Code)
		return
	}

	err := h.PermissionService.CreatePermissionGroup(r.Context(), request)
	if err != nil {
		h.logger.Errorf("[CreatePermissionGroup] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	h.logger.Infof("[CreatePermissionGroup] request sent successfully by user: %s", userContext.UserID)
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupRequestCreated, nil)
}

// Get Permission Groups
//
//	@Summary		Get Permission Groups
//	@Description	Retrieves a paginated list of permission groups
//	@Tags			Permission
//	@Security		BearerAuth
//	@Produce		json
//	@Param			page		query		int	false	"Page number"
//	@Param			per_page	query		int	false	"Items per page"
//	@Success		200			{object}	localization.StandardResponse{data=permission.PaginatedPermissionGroupResponse}
//	@Failure		400,500		{object}	localization.StandardResponse{data=nil}
//	@Router			/permissions [get]
func (h *PermissionHandler) GetPermissionGroups(w http.ResponseWriter, r *http.Request) {
	filterparams := common_utils.ExtractFilterParams(r)
	permissionGroups, err := h.PermissionService.GetPermissionGroups(r.Context(), filterparams)
	if err != nil {
		h.logger.Errorf("[GetPermissionGroups] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupsFetched, permissionGroups)
}

// Get Permission Group by Name
//
//	@Summary		Get Permission Group by ID
//	@Description	Retrieves a permission group by its ID
//	@Tags			Permission
//	@Security		BearerAuth
//	@Produce		json
//	@Param			ID			path		string	true	"ID"
//	@Success		200			{object}	localization.StandardResponse{data=permission.PermissionCategoryResponse}
//	@Failure		400,404,500	{object}	localization.StandardResponse{data=nil}
//	@Router			/permissions/{group_name} [get]
func (h *PermissionHandler) GetPermissionGroup(w http.ResponseWriter, r *http.Request) {
	permissionGroup, err := h.PermissionService.GetPermissionGroup(chi.URLParam(r, "group_name"))
	if err != nil {
		h.logger.Errorf("[GetPermissionGroup] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupFetched, permissionGroup)
}

func (h *PermissionHandler) GetPermissionGroupById(w http.ResponseWriter, r *http.Request) {
	permissionGroup, err := h.PermissionService.GetPermissionGroupById(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		h.logger.Errorf("[GetPermissionGroup] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}
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
	id := chi.URLParam(r, "id")
	if id == "" {
		localization.SendErrorResponse(w, localization.ErrorGroupNameRequired, nil, nil)
		return
	}

	var request permission.UpdatePermissionGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Errorf("[UpdatePermissionGroup] failed to decode request: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	request.Id = id

	if err := request.Validate(); err != nil {
		h.logger.Warnf("[UpdatePermissionGroup] validation failed: %v", err)
		localization.SendBadRequestResponse(w, err.Error())
		return
	}

	userContext := common_utils.ExtractUserContext(r)
	if common_utils.IsIncomplete(userContext) {
		localization.SendErrorResponse(w, localization.ErrorIncompleteUserInfo, nil, nil)
		return
	}

	err := h.PermissionService.UpdatePermissionGroup(r.Context(), request)
	if err != nil {
		h.logger.Errorf("[UpdatePermissionGroup] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

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
	categories, err := h.PermissionService.GetAllPermissionCategoriesWithPermissions(r.Context())
	if err != nil {
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

	localization.SendSuccessResponse(w, localization.SuccessPermissionCategoriesFetched, grouped)
}

func (h *PermissionHandler) GetPermissionCategoriesByDepartment(w http.ResponseWriter, r *http.Request) {
	departmentID := chi.URLParam(r, "department_id")
	if departmentID == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorDepartmentIDRequired.Code)
		return
	}
	permissionCategory, err := h.PermissionService.GetPermissionCategoriesByDepartment(r.Context(), departmentID)
	if err != nil {
		h.logger.Errorf("[GetPermissionCategoriesByDepartment] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessPermissionCategoriesFetched, permissionCategory)
}

func (h *PermissionHandler) GetPermissionGroupsByDepartment(w http.ResponseWriter, r *http.Request) {
	departmentID := chi.URLParam(r, "department_id")
	if departmentID == "" {
		localization.SendErrorByCodeResponse(w, localization.ErrorDepartmentIDRequired.Code)
		return
	}
	filterParam := common_utils.ExtractFilterParams(r)

	permissionGroups, err := h.PermissionService.GetPermissionGroupsByDepartment(r.Context(), departmentID, filterParam)
	if err != nil {
		h.logger.Errorf("[GetPermissionGroupsByDepartment] service error: %v", err)
		localization.SendErrorByCodeResponse(w, err.Error())
		return
	}

	localization.SendSuccessResponse(w, localization.SuccessPermissionGroupsFetched, permissionGroups)
}
