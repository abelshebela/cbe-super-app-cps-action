package cpsuser

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/cps_user/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type cpsUserService struct {
	repo              storage.CpsUserRepository
	roleRepo          storage.RoleRepository
	approverRepo      storage.CPSActionApproveIndexRepository
	permissionService service.PermissionService
	departmentRepo    storage.DepartmentRepository
	logger            shared_utils.Logger
	cpsService        service.CPSActionService
}

func NewCPSUserService(repo storage.CpsUserRepository, roleRepo storage.RoleRepository, approverRepo storage.CPSActionApproveIndexRepository, departmentRepo storage.DepartmentRepository, permission service.PermissionService, cps service.CPSActionService, logger shared_utils.Logger) service.CPSUserService {
	return &cpsUserService{
		repo:              repo,
		roleRepo:          roleRepo,
		approverRepo:      approverRepo,
		permissionService: permission,
		cpsService:        cps,
		logger:            logger,
		departmentRepo:    departmentRepo,
	}
}

func (s *cpsUserService) CreateUserRequest(ctx context.Context, req cpsuser.CreateUserRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateUserRequest", "CPSUser", "CreateUserRequest")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)

	normalized := local_util.FormatPhoneNumber(req.PhoneNumber)
	req.PhoneNumber = normalized
	exists, err := core.UsernameExists(ctx, s.repo, req.UserName)
	if err != nil {
		span.AddEvent("failed to check username existence", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	if exists {
		span.AddEvent("username already exists", trace.WithAttributes(attribute.String("username", req.UserName)))
		return errors.New(localization.ErrorUserAlreadyExists.Code)
	}
	emailCheck, err := core.EmailExists(ctx, "", s.repo, req.Email)
	if err != nil {
		span.AddEvent("failed to check email existence", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	if emailCheck {
		span.AddEvent("email already exists", trace.WithAttributes(attribute.String("email", req.Email)))
		return errors.New(localization.ErrorExistEmail.Code)
	}
	phoneCheck, err := core.PhoneNumberExists(ctx, "", s.repo, req.PhoneNumber)
	if err != nil {
		span.AddEvent("failed to check phone number existence", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	if phoneCheck {
		span.AddEvent("phone number already exists", trace.WithAttributes(attribute.String("phone_number", req.PhoneNumber)))
		return errors.New(localization.ErrorExistPhoneNumber.Code)
	}

	// department validation
	if req.Department.IsZero() {
		span.AddEvent("department is zero", trace.WithAttributes(attribute.String("error", "department is zero")))
		return errors.New(localization.ErrorInvalidRequest.Code)
	}
	dep, err := s.departmentRepo.FindByID(ctx, req.Department.Hex())
	if err != nil {
		span.AddEvent("failed to find department by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	if !dep.Enabled || dep.IsDeleted {
		span.AddEvent("department not found", trace.WithAttributes(attribute.String("department_id", req.Department.Hex())))
		return errors.New(localization.ErrorDepartmentNotFound.Code)
	}

	// Populate permission categories if provided
	// var populatedCategories []cpsuser.PermissionCategoryResponse
	// if len(req.PermissionCategory) > 0 {
	// 	populated, err := s.permissionService.GetPopulatedPermissionCategories(ctx, req.PermissionCategory)
	// 	if err != nil {
	// 		span.AddEvent("failed to populate permission categories", trace.WithAttributes(attribute.String("error", err.Error())))
	// 		s.logger.Errorf("[CreateUserRequest] failed to populate permission categories: %v", err)
	// 		return err
	// 	}
	// 	populatedCategories = populated
	// }
	// var populatedGroups []cpsuser.PermissionGroupResponse
	// if len(req.PermissionGroups) > 0 {
	// 	populated, err := s.permissionService.GetPopulatedPermissionGroups(ctx, req.PermissionGroups)
	// 	if err != nil {
	// 		span.AddEvent("failed to populate permission groups", trace.WithAttributes(attribute.String("error", err.Error())))
	// 		s.logger.Errorf("[CreateUserRequest] failed to populate permission groups: %v", err)
	// 		return err
	// 	}
	// 	populatedGroups = populated
	// }

	cpsUser := core.CPSUModel(req)
	cpsUser.JobTitle = req.JobTitle
	cpsActionModel := lib.CpsModelBuilder("", makerData, nil, cpsUser, string(constants.RequestCpsUserCreate), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("failed to create CPS action", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("[CreateUserRequest] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[CreateUserRequest] CPS user creation request created successfully")
	return nil
}

func (s *cpsUserService) UpdateUserRequest(ctx context.Context, usercode string, req cpsuser.UpdateUserRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateUserRequest", "CPSUser", "UpdateUserRequest")
	defer span.End()

	currentUser, err := s.repo.FindByID(ctx, usercode)
	if err != nil {
		span.AddEvent("failed to find user by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if req.PhoneNumber != "" {
		normalized := local_util.FormatPhoneNumber(req.PhoneNumber)
		req.PhoneNumber = normalized

		phoneCheck, err := core.PhoneNumberExists(ctx, currentUser.UserCode, s.repo, req.PhoneNumber)
		if err != nil {
			span.AddEvent("failed to check phone number existence", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
		if phoneCheck {
			span.AddEvent("phone number already exists", trace.WithAttributes(attribute.String("phone_number", req.PhoneNumber)))
			return errors.New(localization.ErrorExistPhoneNumber.Code)
		}
	}
	if req.UserName != "" {
		if currentUser.UserName != req.UserName {
			exists, err := core.UsernameExists(ctx, s.repo, req.UserName)
			if err != nil {
				span.AddEvent("failed to check username existence", trace.WithAttributes(attribute.String("error", err.Error())))
				return err
			}
			if exists {
				span.AddEvent("username already exists", trace.WithAttributes(attribute.String("username", req.UserName)))
				return errors.New(localization.ErrorUserAlreadyExists.Code)
			}
		}
	}
	if req.Email != "" {
		emailCheck, err := core.EmailExists(ctx, currentUser.UserCode, s.repo, req.Email)
		if err != nil {
			span.AddEvent("failed to check email existence", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
		if emailCheck {
			span.AddEvent("email already exists", trace.WithAttributes(attribute.String("email", req.Email)))
			return errors.New(localization.ErrorExistEmail.Code)
		}

	}

	// var dep *model.Department
	// if !req.Department.IsZero() {
	// 	d, err := s.departmentRepo.FindByID(ctx, req.Department.Hex())
	// 	if err != nil {
	// 		span.AddEvent("failed to find department by id", trace.WithAttributes(attribute.String("error", err.Error())))
	// 		return err
	// 	}
	// 	if !d.Enabled || d.IsDeleted {
	// 		span.AddEvent("department not found", trace.WithAttributes(attribute.String("department_id", req.Department.Hex())))
	// 		return errors.New(localization.ErrorDepartmentNotFound.Code)
	// 	}
	// 	dep = d
	// }

	// var populatedCategories []cpsuser.PermissionCategoryResponse
	// if len(req.PermissionCategory) > 0 {
	// 	populated, err := s.permissionService.GetPopulatedPermissionCategories(ctx, req.PermissionCategory)
	// 	if err != nil {
	// 		span.AddEvent("failed to populate permission categories", trace.WithAttributes(attribute.String("error", err.Error())))
	// 		return err
	// 	}
	// 	populatedCategories = populated
	// }

	// var populatedGroups []cpsuser.PermissionGroupResponse
	// if len(req.PermissionGroups) > 0 {
	// 	populated, err := s.permissionService.GetPopulatedPermissionGroups(ctx, req.PermissionGroups)
	// 	if err != nil {
	// 		span.AddEvent("failed to populate permission groups", trace.WithAttributes(attribute.String("error", err.Error())))
	// 		return err
	// 	}
	// 	populatedGroups = populated
	// }

	// payload := map[string]interface{}{
	// 	"user": req,
	// }
	// if len(populatedCategories) > 0 {
	// 	payload["permission_categories"] = populatedCategories
	// }
	// if len(populatedGroups) > 0 {
	// 	payload["permission_groups"] = populatedGroups
	// }
	// if dep != nil {
	// 	payload["portal_cards"] = dep.PortalCards
	// }

	updated := cpsuser.UpdateUserRequest{
		UserName:    req.UserName,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Gender:      req.Gender,
		Email:       req.Email,
		JobTitle:    req.JobTitle,
	}

	makerData := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(
		usercode,
		makerData,
		currentUser,
		updated,
		string(constants.RequestCpsUserUpdate),
		constants.UPDATE,
	)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("failed to create CPS action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}

func (s *cpsUserService) DeleteUserRequest(ctx context.Context, userCode string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteUserRequest", "CPSUser", "DeleteUserRequest")
	defer span.End()

	if userCode == "" {
		span.AddEvent("user code is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return errors.New(localization.ErrorUserCodeRequired.Code)
	}

	existing, err := s.repo.FindByID(ctx, userCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			span.AddEvent("user not found", trace.WithAttributes(attribute.String("user_code", userCode)))
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		span.AddEvent("failed to find user by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	maker := local_util.ExtractUserFromContext(ctx)
	payload := model.CPSUser{
		UserCode:  userCode,
		IsDeleted: true,
	}

	cpsAction := lib.CpsModelBuilder(userCode, maker, existing, payload, string(constants.RequestCpsUserDelete), constants.DELETE)
	err = s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		span.AddEvent("failed to create cps action", trace.WithAttributes(attribute.String("error", err.Error())))
	}
	return err
}

func (s *cpsUserService) EnableUser(ctx context.Context, userCode string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableUser", "CPSUser", "EnableUser")
	defer span.End()

	if userCode == "" {
		span.AddEvent("user code is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return errors.New(localization.ErrorUserCodeRequired.Code)
	}
	prev, err := s.repo.FindByID(ctx, userCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			span.AddEvent("user not found", trace.WithAttributes(attribute.String("user_code", userCode)))
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		span.AddEvent("failed to find user by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if prev.Enabled {
		span.AddEvent("user already enabled", trace.WithAttributes(attribute.String("user_code", userCode)))
		return errors.New(localization.ErrorUserAlreadyEnabled.Code)
	}

	updated := *prev
	updated.Enabled = true
	updated.LoginAttemptCount = 0
	updated.PasswordDisable = false

	maker := local_util.ExtractUserFromContext(ctx)
	cpsAction := lib.CpsModelBuilder(userCode, maker, prev, updated, string(constants.RequestCpsUserEnable), constants.UPDATE)
	err = s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		span.AddEvent("failed to create cps action", trace.WithAttributes(attribute.String("error", err.Error())))
	}
	return err
}

func (s *cpsUserService) DisableUser(ctx context.Context, userCode string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableUser", "CPSUser", "DisableUser")
	defer span.End()

	if userCode == "" {
		span.AddEvent("user code is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return errors.New(localization.ErrorUserCodeRequired.Code)
	}

	prev, err := s.repo.FindByID(ctx, userCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			span.AddEvent("user not found", trace.WithAttributes(attribute.String("user_code", userCode)))
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		span.AddEvent("failed to find user by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if !prev.Enabled {
		span.AddEvent("user already disabled", trace.WithAttributes(attribute.String("user_code", userCode)))
		return errors.New(localization.ErrorUserAlreadyDisabled.Code)
	}

	updated := prev
	updated.Enabled = false
	updated.PasswordDisable = true

	maker := local_util.ExtractUserFromContext(ctx)

	cpsaction := lib.CpsModelBuilder(userCode, maker, prev, &updated, string(constants.RequestCpsUserDisable), constants.UPDATE)
	err = s.cpsService.CreateCPSAction(ctx, &cpsaction)
	if err != nil {
		span.AddEvent("failed to create cps action", trace.WithAttributes(attribute.String("error", err.Error())))
	}
	return err
}

func (s *cpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*cpsuser.CPSUserDTO, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchUserByUserCode", "CPSUser", "FetchUserByUserCode")
	defer span.End()

	if userCode == "" {
		span.AddEvent("user code is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return nil, errors.New(localization.ErrorUserCodeRequired.Code)
	}

	user, err := s.repo.FindByID(ctx, userCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			span.AddEvent("user not found", trace.WithAttributes(attribute.String("user_code", userCode)))
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		span.AddEvent("failed to find user by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}

	roles, err := s.roleRepo.FindByName(ctx, user.JobTitle)
	if err != nil {
		span.AddEvent("failed to find role by name", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}

	makerAlloc, checkerAlloc, auditorAlloc, err := s.approverRepo.PopulateUserApproverAllocations(ctx, roles.ID.Hex())
	if err != nil {
		span.AddEvent("failed to populate user approver allocations", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}

	return core.ConvertToDTO(user, makerAlloc, checkerAlloc, auditorAlloc), nil
}

func (s *cpsUserService) GetPopulatedCpsUser(ctx context.Context, userCode string) (*cpsuser.CpsUserResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPopulatedCpsUser", "CPSUser", "GetPopulatedCpsUser")
	defer span.End()

	if userCode == "" {
		span.AddEvent("user code is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return nil, errors.New(localization.ErrorUserCodeRequired.Code)
	}

	user, err := s.repo.GetPopulatedByID(ctx, userCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			span.AddEvent("user not found", trace.WithAttributes(attribute.String("user_code", userCode)))
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		span.AddEvent("failed to get populated by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err

	}

	return user, nil
}

func (s *cpsUserService) GetCpsUserDetail(ctx context.Context, userCode string) (*cpsuser.CpsUserResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCpsUserDetail", "CPSUser", "GetCpsUserDetail")
	defer span.End()

	if userCode == "" {
		span.AddEvent("user code is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return nil, errors.New(localization.ErrorUserCodeRequired.Code)
	}

	populated, err := s.repo.GetPopulatedByID(ctx, userCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			span.AddEvent("user not found", trace.WithAttributes(attribute.String("user_code", userCode)))
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		span.AddEvent("failed to get populated by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}

	// detail := cpsuser.BuildCpsUserDetail(populated)
	return populated, nil
}

func (s *cpsUserService) GetAllCPSUsers(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*cpsuser.CPSUserWithDepartment], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllCPSUsers", "CPSUser", "GetAllCPSUsers")
	defer span.End()

	if filter == nil {
		f := types.Filter{}
		filter = &f
	}

	users, err := s.repo.FindAllWithPagination(ctx, *filter)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}

	return users, nil
}

func (s *cpsUserService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "CPSUser", "Authorize")
	defer span.End()

	action.MakerActionTime = time.Now()
	action.LastModifiedAt = action.MakerActionTime

	switch action.RequestAction {
	case string(constants.RequestCpsUserCreate):

		cur, err := local_util.JsonUnmarshal[model.CPSUser](action.CurrentAction)
		// cur, err := core.BindCPSUserFromAction(action.CurrentAction)
		if err != nil {
			span.AddEvent("failed to bind cps user from action", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Create(ctx, cur); err != nil {
			span.AddEvent("failed to create user", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		return action, nil

	case string(constants.RequestCpsUserUpdate):
		cur, err := local_util.JsonUnmarshal[model.CPSUser](action.CurrentAction)
		// cur, err := core.BindCPSUserUpdateFromAction(action.CurrentAction)
		if err != nil {
			span.AddEvent("failed to bind cps user update from action", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if err := s.repo.Update(ctx, action.UniqueId, cur); err != nil {
			span.AddEvent("failed to update user", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		return action, nil

	case string(constants.RequestCpsUserDelete):
		_, err := core.BindCPSUserFromAction(action.CurrentAction)
		if err != nil {
			span.AddEvent("failed to bind cps user from action", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Delete(ctx, action.UniqueId); err != nil {
			span.AddEvent("failed to delete user", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		return action, nil

	case string(constants.RequestCpsUserEnable):
		_, err := core.BindCPSUserFromAction(action.CurrentAction)
		if err != nil {
			span.AddEvent("failed to bind cps user from action", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, action.UniqueId, true); err != nil {
			span.AddEvent("failed to enable user", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		return action, nil

	case string(constants.RequestCpsUserDisable):
		_, err := core.BindCPSUserFromAction(action.CurrentAction)
		if err != nil {
			span.AddEvent("failed to bind cps user from action", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, action.UniqueId, false); err != nil {
			span.AddEvent("failed to disable user", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		return action, nil
	}

	span.AddEvent("unhandled server error", trace.WithAttributes(attribute.String("error", "unhandled server error")))
	return nil, errors.New("UNHANDLED_SERVER_ERROR")
}
