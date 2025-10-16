package cpsuser

import (
	"context"
	"errors"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/cps_user/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type cpsUserService struct {
	repo              storage.CpsUserRepository
	permissionService service.PermissionService
	departmentRepo    storage.DepartmentRepository
	logger            shared_utils.Logger
	cpsService        service.CPSActionService
}

func NewCPSUserService(repo storage.CpsUserRepository, departmentRepo storage.DepartmentRepository, permission service.PermissionService, cps service.CPSActionService, logger shared_utils.Logger) service.CPSUserService {
	return &cpsUserService{
		repo:              repo,
		permissionService: permission,
		cpsService:        cps,
		logger:            logger,
		departmentRepo:    departmentRepo,
	}
}

func (s *cpsUserService) CreateUserRequest(ctx context.Context, req cpsuser.CreateUserRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)

	normalized := local_util.FormatPhoneNumber(req.PhoneNumber)
	req.PhoneNumber =normalized
	exists, err := core.UsernameExists(ctx, s.repo, req.UserName)
	if err != nil {
		return err
	}
	if exists {
		return errors.New(localization.ErrorUserAlreadyExists.Code)
	}
	emailCheck, err := core.EmailExists(ctx, s.repo, req.Email)
	if err != nil {
		return err
	}
	if emailCheck {
		return errors.New(localization.ErrorExistEmail.Code)
	}
	phoneCheck, err := core.PhoneNumberExists(ctx, s.repo, req.PhoneNumber)
	if err != nil {
		return err
	}
	if phoneCheck {
		return errors.New(localization.ErrorExistPhoneNumber.Code)
	}

	// department validation
	if req.Department.IsZero() {
		return errors.New(localization.ErrorInvalidRequest.Code)
	}
	dep, err := s.departmentRepo.FindByID(ctx, req.Department.Hex())
	if err != nil {
		return err
	}
	if !dep.Enabled || dep.IsDeleted {
		return errors.New(localization.ErrorDepartmentNotFound.Code)
	}

	if len(req.PermissionCategory) > 0 {
		categoryIDs := make([]string, len(req.PermissionCategory))
		for i, id := range req.PermissionCategory {
			categoryIDs[i] = id.Hex()
		}
		if _, err := s.permissionService.ValidatePermissionCategories(ctx, categoryIDs); err != nil {
			return err
		}
	}

	if len(req.PermissionGroups) > 0 {
		groupIDs := make([]string, len(req.PermissionGroups))
		for i, id := range req.PermissionGroups {
			groupIDs[i] = id.Hex()
		}
		if _, err := s.permissionService.ValidatePermissionGroups(ctx, groupIDs); err != nil {
			return err
		}
	}

	cpsUser := core.CPSUModel(req)

	cpsActionModel := lib.CpsModelBuilder("", makerData, nil, cpsUser, string(constants.RequestCpsUserCreate), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		return err
	}
	return nil
}

func (s *cpsUserService) UpdateUserRequest(ctx context.Context, req cpsuser.UpdateUserRequest) error {
	if req.PhoneNumber!= ""{
	normalized := local_util.FormatPhoneNumber(req.PhoneNumber)
	req.PhoneNumber =normalized

	phoneCheck, err := core.PhoneNumberExists(ctx, s.repo, req.PhoneNumber)
	if err != nil {
		return err
	}
	if phoneCheck {
		return errors.New(localization.ErrorExistPhoneNumber.Code)
	}
	}
	if req.UserName != "" {

		currentUser, err := s.repo.FindByID(ctx, req.UserCode)
		if err != nil {
			return err
		}

		if currentUser.UserName != req.UserName {
			exists, err := core.UsernameExists(ctx, s.repo, req.UserName)
			if err != nil {
				return err
			}
			if exists {
				return errors.New(localization.ErrorUserAlreadyExists.Code)
			}
		}
	}
	if req.Email != ""{
	emailCheck, err := core.EmailExists(ctx, s.repo, req.Email)
	if err != nil {
		return err
	}
	if emailCheck {
		return errors.New(localization.ErrorExistEmail.Code)
	}
	
	}
	

	// department validation - only if department is being updated
	if !req.Department.IsZero() {
		dep, err := s.departmentRepo.FindByID(ctx, req.Department.Hex())
		if err != nil {
			return err
		}
		if !dep.Enabled || dep.IsDeleted {
			return errors.New(localization.ErrorDepartmentNotFound.Code)
		}
	}

	if len(req.PermissionCategory) > 0 {
		categoryIDs := make([]string, len(req.PermissionCategory))
		for i, id := range req.PermissionCategory {
			categoryIDs[i] = id.Hex()
		}
		if _, err := s.permissionService.ValidatePermissionCategories(ctx, categoryIDs); err != nil {
			return err
		}
	}

	if len(req.PermissionGroups) > 0 {
		groupIDs := make([]string, len(req.PermissionGroups))
		for i, id := range req.PermissionGroups {
			groupIDs[i] = id.Hex()
		}
		if _, err := s.permissionService.ValidatePermissionGroups(ctx, groupIDs); err != nil {
			return err
		}
	}

	makerData := local_util.ExtractUserFromContext(ctx)
	cpsActionModel := lib.CpsModelBuilder(
		req.UserCode,
		makerData,
		nil,
		req,
		string(constants.RequestCpsUserUpdate),
		constants.UPDATE,
	)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		return err
	}
	return nil
}

func (s *cpsUserService) DeleteUserRequest(ctx context.Context, userCode string) error {
	if userCode == "" {
		return errors.New(localization.ErrorUserCodeRequired.Code)
	}

	existing, err := s.repo.FindByID(ctx, userCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		return err
	}

	maker := local_util.ExtractUserFromContext(ctx)
	payload := model.CPSUser{
		UserCode:  userCode,
		IsDeleted: true,
	}

	cpsAction := lib.CpsModelBuilder(userCode, maker, existing, payload, string(constants.RequestCpsUserDelete), constants.DELETE)
	return s.cpsService.CreateCPSAction(ctx, &cpsAction)
}

func (s *cpsUserService) EnableUser(ctx context.Context, userCode string) error {
	if userCode == "" {
		return errors.New(localization.ErrorUserCodeRequired.Code)
	}

	prev, err := s.repo.FindByID(ctx, userCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		return err
	}

	if prev.Enabled {
		return errors.New(localization.ErrorUserAlreadyEnabled.Code)
	}

	updated := *prev
	updated.Enabled = true
	updated.PasswordDisable = false

	maker := local_util.ExtractUserFromContext(ctx)
	cpsAction := lib.CpsModelBuilder(userCode, maker, prev, updated, string(constants.RequestCpsUserEnable), constants.UPDATE)
	return s.cpsService.CreateCPSAction(ctx, &cpsAction)
}

func (s *cpsUserService) DisableUser(ctx context.Context, userCode string) error {
	if userCode == "" {
		return errors.New(localization.ErrorUserCodeRequired.Code)
	}

	prev, err := s.repo.FindByID(ctx, userCode)
	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New(localization.ErrorResourceNotFound.Code)
		}
		return err
	}

	if !prev.Enabled {
		return errors.New(localization.ErrorUserAlreadyDisabled.Code)
	}

	updated := prev
	updated.Enabled = false
	updated.PasswordDisable = true

	maker := local_util.ExtractUserFromContext(ctx)

	cpsaction := lib.CpsModelBuilder(userCode, maker, prev, &updated, string(constants.RequestCpsUserDisable), constants.UPDATE)
	return s.cpsService.CreateCPSAction(ctx, &cpsaction)
}

func (s *cpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*cpsuser.CPSUserDTO, error) {
	if userCode == "" {
		return nil, errors.New(localization.ErrorUserCodeRequired.Code)
	}

	user, err := s.repo.FindByID(ctx, userCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		return nil, err
	}

	return core.ConvertToDTO(user), nil
}

func (s *cpsUserService) GetPopulatedCpsUser(ctx context.Context, userCode string) (*cpsuser.CpsUserResponse, error) {
	if userCode == "" {
		return nil, errors.New(localization.ErrorUserCodeRequired.Code)
	}

	user, err := s.repo.GetPopulatedByID(ctx, userCode)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) || err.Error() == localization.ErrorResourceNotFound.Code {
			return nil, errors.New(localization.ErrorResourceNotFound.Code)
		}
		return nil, err
	}

	return user, nil
}
func (s *cpsUserService) GetAllCPSUsers(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*cpsuser.CPSUserDTO], error) {
	if filter == nil {
		f := types.Filter{}
		filter = &f
	}

	users, err := s.repo.FindAllWithPagination(ctx, *filter)
	if err != nil {
		return nil, err
	}

	dtos := make([]*cpsuser.CPSUserDTO, len(users.Data))
	for i, u := range users.Data {
		dtos[i] = core.ConvertToDTO(u)
	}

	return &types.PaginatedResponse[[]*cpsuser.CPSUserDTO]{
		Data: dtos,
		Meta: users.Meta,
	}, nil
}

func (s *cpsUserService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	action.MakerActionTime = time.Now()
	action.LastModifiedAt = action.MakerActionTime

	switch action.RequestAction {
	case string(constants.RequestCpsUserCreate):

		cur, err := core.BindCPSUserFromAction(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Create(ctx, &cur); err != nil {
			return nil, err
		}
		return action, nil

	case string(constants.RequestCpsUserUpdate):
		cur, err := core.BindCPSUserUpdateFromAction(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		update := core.CPSUUpdateModel(cur)
		if err := s.repo.Update(ctx, action.UniqueId, update); err != nil {
			return nil, err
		}
		return action, nil

	case string(constants.RequestCpsUserDelete):
		_, err := core.BindCPSUserFromAction(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Delete(ctx, action.UniqueId); err != nil {
			return nil, err
		}
		return action, nil

	case string(constants.RequestCpsUserEnable):
		_, err := core.BindCPSUserFromAction(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, action.UniqueId, true); err != nil {
			return nil, err
		}
		return action, nil

	case string(constants.RequestCpsUserDisable):
		_, err := core.BindCPSUserFromAction(action.CurrentAction)
		if err != nil {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, action.UniqueId, false); err != nil {
			return nil, err
		}
		return action, nil
	}

	return nil, errors.New("UNHANDLED_SERVER_ERROR")
}
