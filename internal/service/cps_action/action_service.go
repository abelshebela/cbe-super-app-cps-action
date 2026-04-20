package cpsaction

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	actionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type cpsActionService struct {
	repo         storage.CPSActionRepository
	roles        storage.CPSActionRoleRepository
	logger       utils.Logger
	dispatcher   Dispatcher
	minioClient  *s3.Client
	buckerName   string
	minioBaseURL string
	cfg          config.VaultConfig
}

// IsMakerOnlyForRequest returns true if the module mapped from requestAction is configured as maker-only in CPSActionRole.
func (ca *cpsActionService) IsMakerOnlyForRequest(ctx context.Context, requestAction string) (bool, error) {
	if mod, ok := ResolveModuleForRA(constants.RequestAction(requestAction)); ok && ca.roles != nil {
		role, err := ca.roles.FindByActionName(ctx, mod)
		if err != nil || role == nil {
			return false, errors.New(localization.ErrorOperationNotAllowed.Code)
		}
		return role.IsMakerOnly, nil
	}
	return false, errors.New(localization.ErrorOperationNotAllowed.Code)
}

func (ca *cpsActionService) AuditorClaim(ctx context.Context, actionCode string, activeGroup int) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "AuditorClaim", "CPSAction", "AuditorClaim")
	defer span.End()
	act, err := ca.repo.SanitizedFindOne(ctx, bson.M{"action_code": actionCode})
	if err != nil || act == nil {
		return errors.New(localization.ErrorResourceNotFound.Code)
	}

	current := int64(0)
	if act.CurrentAuditorIndex > 0 {
		current = int64(act.CurrentAuditorIndex)
	}

	expected := int64(activeGroup)
	if current != 0 && current != expected {
		return errors.New(localization.ErrorOperationNotAllowed.Code)
	}

	upd := model.CPSAction{ActionCode: actionCode}
	upd.AuditorStatus = model.AuditorStatus(constants.AUDITORINPROGRESS)
	_, err = ca.repo.Update(ctx, actionCode, upd, "", nil)
	return err
}

func (ca *cpsActionService) AuditorMark(ctx context.Context, actionCode string, auditor model.Auditor, activeGroup int) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "AuditorMark", "CPSAction", "AuditorMark")
	defer span.End()
	act, err := ca.repo.SanitizedFindOne(ctx, bson.M{"action_code": actionCode})
	if err != nil || act == nil {
		ca.logger.Errorf("[CpsActionSvc][AuditorMark] find err: %v", err)
		return errors.New(localization.ErrorActionNotFound.Code)
	}

	// prevent multiple marks within the same group (any-one quorum)
	grp := int(activeGroup)
	for _, au := range act.AuditorUsers {
		if int(au.AuditorIndex) == grp {
			return errors.New(localization.ErrorOperationNotAllowed.Code)
		}
	}

	// append this auditor
	nextUsers := append(act.AuditorUsers, auditor)

	upd := model.CPSAction{ActionCode: actionCode}
	upd.AuditorUsers = nextUsers

	// advance group or finish
	if act.AuditorCount > 0 && int32(activeGroup) >= act.AuditorCount {
		upd.AuditorStatus = model.AuditorStatus(constants.AUDITORCHECKED)
		upd.CurrentAuditorIndex = float64(activeGroup)
	} else {
		upd.AuditorStatus = model.AuditorStatus(constants.AUDITORINPROGRESS)
		upd.CurrentAuditorIndex = float64(activeGroup + 1)
	}
	_, err = ca.repo.UpdateByActionCode(ctx, actionCode, upd)
	return err
}

func NewCPSActionService(roles storage.CPSActionRoleRepository, repo storage.CPSActionRepository, logger utils.Logger, dispatcher Dispatcher, minioClient *s3.Client, bucketName string, minioBaseURL string, cfg config.VaultConfig,
) service.CPSActionService {
	return &cpsActionService{
		repo:         repo,
		logger:       logger,
		roles:        roles,
		dispatcher:   dispatcher,
		minioClient:  minioClient,
		buckerName:   bucketName,
		minioBaseURL: minioBaseURL,
		cfg:          cfg,
	}
}

func (ca *cpsActionService) CreateCPSAction(ctx context.Context, cpsAction *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateCPSAction", "CPSAction", "CreateCPSAction")
	defer span.End()
	var existing *model.CPSAction

	var err error
	ca.logger.Infof("[CpsActionSvc][Create] action: %s", cpsAction.RequestAction)

	roleCode, _ := ctx.Value(constants.ContextKey("role_code")).(string)
	actionName, _ := ctx.Value(constants.ContextKey("action_name")).(string)

	if strings.Contains(cpsAction.RequestAction, string(constants.CREATE)) {
		cpsAction.RoleCode = roleCode
		err = ca.repo.Save(ctx, cpsAction)
		if err != nil {
			span.AddEvent("failed to save cps action", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
		return nil
	}

	RAList := local_util.GetRAListForUpdateAction(constants.RequestAction(cpsAction.RequestAction), actionName, RequestActionGroups)

	reqs := ca.pendingLockRequestActions(actionName, cpsAction.RequestAction)
	if len(reqs) > 0 {
		for _, r := range RAList {
			if slices.Contains(reqs, string(r)) {
				continue
			}
			reqs = append(reqs, r)
		}
		existing, err = ca.GetPendingCPSActionByRoleAndRequestActions(ctx, cpsAction.UniqueId, reqs)
		if err != nil && err.Error() != localization.ErrorActionNotFound.Code {
			span.AddEvent("failed to get cps action by role and request actions", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
	}

	if existing != nil {
		span.AddEvent("pending cps action exists", trace.WithAttributes(attribute.String("error", "pending cps action exists")))
		ctx = context.WithValue(ctx, constants.ContextKey("existing_action_code"), existing.ActionCode)
		ctx = context.WithValue(ctx, constants.ContextKey("existing_action_status"), existing.ActionStatus)
		return errors.New(localization.ErrorPendingCpsActionExists.Code)
	}

	cpsAction.RoleCode = roleCode
	err = ca.repo.Save(ctx, cpsAction)
	if err != nil {
		span.AddEvent("failed to save cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}

func (ca *cpsActionService) pendingLockRequestActions(actionName string, requestAction string) []string {
	normalize := func(s string) string {
		return strings.ToUpper(strings.TrimSpace(s))
	}
	isCreate := func(s string) bool {
		return strings.Contains(normalize(s), constants.CREATE)
	}

	// isDelete := func(s string) bool {
	// 	return strings.Contains(normalize(s), constants.DELETE)
	// }

	defaultReq := []string{normalize(requestAction)}
	if actionName == "" {
		return defaultReq
	}

	lst, ok := RequestActionGroups[actionName]
	if !ok {
		return defaultReq
	}

	wantCreateOnly := isCreate(requestAction)
	// wantDeleteOnly := isDelete(requestAction)
	seen := map[string]struct{}{}
	reqs := make([]string, 0, len(lst))

	for _, ra := range lst {
		key := string(ra)
		if _, ok := seen[key]; ok {
			continue
		}

		isKeyCreate := isCreate(key)
		// isKeyDelete := isDelete(key)
		if wantCreateOnly != isKeyCreate {
			continue
		}

		// if wantDeleteOnly != isKeyDelete {
		// 	continue
		// }

		seen[key] = struct{}{}
		reqs = append(reqs, key)
	}

	if len(reqs) == 0 {
		return defaultReq
	}
	return reqs
}

func (ca *cpsActionService) ApproveCPSAction(ctx context.Context, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "ApproveCPSAction", "CPSAction", "ApproveCPSAction")
	defer span.End()
	ca.logger.Infof("[CpsActionSvc][Approve] action: %s", action.ActionCode)

	mod, ok := ResolveModuleForRA(constants.RequestAction(action.RequestAction))
	if !ok {
		span.AddEvent("failed to resolve module for request action", trace.WithAttributes(attribute.String("error", "failed to resolve module for request action")))
		ca.logger.Errorf("[CpsActionSvc][Approve] failed to resolve module for request action: %s", action.RequestAction)
		return errors.New(localization.ErrorOperationNotAllowed.Code)
	}

	data, err := ca.repo.Update(ctx, action.ActionCode, *action, mod, RequestActionGroups)
	if err != nil {
		span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		ca.logger.Errorf("[CpsActionSvc][Approve] update err: %v", err)
		return err
	}

	if action.ActionStatus != string(constants.Approved) {
		return nil
	}

	approve, err := ca.dispatcher.Authorize(ctx, data)
	if err != nil && approve == nil {
		span.AddEvent("failed to authorize cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		ca.logger.Errorf("[CpsActionSvc][Approve] authorize err: %v", err)
		RollErr := ca.RollBack(ctx, data)
		if RollErr != nil {
			span.AddEvent("failed to roll back cps action", trace.WithAttributes(attribute.String("error", RollErr.Error())))
			ca.logger.Errorf("[CpsActionSvc][Approve] rollback err: %v", RollErr)
			return RollErr
		}
		if err.Error() == localization.ErrorTimeoutError.Code {
			return err
		}
		return err
	}
	return nil

}
func (ca *cpsActionService) RejectCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "RejectCPSAction", "CPSAction", "RejectCPSAction")
	defer span.End()
	_, err := ca.repo.Update(ctx, action_code, *action, "", nil)
	if err != nil {
		span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}
func (ca *cpsActionService) CancelCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CancelCPSAction", "CPSAction", "CancelCPSAction")
	defer span.End()
	_, err := ca.repo.Update(ctx, action_code, *action, "", nil)
	if err != nil {
		span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}
func (ca *cpsActionService) GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionsByDepartment", "CPSAction", "GetCPSActionsByDepartment")
	defer span.End()
	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, department)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ca *cpsActionService) GetCPSActionsForApprover(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionsForApprover", "CPSAction", "GetCPSActionsForApprover")
	defer span.End()
	createdAtFrom, _ := filterParams.Filters["created_at_from"].(string)
	createdAtTo, _ := filterParams.Filters["created_at_to"].(string)
	result, err := ca.repo.SanitizedFindAllWithPaginationForApprover(ctx, userID, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, "", err
	}

	if filterParams.Filters["action"] == "export" { // checked
		// Validate date filters for export
		if createdAtFrom == "" || createdAtTo == "" {
			span.AddEvent("missing date filters for export")
			ca.logger.Errorf("[CpsActionSvc][Export] missing required date filters: from=%q, to=%q", createdAtFrom, createdAtTo)
			return nil, "", errors.New(localization.ErrorRequiredFieldMissing.Code)
		}

		filterParams.Filters["created_at_to"] = createdAtTo
		filterParams.Filters["created_at_from"] = createdAtFrom

		url, err := lib.FileExporterForCPSAction(ctx, ca.cfg, ca.minioClient, ca.buckerName, filterParams, result.Data, CpsActionCSVHeader, ca.logger)
		if err != nil {
			span.AddEvent("failed to export CPS actions", trace.WithAttributes(attribute.String("error", err.Error())))
			ca.logger.Errorf("[CpsActionSvc][Export] export CPS actions err: %v", err)
			return nil, "", err
		}
		return result, url, nil
	}
	return result, "", nil
}

func (ca *cpsActionService) exportCPSActions(ctx context.Context, filterParams *types.Filter, result *types.PaginatedResponse[[]*model.CPSAction]) (string, error) {
	startDate, endDate, err := local_util.FormatDateRangeToUTCStrings(filterParams.Filters["created_at_from"].(string), filterParams.Filters["created_at_to"].(string))
	if err != nil {
		ca.logger.Errorf("[CpsActionSvc][Export] format date range to UTC strings err: %v", err)
		return "", errors.New(localization.ErrorInvalidDateFormat.Code)
	}

	url, err := lib.ProduceFileFromData(ctx, lib.FileProducerConfig{
		FilePrefix: "cps_actions",
		Header:     CpsActionCSVHeader([]string{"Action Code", "Action Name", "Action Status", "Action Type", "Action Date"}),
		FileType:   lib.FileTypeCSV,
		ObjectName: fmt.Sprintf("cps_actions_%s_to_%s_%d.csv", startDate.Format("20060102"), endDate.Format("20060102"), time.Now().Unix()),
	}, result.Data, func(action *model.CPSAction) ([]string, error) {
		return []string{
			action.ActionCode,
			action.MakerName, // Assuming Action Name is Maker Name
			action.ActionStatus,
			action.ActionType,
			formatTime(action.CreatedAt), // Action Date
		}, nil
	}, func(ctx context.Context, file *os.File, size int64, objectName string, ft lib.FileType) (string, error) {

		switch ft {
		case lib.FileTypeCSV:
			return lib.UploadCSVToMinio(ctx, ca.minioClient, ca.buckerName, file, size, ca.cfg, objectName, ca.logger)

		case lib.FileTypePDF:
			return lib.UploadPDFToMinio(ctx, ca.minioClient, ca.buckerName, file, size, ca.cfg, objectName, ca.logger)

		default:
			return "", fmt.Errorf("unsupported file type: %s", ft)
		}
	})
	if err != nil {
		ca.logger.Errorf("[CpsActionSvc][Export] produce file from data err: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	return url, nil
}
func (ca *cpsActionService) GetCPSActionsForAuditor(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionsForAuditor", "CPSAction", "GetCPSActionsForAuditor")
	defer span.End()
	createdAtFrom, _ := filterParams.Filters["created_at_from"].(string)
	createdAtTo, _ := filterParams.Filters["created_at_to"].(string)
	result, err := ca.repo.SanitizedFindAllWithPaginationForAuditor(ctx, userID, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, "", err
	}

	if filterParams.Filters["action"] == "export" { // checked
		// Validate date filters for export
		if createdAtFrom == "" || createdAtTo == "" {
			span.AddEvent("missing date filters for export")
			ca.logger.Errorf("[CpsActionSvc][Export] missing required date filters: from=%q, to=%q", createdAtFrom, createdAtTo)
			return nil, "", errors.New(localization.ErrorRequiredFieldMissing.Code)
		}

		filterParams.Filters["created_at_to"] = createdAtTo
		filterParams.Filters["created_at_from"] = createdAtFrom

		url, err := lib.FileExporterForCPSAction(ctx, ca.cfg, ca.minioClient, ca.buckerName, filterParams, result.Data, CpsActionCSVHeader, ca.logger)
		if err != nil {
			span.AddEvent("failed to export CPS actions", trace.WithAttributes(attribute.String("error", err.Error())))
			ca.logger.Errorf("[CpsActionSvc][Export] export CPS actions err: %v", err)
			return nil, "", err
		}
		return result, url, nil
	}
	return result, "", nil
}

func (ca *cpsActionService) GetCPSActions(ctx context.Context, userID, role string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActions", "CPSAction", "GetCPSActions")
	defer span.End()

	result, err := ca.repo.SanitizedFindAllWithPaginationCPSActions(ctx, userID, role, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}

	return result, nil
}

func (ca *cpsActionService) GetCPSActionByID(ctx context.Context, id, department string) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionByID", "CPSAction", "GetCPSActionByID")
	defer span.End()
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		span.AddEvent("failed to parse the string to bson object", trace.WithAttributes(attribute.String("error", err.Error())))
		ca.logger.Errorf("[CpsActionSvc][GetByID] parse id err")
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	action, err := ca.repo.SanitizedFindOne(ctx, bson.M{"_id": objID, "department": department})
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}
func (ca *cpsActionService) GetCPSActionByUniqueID(ctx context.Context, requestAction, role_code string) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionByUniqueID", "CPSAction", "GetCPSActionByUniqueID")
	defer span.End()

	filter := bson.M{
		"role_code":      role_code,
		"action_status":  string(constants.Pending),
		"request_action": requestAction,
	}

	action, err := ca.repo.SanitizedFindOne(ctx, filter)
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}

func (ca *cpsActionService) GetPendingCPSActionByRoleAndRequestActions(ctx context.Context, uniqueId string, requestActions []string) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPendingCPSActionByRoleAndRequestActions", "CPSAction", "GetPendingCPSActionByRoleAndRequestActions")
	defer span.End()

	filter := bson.M{
		"unique_id":      uniqueId,
		"action_status":  string(constants.Pending),
		"request_action": bson.M{"$in": requestActions},
	}

	action, err := ca.repo.SanitizedFindOne(ctx, filter)
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		ca.logger.Errorf("[CpsActionSvc][GetPendingByRoleAndRA] find one err: %v", err)
		return nil, err
	}

	ca.logger.Infof("[CpsActionSvc][GetPendingByRoleAndRA] found pending action: %s for uniqueId: %s and requestActions: %v", action.ActionCode, uniqueId, requestActions)
	return action, nil
}

func (ca *cpsActionService) GetCPSActionByForUpdate(ctx context.Context) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionByForUpdate", "CPSAction", "GetCPSActionByUniqueID")
	defer span.End()

	var reqs []string
	seen := map[string]struct{}{}

	actionName, _ := ctx.Value(constants.ContextKey("action_name")).(string)

	if lst, ok := RequestActionGroups[actionName]; ok {
		for _, ra := range lst {
			key := string(ra)
			if _, ok := seen[key]; ok {
				continue
			}

			if strings.Contains(key, constants.CREATE) {
				continue
			}

			seen[key] = struct{}{}
			reqs = append(reqs, key)
		}
	}

	filter := bson.M{
		"role_code":      ctx.Value(constants.ContextKey("role_code")).(string),
		"action_status":  string(constants.Pending),
		"request_action": bson.M{"$in": reqs},
	}

	action, err := ca.repo.SanitizedFindOne(ctx, filter)
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}

func (ca *cpsActionService) GetCPSActionByActionCode(ctx context.Context, uniqueID, department string) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionByActionCode", "CPSAction", "GetCPSActionByActionCode")
	defer span.End()

	action, err := ca.repo.SanitizedFindOne(ctx, bson.M{"action_code": uniqueID})
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return action, nil
}

// RollBack reverts a CPS action after a failed Authorize call.
// It receives the full DB document (data) so it can inspect CheckerCount, AuditorCount, etc.
//   - Maker-only (CheckerCount == 0): soft-delete the action so it does not block future creations.
//   - Multi-checker (CheckerCount > 0): revert the last checker and reset status to Pending.
//   - Auditor fields are also reset when AuditorCount > 0.
func (ca *cpsActionService) RollBack(ctx context.Context, data *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "RollBack", "CPSAction", "RollBack")
	defer span.End()

	filter := bson.M{"action_code": data.ActionCode}

	if data.CheckerCount == 0 {
		err := ca.repo.UpdateCustome(ctx, filter, bson.M{
			"action_status":   string(constants.Canceled),
			"canceled_reason": constants.RoleBackReason,
		})

		if err != nil {
			span.AddEvent("failed to soft-delete maker-only cps action", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
		ca.logger.Infof("[CpsActionSvc][RollBack] soft-deleted: %s", data.ActionCode)
		return nil
	}

	previousCheckers := data.CheckerUsers
	if len(previousCheckers) > 0 {
		previousCheckers = previousCheckers[:len(previousCheckers)-1]
	}

	previousIndex := data.CurrentCheckerIndex - 1
	if previousIndex < 0 {
		previousIndex = 0
	}

	update := bson.M{
		"action_status":         string(constants.Pending),
		"checker_users":         previousCheckers,
		"current_checker_index": previousIndex,
	}

	if data.AuditorCount > 0 {
		update["auditor_users"] = []model.Auditor{}
		update["auditor_status"] = string(model.AUDITORNOTCHECKED)
		update["current_auditor_index"] = 0
	}

	err := ca.repo.UpdateCustome(ctx, filter, update)
	if err != nil {
		span.AddEvent("failed to roll back cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	ca.logger.Infof("[CpsActionSvc][RollBack] reverted to pending: %s", data.ActionCode)
	return nil
}

func (ca *cpsActionService) GetActionCountsByDepartemnt(ctx context.Context, department string) (*actionDto.CPSActionCountResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetActionCountsByDepartemnt", "CPSAction", "GetActionCountsByDepartemnt")
	defer span.End()
	count, err := ca.repo.GetCountByDepartment(ctx, department)
	if err != nil {
		span.AddEvent("failed to get count by department", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return count, nil
}

func (ca *cpsActionService) GetUserAuthorizerIndex(ctx context.Context, requestAction constants.RequestAction) (imodel.CPSActionApproveIndex, error) {
	roleCode, _ := ctx.Value(constants.ContextKey("role_code")).(string)
	var approverData imodel.CPSActionApproveIndex

	ca.logger.Infof("[CpsActionSvc][GetAuthIdx] role: %s action: %s", roleCode, requestAction)
	if mod, ok := ResolveModuleForRA(constants.RequestAction(requestAction)); ok && ca.roles != nil {
		if approver, err := ca.roles.FindApproverByActionName(ctx, strings.ToUpper(mod), roleCode); err == nil {
			approverData = approver
		}
	} else {
		return imodel.CPSActionApproveIndex{}, errors.New(localization.ErrorOperationNotAllowed.Message)
	}

	return approverData, nil
}

func (ca *cpsActionService) GetUserCreatedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetUserCreatedActions", "CPSAction", "GetUserCreatedActions")
	defer span.End()
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	filterParams.Filters["maker_id"] = userID
	createdAtFrom, _ := filterParams.Filters["created_at_from"].(string)
	createdAtTo, _ := filterParams.Filters["created_at_to"].(string)
	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")
	if err != nil {
		span.AddEvent("failed to find pending cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, "", err
	}

	if filterParams.Filters["action"] == "export" {
		// Validate date filters for export
		if createdAtFrom == "" || createdAtTo == "" {
			span.AddEvent("missing date filters for export")
			ca.logger.Errorf("[CpsActionSvc][Export] missing required date filters: from=%q, to=%q", createdAtFrom, createdAtTo)
			return nil, "", errors.New(localization.ErrorRequiredFieldMissing.Code)
		}

		filterParams.Filters["created_at_to"] = createdAtTo
		filterParams.Filters["created_at_from"] = createdAtFrom

		ca.logger.Infof("[CpsActionSvc][Export] export CPS actions with filters: %v and length: %v", filterParams.Filters, len(result.Data))
		url, err := lib.FileExporterForCPSAction(ctx, ca.cfg, ca.minioClient, ca.buckerName, filterParams, result.Data, CpsActionCSVHeader, ca.logger)
		if err != nil {
			span.AddEvent("failed to export CPS actions", trace.WithAttributes(attribute.String("error", err.Error())))
			ca.logger.Errorf("[CpsActionSvc][Export] export CPS actions err: %v", err)
			return nil, "", err
		}
		return result, url, nil
	}

	return result, "", nil
}

func (ca *cpsActionService) GetUserCheckedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetUserCheckedActions", "CPSAction", "GetUserCheckedActions")
	defer span.End()
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	filterParams.Filters["maker_id"] = userID
	createdAtFrom, _ := filterParams.Filters["created_at_from"].(string)
	createdAtTo, _ := filterParams.Filters["created_at_to"].(string)
	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")
	if err != nil {
		span.AddEvent("failed to find approver cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, "", err
	}

	if filterParams.Filters["action"] == "export" {
		// Validate date filters for export
		if createdAtFrom == "" || createdAtTo == "" {
			span.AddEvent("missing date filters for export")
			ca.logger.Errorf("[CpsActionSvc][Export] missing required date filters: from=%q, to=%q", createdAtFrom, createdAtTo)
			return nil, "", errors.New(localization.ErrorRequiredFieldMissing.Code)
		}

		filterParams.Filters["created_at_to"] = createdAtTo
		filterParams.Filters["created_at_from"] = createdAtFrom

		url, err := lib.FileExporterForCPSAction(ctx, ca.cfg, ca.minioClient, ca.buckerName, filterParams, result.Data, CpsActionCSVHeader, ca.logger)
		if err != nil {
			span.AddEvent("failed to export CPS actions", trace.WithAttributes(attribute.String("error", err.Error())))
			ca.logger.Errorf("[CpsActionSvc][Export] export CPS actions err: %v", err)
			return nil, "", err
		}
		return result, url, nil
	}
	return result, "", nil
}

func (ca *cpsActionService) ExportCpsActionData(
	ctx context.Context, req []string, filterMap *types.Filter,
	exportType string,
) (string, error) {

	var filterFields []string
	var rowCount int

	// 1 Create temp file
	tmpFile, err := os.CreateTemp("", "cps_actions_*.csv")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)

	startDate, endDate, err := local_util.FormatDateRangeToUTCStrings(filterMap.Filters["created_at_from"].(string), filterMap.Filters["created_at_to"].(string))
	if err != nil {
		ca.logger.Errorf("[CpsActionSvc][Export] format date range to UTC strings err: %v", err)
		return "", errors.New(localization.ErrorInvalidDateFormat.Code)
	}
	filterMap.Filters["created_at_from"] = startDate
	filterMap.Filters["created_at_to"] = endDate
	// 2️ Write Header
	if filterMap.Filters == nil {
		filterMap.Filters = map[string]interface{}{}
	}

	fields, ok := filterMap.Filters["fields"].([]string)
	if ok {
		filterFields = fields
	}
	// delete(filterMap.Filters, "fields")s

	if err := writer.Write(CpsActionCSVHeader(filterFields)); err != nil {
		return "", fmt.Errorf("write header: %w", err)
	}
	//==================================

	actions, err := ca.repo.ActionByDateRange(ctx, *filterMap, req)
	if err != nil {
		ca.logger.Errorf("[CpsActionSvc][Export] find all with date range err: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	for _, action := range actions {
		rowCount++
		if err := ca.processCPSAction(writer, action); err != nil {
			ca.logger.Errorf("[CpsActionSvc][Export] process CPS action err: %v", err)
			return "", errors.New(localization.ErrorUnexpectedError.Code)
		}
	}

	if rowCount == 0 {
		ca.logger.Infof("[CpsActionSvc][Export] no data found in date range %v - %v",
			filterMap.Filters["created_at_from"], filterMap.Filters["created_at_to"])
		return "", errors.New(localization.CpsActionDataNotFoundInDateRange.Code)
	}

	// 4️Upload to MinIO
	objectName := fmt.Sprintf(
		"cps_actions_%s_to_%s_%d.csv",
		startDate.Format("20060102"),
		endDate.Format("20060102"),
		time.Now().Unix(),
	)

	if _, err := tmpFile.Seek(0, 0); err != nil {
		ca.logger.Errorf("[CpsActionSvc][Export] seek temp file err: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	stat, err := tmpFile.Stat()
	if err != nil {
		ca.logger.Errorf("[CpsActionSvc][Export] stat temp file err: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	// exportType comes from the handler (e.g. query file_type=csv); default to csv for this endpoint.
	ft := strings.TrimSpace(strings.ToLower(exportType))
	if ft == "" {
		ft = "csv"
	}

	var url string
	if ft == "csv" {
		url, err = lib.UploadCSVToMinio(ctx, ca.minioClient, ca.buckerName, tmpFile, stat.Size(), ca.cfg, objectName, ca.logger)
		if err != nil {
			ca.logger.Errorf("[CpsActionSvc][Export] upload to MinIO err: %v", err)
			return "", errors.New(localization.CpsActionDataExportedError.Code)
		}
	} else {
		url, err = lib.UploadPDFToMinio(ctx, ca.minioClient, ca.buckerName, tmpFile, stat.Size(), ca.cfg, objectName, ca.logger)
		if err != nil {
			ca.logger.Errorf("[CpsActionSvc][Export] upload to MinIO err: %v", err)
			return "", errors.New(localization.CpsActionDataExportedError.Code)
		}
	}

	baseURL := strings.TrimSuffix(ca.minioBaseURL, "/")
	if baseURL != "" {
		url = fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(objectName, "/"))
	}

	return url, nil
}

func (ca *cpsActionService) processCPSAction(
	writer *csv.Writer,
	action *model.CPSAction,
) error {

	row, err := BuildCPSActionRow(action)
	if err != nil {
		return err
	}

	return writer.Write(row)
}

func BuildCPSActionRow(a *model.CPSAction) ([]string, error) {
	// Extract auditor names
	var auditorNames []string
	for _, auditor := range a.AuditorUsers {
		auditorNames = append(auditorNames, auditor.AuditorName)
	}
	auditorNamesStr := strings.Join(auditorNames, ", ")

	return []string{
		a.ID.Hex(),
		a.ActionCode,
		// a.UniqueId,
		a.MakerID,
		a.MakerName,
		a.MakerPhoneNumber,
		// string(checkerJSON),
		// string(auditorJSON),
		auditorNamesStr,
		// strconv.Itoa(int(a.AuditorCount)),
		string(a.AuditorStatus),
		// fmt.Sprintf("%f", a.CurrentAuditorIndex),
		// strconv.Itoa(int(a.CheckerCount)),
		// fmt.Sprintf("%f", a.CurrentCheckerIndex),
		// a.RoleCode,
		// a.RejectionReason,
		// a.CanceledReason,
		// string(prevJSON),
		// string(currJSON),
		a.ActionStatus,
		a.ActionType,
		// strconv.FormatBool(a.IsDeleted),
		a.RequestAction,
		// strconv.FormatInt(a.Version, 10),
		// a.ReversedByRoleID,
		// a.ReversedByID,
		// a.ReversedByName,
		// formatTime(a.ReversedAt),
		formatTime(a.CreatedAt),
		formatTime(a.LastModifiedAt),
		formatTime(a.MakerActionTime),
		formatTime(a.LastModifiedAt), // Using LastModifiedAt as CheckerActionTime
	}, nil
}

func CpsActionCSVHeader(fields []string) []string {
	if fields != nil {
		return fields
	}

	// default header if no specific fields are requested. The order of fields should match the order in BuildCPSActionRow.
	var defaultHeader = []string{
		"ID",
		"Action Code",
		"Maker ID",
		"Maker Name",
		"Maker Phone Number",
		"Auditor Names",
		"Auditor Status",
		"Action Status",
		"Action Type",
		"Request Action",
		"Created At",
		"Last Modified At",
		"Maker Action Time",
		"Checker Action Time",
	}

	return defaultHeader
}

func formatTime(v any) string {
	switch t := v.(type) {
	case time.Time:
		if t.IsZero() {
			return ""
		}
		return t.Format(time.RFC3339)
	case *time.Time:
		if t == nil || t.IsZero() {
			return ""
		}
		return t.Format(time.RFC3339)
	default:
		return fmt.Sprintf("%v", v)
	}
}
