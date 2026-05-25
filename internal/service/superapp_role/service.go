package superapprole

import (
	"context"
	"errors"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
	cps_roles_dto "cbe-super-app-cps-action/internal/constants/dto/cps_roles"
	superapproledto "cbe-super-app-cps-action/internal/constants/dto/superapp_role"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	imodel "cbe-super-app-cps-action/internal/constants/model"

	"github.com/hugokessem/coreio/core"
	cifLimit "github.com/hugokessem/coreio/lib/core/customer/customer_limit_fetch_by_cif"
	climit "github.com/hugokessem/coreio/lib/core/customer/customer_limit_fetch_by_service"
	shared_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type superAppRoleService struct {
	repo       storage.SuperAppRoleRepository
	cpsService service.CPSActionService
	core       core.CBECoreAPIInterface
	logger     utils.Logger
}

func NewSuperAppRoleService(repo storage.SuperAppRoleRepository, cpsService service.CPSActionService, coreInterface core.CBECoreAPIInterface, logger utils.Logger) *superAppRoleService {
	return &superAppRoleService{
		repo:       repo,
		cpsService: cpsService,
		core:       coreInterface,
		logger:     logger,
	}
}

func (s *superAppRoleService) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.SuperAppRoleGroup], error) {
	return s.repo.FindAllWithPagination(ctx, filterParam)
}

func (s *superAppRoleService) GetTransferLimitByRole(ctx context.Context, superappRole string, filterParam types.Filter) (*types.PaginatedResponse[[]cps_roles_dto.ServiceLevelLimitResponse], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	serviceCode := "GLOBAL-" + strings.ToUpper(superappRole)
	response, err := s.core.CustomerLimitFetchByService(core.CustomerLimitFetchByServiceParam{
		ServiceCode: serviceCode,
	})
	if err != nil {
		log.Errorf("[SuperAppRole][GetTransferLimitByRole] core call err: %v", err)
		return nil, err
	}

	if !response.Success {
		log.Warnf("[SuperAppRole][GetTransferLimitByRole] core returned failure for role: %s", superappRole)
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	results := buildLimitResults(response)
	return paginate(results, filterParam), nil
}

func buildLimitResults(response *core.CustomerLimitFetchByServiceResult) []cps_roles_dto.ServiceLevelLimitResponse {
	serviceMap := make(map[string]*cps_roles_dto.ServiceLevelLimitResponse)

	if response.Detail != nil && response.Detail.GChannelType != nil {
		for _, ch := range response.Detail.GChannelType.MChannelType {
			applyChannelToMap(serviceMap, ch)
		}
	}

	results := make([]cps_roles_dto.ServiceLevelLimitResponse, 0, len(serviceMap))
	for _, v := range serviceMap {
		results = append(results, *v)
	}
	return results
}

func applyChannelToMap(serviceMap map[string]*cps_roles_dto.ServiceLevelLimitResponse, ch climit.MChannelType) {
	if ch.SGServiceTypes == nil {
		return
	}
	for _, svc := range ch.SGServiceTypes.GServiceType {
		if _, exists := serviceMap[svc.Name]; !exists {
			serviceMap[svc.Name] = &cps_roles_dto.ServiceLevelLimitResponse{Name: svc.Name}
		}
		entry := serviceMap[svc.Name]
		switch ch.ChannelType {
		case "APP":
			entry.SuperAppMaxLimit = svc.CHANNELMAXLIMIT
			entry.SuperAppTranFreq = svc.CHANNELCOUNT
		case "USSD":
			entry.USSDMaxLimit = svc.CHANNELMAXLIMIT
			entry.USSDTranFreq = svc.CHANNELCOUNT
		}
	}
}

func paginate(results []cps_roles_dto.ServiceLevelLimitResponse, filterParam types.Filter) *types.PaginatedResponse[[]cps_roles_dto.ServiceLevelLimitResponse] {
	total := int64(len(results))
	page := filterParam.Page
	if page < 1 {
		page = 1
	}
	limit := filterParam.PerPage
	if limit < 1 {
		limit = 10
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	start := (page - 1) * limit
	if start > int(total) {
		start = int(total)
	}
	end := start + limit
	if end > int(total) {
		end = int(total)
	}

	meta := types.PaginationMeta{
		TotalDocs:     total,
		Limit:         limit,
		TotalPages:    totalPages,
		Page:          page,
		PagingCounter: start + 1,
		HasPrevPage:   page > 1,
		HasNextPage:   page < totalPages,
	}
	if meta.HasPrevPage {
		prev := page - 1
		meta.PrevPage = &prev
	}
	if meta.HasNextPage {
		next := page + 1
		meta.NextPage = &next
	}

	return &types.PaginatedResponse[[]cps_roles_dto.ServiceLevelLimitResponse]{
		Data: results[start:end],
		Meta: meta,
	}
}

func (s *superAppRoleService) GetGlobalLimitByRole(ctx context.Context, superappRole string) (*cps_roles_dto.GlobalLimitResponse, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	customerNumber := "GLOBAL-" + strings.ToUpper(superappRole)
	response, err := s.core.CustomerLimitFetchByCustomerNumber(core.CustomerLimitFetchByCIFParam{
		CustomerNumber: customerNumber,
	})
	if err != nil {
		log.Errorf("[SuperAppRole][GetGlobalLimitByRole] core call err: %v", err)
		return nil, err
	}
	if !response.Success {
		log.Warnf("[SuperAppRole][GetGlobalLimitByRole] core returned failure for role: %s", superappRole)
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	return buildGlobalLimitResponse(response.Detail), nil
}

func buildGlobalLimitResponse(detail *cifLimit.CustomerLimitType) *cps_roles_dto.GlobalLimitResponse {
	resp := &cps_roles_dto.GlobalLimitResponse{}
	if detail == nil || detail.GUserChannel == nil {
		return resp
	}
	for _, ch := range detail.GUserChannel.MUserChannel {
		channel := cps_roles_dto.GlobalLimitChannelResponse{ChannelType: ch.UserChannelType}
		if ch.SGServiceType != nil {
			for _, svc := range ch.SGServiceType.Services {
				channel.Services = append(channel.Services, cps_roles_dto.GlobalLimitServiceEntry{
					ServiceType: svc.Name,
					MaxAmount:   svc.ServiceMaxAmt,
					MaxCount:    svc.UserMaxCnt,
				})
			}
		}
		resp.Channels = append(resp.Channels, channel)
	}
	return resp
}

func (s *superAppRoleService) EnableByRole(ctx context.Context, superappRole string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	cpsActionData := lib.CpsModelBuilder(
		superappRole,
		makerUser,
		map[string]interface{}{"superapp_role": superappRole},
		map[string]interface{}{"superapp_role": superappRole, "is_enabled": true},
		string(constants.RequestEnableSuperAppRole),
		constants.UPDATE,
	)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[SuperAppRole][EnableByRole] cps action err: %v", err)
		return err
	}
	return nil
}

func (s *superAppRoleService) DisableByRole(ctx context.Context, superappRole string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	cpsActionData := lib.CpsModelBuilder(
		superappRole,
		makerUser,
		map[string]interface{}{"superapp_role": superappRole, "is_enabled": true},
		map[string]interface{}{"superapp_role": superappRole, "is_enabled": false},
		string(constants.RequestDisableSuperAppRole),
		constants.UPDATE,
	)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[SuperAppRole][DisableByRole] cps action err: %v", err)
		return err
	}
	return nil
}

func (s *superAppRoleService) DeleteByRole(ctx context.Context, superappRole string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	cpsActionData := lib.CpsModelBuilder(
		superappRole,
		makerUser,
		map[string]interface{}{"superapp_role": superappRole},
		map[string]interface{}{"superapp_role": superappRole},
		string(constants.RequestDeleteSuperAppRole),
		constants.DELETE,
	)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[SuperAppRole][DeleteByRole] cps action err: %v", err)
		return err
	}
	return nil
}

func (s *superAppRoleService) GetAccessListsByRole(ctx context.Context, superappRole string) ([]imodel.APPAccessList, []imodel.APPAccessList, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	exists, err := s.repo.RoleExists(ctx, superappRole)
	if err != nil {
		log.Errorf("[SuperAppRole][GetAccessListsByRole] role check err: %v", err)
		return nil, nil, err
	}
	if !exists {
		return nil, nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	globallyEnabled, err := s.repo.FindGloballyEnabledAccessLists(ctx)
	if err != nil {
		log.Errorf("[SuperAppRole][GetAccessListsByRole] global enabled fetch err: %v", err)
		return nil, nil, err
	}

	globallyDisabled, err := s.repo.FindGloballyDisabledAccessLists(ctx)
	if err != nil {
		log.Errorf("[SuperAppRole][GetAccessListsByRole] global disabled fetch err: %v", err)
		return nil, nil, err
	}

	roleBlocked, err := s.repo.FindRoleBlockedAccessLists(ctx, superappRole)
	if err != nil {
		log.Errorf("[SuperAppRole][GetAccessListsByRole] role blocked fetch err: %v", err)
		return nil, nil, err
	}

	// disabled = distinct union of role-blocked + globally disabled
	disabledMap := make(map[string]imodel.APPAccessList, len(roleBlocked)+len(globallyDisabled))
	for _, al := range roleBlocked {
		disabledMap[al.ID] = al
	}
	for _, al := range globallyDisabled {
		disabledMap[al.ID] = al
	}
	disabled := make([]imodel.APPAccessList, 0, len(disabledMap))
	for _, al := range disabledMap {
		disabled = append(disabled, al)
	}

	// enabled = globally enabled - role-blocked
	blockedSet := make(map[string]struct{}, len(roleBlocked))
	for _, al := range roleBlocked {
		blockedSet[al.ID] = struct{}{}
	}
	enabled := make([]imodel.APPAccessList, 0, len(globallyEnabled))
	for _, al := range globallyEnabled {
		if _, blocked := blockedSet[al.ID]; !blocked {
			enabled = append(enabled, al)
		}
	}

	return enabled, disabled, nil
}

func (s *superAppRoleService) BulkDisableAccessLists(ctx context.Context, superappRole string, req superapproledto.BulkAccessListByRoleRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	exists, err := s.repo.RoleExists(ctx, superappRole)
	if err != nil {
		log.Errorf("[SuperAppRole][BulkDisableAccessLists] role check err: %v", err)
		return err
	}
	if !exists {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}
	log.Infof("[SuperAppRole][BulkDisableAccessLists] starting bulk disable for role %s with access list IDs: %v", superappRole, req.AccessListIDs)
	//==================
	objects, err := s.repo.FindAccessListsByIDs(ctx, req.AccessListIDs)
	if err != nil {
		log.Errorf("[SuperAppRole][BulkDisableAccessLists] fetch ids err: %v", err)
		return err
	}
	if len(objects) != len(req.AccessListIDs) {
		log.Errorf("[SuperAppRole][BulkDisableAccessLists] requested %d IDs but found %d", len(req.AccessListIDs), len(objects))
		return errors.New(localization.ErrorSomeAccesslistNotFound.Code)
	}
	log.Infof("[SuperAppRole][BulkDisableAccessLists] found access lists: %v", extractIDs(objects))

	for _, obj := range objects {
		if !obj.Enabled {
			log.Errorf("[SuperAppRole][BulkDisableAccessLists] access list %s is globally disabled", obj.ID)
			return errors.New(localization.ErrorSomeAccesslistAlreadyDisabled.Code)
		}
	}

	roleBlocked, err := s.repo.FindRoleBlockedAccessLists(ctx, superappRole)
	if err != nil {
		log.Errorf("[SuperAppRole][BulkDisableAccessLists] role blocked fetch err: %v", err)
		return err
	}
	blockedSet := make(map[string]struct{}, len(roleBlocked))
	for _, al := range roleBlocked {
		blockedSet[al.ID] = struct{}{}
	}
	for _, obj := range objects {
		if _, alreadyBlocked := blockedSet[obj.ID]; alreadyBlocked {
			log.Errorf("[SuperAppRole][BulkDisableAccessLists] access list %s is already role-disabled for role %s", obj.ID, superappRole)
			return errors.New(localization.ErrorSomeAccesslistAlreadyDisabled.Code)
		}
	}

	prev := setEnabled(objects, true)
	curr := setEnabled(objects, false)

	log.Infof("[SuperAppRole][BulkDisableAccessLists] prepared prev data: %v and curr data: %v states for role %s: prev enabled IDs: %v, curr disabled IDs: %v", prev, curr, superappRole, extractIDs(prev), extractIDs(curr))
	cpsActionData := lib.CpsModelBuilder(
		superappRole,
		makerUser,
		prev,
		curr,
		string(constants.RequestBulkDisableAccessListByRole),
		constants.UPDATE,
	)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[SuperAppRole][BulkDisableAccessLists] cps action err: %v", err)
		return err
	}
	return nil
}

func (s *superAppRoleService) BulkEnableAccessLists(ctx context.Context, superappRole string, req superapproledto.BulkAccessListByRoleRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	makerUser := local_util.ExtractUserFromContext(ctx)

	exists, err := s.repo.RoleExists(ctx, superappRole)
	if err != nil {
		log.Errorf("[SuperAppRole][BulkEnableAccessLists] role check err: %v", err)
		return err
	}
	if !exists {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	log.Infof("[SuperAppRole][BulkEnableAccessLists] starting bulk enable for role %s with access list IDs: %v", superappRole, req.AccessListIDs)
	//==================
	objects, err := s.repo.FindAccessListsByIDs(ctx, req.AccessListIDs)
	if err != nil {
		log.Errorf("[SuperAppRole][BulkEnableAccessLists] fetch ids err: %v", err)
		return err
	}

	if len(objects) != len(req.AccessListIDs) {
		log.Errorf("[SuperAppRole][BulkEnableAccessLists] requested %d IDs but found %d", len(req.AccessListIDs), len(objects))
		return errors.New(localization.ErrorSomeAccesslistNotFound.Code)
	}

	// roleBlocked, err := s.repo.FindRoleBlockedAccessLists(ctx, superappRole)
	// if err != nil {
	// 	log.Errorf("[SuperAppRole][BulkEnableAccessLists] role blocked fetch err: %v", err)
	// 	return err
	// }
	// blockedSet := make(map[string]struct{}, len(roleBlocked))
	// for _, al := range roleBlocked {
	// 	blockedSet[al.ID] = struct{}{}
	// }

	accessListSegData, err := s.repo.FindBlockedAccessListsByIDs(ctx, superappRole, req.AccessListIDs)
	if err != nil {
		log.Errorf("[SuperAppRole][BulkEnableAccessLists] fetch ids err: %v", err)
		return err
	}

	if len(accessListSegData) < len(req.AccessListIDs) {
		log.Errorf("[SuperAppRole][BulkEnableAccessLists] some access lists are not role-blocked for role %s (requested %d, role-blocked %d)", superappRole, len(req.AccessListIDs), len(accessListSegData))
		return errors.New(localization.ErrorSomeAccesslistAlreadyEnabled.Code)
	}

	globallyDisabledServices, err := s.repo.FindGloballyDisabledAccessLists(ctx)
	if err != nil {
		log.Errorf("[SuperAppRole][BulkEnableAccessLists] fetch globally disabled ids err: %v", err)
		return err
	}

	globallyDisabledSet := make(map[string]struct{}, len(globallyDisabledServices))
	for _, al := range globallyDisabledServices {
		globallyDisabledSet[al.ID] = struct{}{}
	}
	for _, obj := range accessListSegData {
		if _, isGloballyDisabled := globallyDisabledSet[obj.ID]; isGloballyDisabled {
			log.Errorf("[SuperAppRole][BulkEnableAccessLists] access list %s is globally disabled, cannot enable for role %s", obj.ID, superappRole)
			return errors.New(localization.ErrorSomeAccesslistAlreadyDisabledGlobaly.Code)
		}
	}

	// for _, obj := range objects {
	// 	if _, inBlockList := blockedSet[obj.ID]; !inBlockList {
	// 		log.Errorf("[SuperAppRole][BulkEnableAccessLists] access list %s is not role-blocked for role %s, already enabled", obj.ID, superappRole)
	// 		return errors.New(localization.ErrorSomeAccesslistAlreadyEnabled.Code)
	// 	}
	// }

	prev := setEnabled(objects, false)
	curr := setEnabled(objects, true)

	log.Infof("[SuperAppRole][BulkEnableAccessLists] prepared prev data: %v and curr data: %v states for role %s: prev enabled IDs: %v, curr disabled IDs: %v", prev, curr, superappRole, extractIDs(prev), extractIDs(curr))

	cpsActionData := lib.CpsModelBuilder(
		superappRole,
		makerUser,
		prev,
		curr,
		string(constants.RequestBulkEnableAccessListByRole),
		constants.UPDATE,
	)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[SuperAppRole][BulkEnableAccessLists] cps action err: %v", err)
		return err
	}
	return nil
}

func setEnabled(objects []imodel.APPAccessList, enabled bool) []imodel.APPAccessList {
	copies := make([]imodel.APPAccessList, len(objects))
	copy(copies, objects)
	for i := range copies {
		copies[i].Enabled = enabled
	}
	return copies
}

func extractIDs(objects []imodel.APPAccessList) []string {
	ids := make([]string, len(objects))
	for i, o := range objects {
		ids[i] = o.ID
	}
	return ids
}

func (s *superAppRoleService) authorizeBulkAccessList(ctx context.Context, role string, currentAction interface{}, repoFn func(context.Context, string, []string) error) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	lists, err := local_util.JsonUnmarshal[[]imodel.APPAccessList](currentAction)
	if err != nil {
		log.Errorf("[SuperAppRole][authorizeBulkAccessList] unmarshal err: %v", err)
		return errors.New(localization.ErrorInvalidActionData.Code)
	}
	if len(*lists) == 0 {
		log.Errorf("[SuperAppRole][authorizeBulkAccessList] currentAction contains no access list entries for role %s", role)
		return errors.New(localization.ErrorNoDataProvided.Code)
	}
	return repoFn(ctx, role, extractIDs(*lists))
}

func (s *superAppRoleService) Authorize(ctx context.Context, action *shared_model.CPSAction) (*shared_model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	log.Infof("[SuperAppRole][Authorize] action: %s, role: %s", action.RequestAction, action.UniqueId)

	role := action.UniqueId
	if role == "" {
		return nil, errors.New(localization.ErrorNoDataProvided.Code)
	}

	switch action.RequestAction {
	case string(constants.RequestEnableSuperAppRole):
		if err := s.repo.EnableByRole(ctx, role); err != nil {
			log.Errorf("[SuperAppRole][Authorize] enable err: %v", err)
			return nil, err
		}
	case string(constants.RequestDisableSuperAppRole):
		if err := s.repo.DisableByRole(ctx, role); err != nil {
			log.Errorf("[SuperAppRole][Authorize] disable err: %v", err)
			return nil, err
		}
	case string(constants.RequestDeleteSuperAppRole):
		if err := s.repo.DeleteByRole(ctx, role); err != nil {
			log.Errorf("[SuperAppRole][Authorize] delete err: %v", err)
			return nil, err
		}
	case string(constants.RequestBulkDisableAccessListByRole):
		return action, s.authorizeBulkAccessList(ctx, role, action.CurrentAction, s.repo.BulkDisableAccessLists)
	case string(constants.RequestBulkEnableAccessListByRole):
		return action, s.authorizeBulkAccessList(ctx, role, action.CurrentAction, s.repo.BulkEnableAccessLists)
	default:
		log.Errorf("[SuperAppRole][Authorize] unsupported action: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	return action, nil
}
