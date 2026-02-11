package roles

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/roles/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"errors"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	sharedmodel "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RoleService struct {
	cpsService     service.CPSActionService
	portalCardRepo storage.PortalCardRepository
	roleRepository storage.JobRoleRepository
	cfg            config.VaultConfig
	logger         utils.Logger
}

func NewRoleService(roleRepo storage.JobRoleRepository, portalCard storage.PortalCardRepository, cpsService service.CPSActionService, cfg config.VaultConfig, logger utils.Logger) service.RoleService {
	return &RoleService{
		cpsService:     cpsService,
		portalCardRepo: portalCard,
		roleRepository: roleRepo,
		cfg:            cfg,
		logger:         logger,
	}
}

func (j *RoleService) Create(ctx context.Context, role imodel.JobRole) error {
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		j.logger.Errorf("[Role Service][Create] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	if err := core.RoleExistenChecker(ctx, constants.CREATE, "", role, j.roleRepository); err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New(err.Error())
		}
	}
	cpsModel := lib.CpsModelBuilder(constants.Empty, maker, nil, role, constants.RequestCreateRole, constants.CREATE)
	return j.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (j *RoleService) Update(ctx context.Context, id string, update imodel.JobRole) error {
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		j.logger.Errorf("[Role Service][Update] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	if err := core.RoleExistenChecker(ctx, constants.CREATE, "", update, j.roleRepository); err != nil {
		if err != mongo.ErrNoDocuments {
			return errors.New(localization.ErrorUnexpectedError.Code)
		}
	}

	prev, err := j.roleRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	newRole := *prev
	if update.Name != "" {
		newRole.Name = update.Name
	}
	if update.Code != "" {
		newRole.Code = update.Code
	}
	if len(update.PortalCards) > 0 {
		newRole.PortalCards = update.PortalCards
	}
	newRole.UpdatedAt = time.Now()

	cpsModel := lib.CpsModelBuilder(id, maker, prev, newRole, constants.RequestUpdateRole, constants.UPDATE)
	return j.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (j *RoleService) FindById(ctx context.Context, id string) (*imodel.JobRole, error) {
	return j.roleRepository.FindByID(ctx, id)
}

func (j *RoleService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.JobRole], error) {
	return j.roleRepository.FindAllWithPagination(ctx, filterParam)
}

func (j *RoleService) Authorize(ctx context.Context, cpsAction *sharedmodel.CPSAction) (*sharedmodel.CPSAction, error) {
	// Turn CurrentAction into Role, attach ID from UniqueId (if present), apply action
	var asAny any
	raw, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		j.logger.Errorf("[JobRole Service][Authorize] marshal CurrentAction failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if err := json.Unmarshal(raw, &asAny); err != nil {
		j.logger.Errorf("[JobRole Service][Authorize] unmarshal CurrentAction failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	raw2, _ := json.Marshal(asAny)
	var role imodel.JobRole
	if err := json.Unmarshal(raw2, &role); err != nil {
		j.logger.Errorf("[JobRole Service][Authorize] map to Role failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if cpsAction.UniqueId != "" {
		if oid, err := bson.ObjectIDFromHex(cpsAction.UniqueId); err == nil {
			role.ID = oid
		}
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateRole):
		role.CreatedAt = time.Now()
		if err := j.roleRepository.Create(ctx, &role); err != nil {
			return nil, err
		}
	case string(constants.RequestUpdateRole):
		role.UpdatedAt = time.Now()
		if err := j.roleRepository.Update(ctx, role.ID.Hex(), &role); err != nil {
			return nil, err
		}
	default:
		j.logger.Errorf("[JobRole Service][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	cpsAction.CurrentAction = role
	return cpsAction, nil
}
