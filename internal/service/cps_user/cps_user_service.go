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
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/service/cps_user/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type cpsUserService struct {
	repo              storage.CpsUserRepository
	permissionService service.PermissionService
	departmentRepo    storage.AppAccessListRepository
	logger            shared_utils.Logger
	cpsService        service.CPSActionService
}

func NewCPSUserService(repo storage.CpsUserRepository, departmentRepo storage.AppAccessListRepository, permission service.PermissionService, cps service.CPSActionService, logger shared_utils.Logger) service.CPSUserService {
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
	// department validation
	if req.Department.IsZero() {
		return errors.New(localization.ErrorInvalidRequest.Code)
	}
	if _, err := s.departmentRepo.FindByID(ctx, req.Department.Hex()); err != nil {
		return err
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

	// department validation
	if req.Department.IsZero() {
		return errors.New(localization.ErrorInvalidRequest.Code)
	}
	if _, err := s.departmentRepo.FindByID(ctx, req.Department.Hex()); err != nil {
		return err
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

	// role validation
	if req.Role != "" && req.Role != "maker" && req.Role != "checker" {
		return errors.New("MAKER_OR_CHECKER")
	}

	
	if req.PhoneNumber != "" {
		normalized := local_util.FormatPhoneNumber(req.PhoneNumber)
		if normalized == "" {
			return errors.New("UNSUPPORTED_PHONE_NUMBER_FORMAT")
		}
		req.PhoneNumber = normalized
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

	maker := local_util.ExtractUserFromContext(ctx)
	cpsAction := lib.CpsModelBuilder(userCode, maker, prev, &updated, string(constants.RequestCpsUserEnable), constants.UPDATE)
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

	updated := *prev
	updated.Enabled = false

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

	switch action.ActionType {
	case string(cpsaction.ActionCreate):
		cur, ok := action.CurrentAction.(model.CPSUser)
		if !ok {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Create(ctx, &cur); err != nil {
			return nil, err
		}
		return action, nil

	case string(cpsaction.ActionUpdate):
		upd, ok := action.CurrentAction.(cpsuser.UpdateUserRequest)
		if !ok {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		update := core.CPSUUpdateModel(upd)
		if err := s.repo.Update(ctx, upd.UserCode, update); err != nil {
			return nil, err
		}
		return action, nil

	case string(cpsaction.ActionDelete):
		cur, ok := action.CurrentAction.(model.CPSUser)
		if !ok {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Delete(ctx, cur.UserCode); err != nil {
			return nil, err
		}
		return action, nil

	case string(cpsaction.ActionEnable):
		cur, ok := action.CurrentAction.(model.CPSUser)
		if !ok {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, cur.UserCode, true); err != nil {
			return nil, err
		}
		return action, nil

	case string(cpsaction.ActionDisable):
		cur, ok := action.CurrentAction.(model.CPSUser)
		if !ok {
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, cur.UserCode, false); err != nil {
			return nil, err
		}
		return action, nil
	}

	return nil, errors.New("UNHANDLED_SERVER_ERROR")
}
