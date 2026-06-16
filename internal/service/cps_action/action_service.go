package cpsaction

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	repo          storage.CPSActionRepository
	roles         storage.CPSActionRoleRepository
	actionLogRepo storage.UserActionLogRepository
	logger        utils.Logger
	dispatcher    Dispatcher
	minioClient   *s3.Client
	buckerName    string
	minioBaseURL  string
	cfg           config.VaultConfig
}

func NewCPSActionService(roles storage.CPSActionRoleRepository, repo storage.CPSActionRepository, actionLogRepo storage.UserActionLogRepository, logger utils.Logger, dispatcher Dispatcher, minioClient *s3.Client, bucketName string, minioBaseURL string, cfg config.VaultConfig,
) service.CPSActionService {
	return &cpsActionService{
		repo:          repo,
		logger:        logger,
		roles:         roles,
		actionLogRepo: actionLogRepo,
		dispatcher:    dispatcher,
		minioClient:   minioClient,
		buckerName:    bucketName,
		minioBaseURL:  minioBaseURL,
		cfg:           cfg,
	}
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
	if err != nil {
		return err
	}
	if ca.actionLogRepo != nil {
		if logErr := ca.actionLogRepo.UpdateAuditorActionStatusByActionCode(ctx, actionCode, string(constants.AUDITORINPROGRESS)); logErr != nil {
			span.AddEvent("failed to update auditor action status on claim", trace.WithAttributes(attribute.String("error", logErr.Error())))
		}
	}
	return nil
}

func (ca *cpsActionService) AuditorMark(ctx context.Context, actionCode string, auditor model.Auditor, activeGroup int) error {
	log := local_util.LoggerFromCtx(ctx, ca.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "AuditorMark", "CPSAction", "AuditorMark")
	defer span.End()
	act, err := ca.repo.SanitizedFindOne(ctx, bson.M{"action_code": actionCode})
	if err != nil || act == nil {
		log.Errorf("[CpsActionSvc][AuditorMark] find err: %v", err)
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
	if err != nil {
		return err
	}

	// Determine the auditor process state after this mark.
	actionAuditorStatus := string(constants.AUDITORINPROGRESS)
	if act.AuditorCount > 0 && int32(activeGroup) >= act.AuditorCount {
		actionAuditorStatus = string(constants.AUDITORCHECKED)
	}

	ca.logUserAction(ctx, act, imodel.AUDITOR, "", imodel.AuditorMark(auditor.AuditorMark), "", fmt.Sprintf("%d", activeGroup))
	if ca.actionLogRepo.AuditorMarkLogsByActionCode(ctx, actionCode, string(auditor.AuditorMark), actionAuditorStatus, false) != nil {
		span.AddEvent("failed to log auditor mark actions by action code", trace.WithAttributes(attribute.String("error", "failed to log auditor mark actions by action code")))
		log.Errorf("[CpsActionSvc][AuditorMark] failed to log auditor mark actions by action code: %s", actionCode)
	}
	return nil
}

func (ca *cpsActionService) logUserAction(ctx context.Context, action *model.CPSAction, responsibility imodel.UserActionResponsibility, givenStatus string, auditorMark imodel.AuditorMark, checkerLevel, auditorLevel string) {
	reqLog := local_util.LoggerFromCtx(ctx, ca.logger)
	roleCode, _ := ctx.Value(constants.ContextKey("role_code")).(string)

	userData := local_util.ExtractUserFromContext(ctx)
	if ca.actionLogRepo == nil {
		return
	}

	// userData, _ := ctx.Value(constants.ContextKey("user_data")).(types.UserContext)

	userOID, _ := bson.ObjectIDFromHex(userData.UserID)
	actionOID := action.ID

	serviceName := ""
	if mod, ok := ResolveModuleForRA(constants.RequestAction(action.RequestAction)); ok {
		serviceName = mod
	}

	actionLog := &imodel.UserActionLog{
		ID:                         bson.NewObjectID(),
		ActionID:                   actionOID,
		ActionCode:                 action.ActionCode,
		GivenActionStatus:          givenStatus,
		GivenAuditorStatus:         auditorMark,
		RequestAction:              constants.RequestAction(action.RequestAction),
		ActionTakenServiceName:     serviceName,
		CheckerLevel:               checkerLevel,
		AuditorLevel:               auditorLevel,
		UserRoleCode:               roleCode,
		UserID:                     userOID,
		Username:                   userData.UserName,
		UserPhone:                  userData.PhoneNumber,
		ActionType:                 imodel.CPSActions,
		UserActionResponsibilities: responsibility,
		CreatedAt:                  time.Now(),
	}

	if err := ca.actionLogRepo.Upsert(ctx, actionLog); err != nil {
		reqLog.Errorf("[CpsActionSvc][logUserAction] failed to log action: %v", err)
	}
}

func (ca *cpsActionService) CreateCPSAction(ctx context.Context, cpsAction *model.CPSAction) error {
	log := local_util.LoggerFromCtx(ctx, ca.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateCPSAction", "CPSAction", "CreateCPSAction")
	defer span.End()
	var existing *model.CPSAction

	var err error
	log.Infof("[CpsActionSvc][Create] action: %s", cpsAction.RequestAction)

	roleCode, _ := ctx.Value(constants.ContextKey("role_code")).(string)
	actionName, _ := ctx.Value(constants.ContextKey("action_name")).(string)

	checksum := computeActionChecksum(cpsAction.RequestAction, cpsAction.UniqueId, cpsAction.CurrentAction)
	cpsAction.Checksum = checksum

	if strings.Contains(cpsAction.RequestAction, string(constants.CREATE)) {
		dup, err := ca.repo.SanitizedFindOne(ctx, bson.M{
			"checksum":      checksum,
			"action_status": "PENDING",
			"is_deleted":    bson.M{"$ne": true},
		})
		if err != nil && err.Error() != localization.ErrorActionNotFound.Code {
			span.AddEvent("checksum lookup failed", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
		if dup != nil {
			if md := types.GetMetadata(ctx); md != nil {
				md.CPSActionCode = dup.ActionCode
			}

			ctx = context.WithValue(ctx, constants.ContextKey("existing_action_code"), dup.ActionCode)
			ctx = context.WithValue(ctx, constants.ContextKey("existing_action_status"), dup.ActionStatus)
			return errors.New(localization.ErrorDuplicatePendingCreateAction.Code)
		}

		tokens := extractUniqueTokens(cpsAction.RequestAction, cpsAction.CurrentAction)
		if len(tokens) > 0 {
			fieldDup, err := ca.repo.SanitizedFindOne(ctx, bson.M{
				"request_action": cpsAction.RequestAction,
				"action_status":  "PENDING",
				"is_deleted":     bson.M{"$ne": true},
				"unique_tokens":  bson.M{"$elemMatch": bson.M{"$in": tokens}},
			})
			if err != nil && err.Error() != localization.ErrorActionNotFound.Code {
				span.AddEvent("unique token lookup failed", trace.WithAttributes(attribute.String("error", err.Error())))
				return err
			}
			if fieldDup != nil {
				if md := types.GetMetadata(ctx); md != nil {
					md.CPSActionCode = fieldDup.ActionCode
				}

				ctx = context.WithValue(ctx, constants.ContextKey("existing_action_code"), fieldDup.ActionCode)
				ctx = context.WithValue(ctx, constants.ContextKey("existing_action_status"), fieldDup.ActionStatus)
				return errors.New(localization.ErrorDuplicatePendingCreateAction.Code)
			}
		}

		cpsAction.RoleCode = roleCode
		cpsAction.UniqueTokens = tokens
		cpsActionResult, err := ca.repo.Save(ctx, cpsAction)
		if err != nil {
			span.AddEvent("failed to save cps action", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
		if len(tokens) > 0 {
			if updErr := ca.repo.UpdateCustome(ctx,
				bson.M{"action_code": cpsActionResult.ActionCode},
				bson.M{"$set": bson.M{"unique_tokens": tokens}},
			); updErr != nil {
				log.Warnf("[CpsActionSvc][Create] failed to write unique_tokens: %v", updErr)
			}
		}
		ca.logUserAction(ctx, &cpsActionResult, imodel.MAKER, "PENDING", "", "", "")
		return nil
	}

	if strings.Contains(cpsAction.RequestAction, string(constants.UPDATE)) {
		reqs := ca.pendingUpdateLockRequestActions(actionName, cpsAction.RequestAction)
		if len(reqs) > 0 {
			existing, err = ca.GetPendingCPSActionByRoleAndRequestActions(ctx, cpsAction.UniqueId, reqs)
			if err != nil && err.Error() != localization.ErrorActionNotFound.Code {
				span.AddEvent("failed to get cps action by role and request actions", trace.WithAttributes(attribute.String("error", err.Error())))
				return err
			}
		}

		if existing != nil {
			if md := types.GetMetadata(ctx); md != nil {
				md.CPSActionCode = existing.ActionCode
			}

			span.AddEvent("pending update cps action exists", trace.WithAttributes(attribute.String("error", "pending update cps action exists")))
			ctx = context.WithValue(ctx, constants.ContextKey("existing_action_code"), existing.ActionCode)
			ctx = context.WithValue(ctx, constants.ContextKey("existing_action_status"), existing.ActionStatus)
			return errors.New(localization.ErrorPendingCpsActionExists.Code)
		}
	}

	if strings.Contains(cpsAction.RequestAction, string(constants.DELETE)) {
		reqs := ca.pendingDeleteLockRequestActions(actionName, cpsAction.RequestAction)
		if len(reqs) > 0 {
			existing, err = ca.GetPendingCPSActionByRoleAndRequestActions(ctx, cpsAction.UniqueId, reqs)
			if err != nil && err.Error() != localization.ErrorActionNotFound.Code {
				span.AddEvent("failed to get cps action by role and request actions", trace.WithAttributes(attribute.String("error", err.Error())))
				return err
			}
		}

		if existing != nil {
			if md := types.GetMetadata(ctx); md != nil {
				md.CPSActionCode = existing.ActionCode
			}

			span.AddEvent("pending delete cps action exists", trace.WithAttributes(attribute.String("error", "pending delete cps action exists")))
			ctx = context.WithValue(ctx, constants.ContextKey("existing_action_code"), existing.ActionCode)
			ctx = context.WithValue(ctx, constants.ContextKey("existing_action_status"), existing.ActionStatus)
			return errors.New(localization.ErrorPendingCpsActionExists.Code)
		}
	}

	if strings.Contains(cpsAction.RequestAction, constants.ENABLE) {
		reqs := ca.pendingEnableLockRequestActions(actionName, cpsAction.RequestAction)
		if len(reqs) > 0 {
			existing, err = ca.GetPendingCPSActionByRoleAndRequestActions(ctx, cpsAction.UniqueId, reqs)
			if err != nil && err.Error() != localization.ErrorActionNotFound.Code {
				span.AddEvent("failed to get cps action by role and request actions", trace.WithAttributes(attribute.String("error", err.Error())))
				return err
			}
		}

		if existing != nil {
			span.AddEvent("pending enable cps action exists", trace.WithAttributes(attribute.String("error", "pending enable cps action exists")))
			ctx = context.WithValue(ctx, constants.ContextKey("existing_action_code"), existing.ActionCode)
			ctx = context.WithValue(ctx, constants.ContextKey("existing_action_status"), existing.ActionStatus)

			if md := types.GetMetadata(ctx); md != nil {
				md.CPSActionCode = existing.ActionCode
			}
			return errors.New(localization.ErrorPendingCpsActionExists.Code)
		}
	}

	if strings.Contains(cpsAction.RequestAction, constants.DISABLE) {
		reqs := ca.pendingDisableLockRequestActions(actionName, cpsAction.RequestAction)
		if len(reqs) > 0 {
			existing, err = ca.GetPendingCPSActionByRoleAndRequestActions(ctx, cpsAction.UniqueId, reqs)
			if err != nil && err.Error() != localization.ErrorActionNotFound.Code {
				span.AddEvent("failed to get cps action by role and request actions", trace.WithAttributes(attribute.String("error", err.Error())))
				return err
			}
		}

		if existing != nil {
			if md := types.GetMetadata(ctx); md != nil {
				md.CPSActionCode = existing.ActionCode
			}

			span.AddEvent("pending disable cps action exists", trace.WithAttributes(attribute.String("error", "pending disable cps action exists")))
			ctx = context.WithValue(ctx, constants.ContextKey("existing_action_code"), existing.ActionCode)
			ctx = context.WithValue(ctx, constants.ContextKey("existing_action_status"), existing.ActionStatus)
			return errors.New(localization.ErrorPendingCpsActionExists.Code)
		}
	}

	cpsAction.RoleCode = roleCode
	cpsActionResult, err := ca.repo.Save(ctx, cpsAction)
	if err != nil {
		span.AddEvent("failed to save cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	ca.logUserAction(ctx, &cpsActionResult, imodel.MAKER, "PENDING", "", "", "")
	return nil
}

func (ca *cpsActionService) pendingLockRequestActions(actionName string, requestAction string) []string {
	normalize := func(s string) string {
		return strings.ToUpper(strings.TrimSpace(s))
	}
	isCreate := func(s string) bool {
		return strings.Contains(normalize(s), constants.CREATE)
	}

	isDisable := func(s string) bool {
		return strings.Contains(normalize(s), constants.DISABLE)
	}

	defaultReq := []string{normalize(requestAction)}
	if actionName == "" {
		return defaultReq
	}

	lst, ok := RequestActionGroups[actionName]
	if !ok {
		return defaultReq
	}

	wantCreateOnly := isCreate(requestAction)
	wantDisableOnly := isDisable(requestAction)
	seen := map[string]struct{}{}
	reqs := make([]string, 0, len(lst))

	for _, ra := range lst {
		key := string(ra)
		if _, ok := seen[key]; ok {
			continue
		}

		isKeyCreate := isCreate(key)
		isKeyDisable := isDisable(key)
		if wantCreateOnly != isKeyCreate {
			continue
		}

		if wantDisableOnly != isKeyDisable {
			continue
		}

		seen[key] = struct{}{}
		reqs = append(reqs, key)
	}

	if len(reqs) == 0 {
		return defaultReq
	}
	return reqs
}

// pendingUpdateLockRequestActions returns UPDATE, ENABLE, and DELETE actions for blocking UPDATE requests.
// UPDATE is blocked by any pending UPDATE, ENABLE, or DELETE action on the same resource.
func (ca *cpsActionService) pendingUpdateLockRequestActions(actionName string, requestAction string) []string {
	normalize := func(s string) string {
		return strings.ToUpper(strings.TrimSpace(s))
	}
	isUpdateEnableDisableOrDelete := func(s string) bool {
		n := normalize(s)
		return strings.Contains(n, constants.UPDATE) || strings.Contains(n, constants.ENABLE) ||
			strings.Contains(n, constants.DELETE)
	}

	defaultReq := []string{normalize(requestAction)}
	if actionName == "" {
		return defaultReq
	}

	lst, ok := RequestActionGroups[actionName]
	if !ok {
		return defaultReq
	}

	seen := map[string]struct{}{}
	reqs := make([]string, 0, len(lst))

	for _, ra := range lst {
		key := string(ra)
		if _, ok := seen[key]; ok {
			continue
		}
		if !isUpdateEnableDisableOrDelete(key) {
			continue
		}
		seen[key] = struct{}{}
		reqs = append(reqs, key)
	}

	if len(reqs) == 0 {
		return defaultReq
	}
	return reqs
}

// pendingDeleteLockRequestActions returns only DELETE actions for blocking DELETE requests
func (ca *cpsActionService) pendingDeleteLockRequestActions(actionName string, requestAction string) []string {
	normalize := func(s string) string {
		return strings.ToUpper(strings.TrimSpace(s))
	}
	isDelete := func(s string) bool {
		return strings.Contains(normalize(s), constants.DELETE)
	}

	defaultReq := []string{normalize(requestAction)}
	if actionName == "" {
		return defaultReq
	}

	lst, ok := RequestActionGroups[actionName]
	if !ok {
		return defaultReq
	}

	wantDeleteOnly := isDelete(requestAction)
	seen := map[string]struct{}{}
	reqs := make([]string, 0, len(lst))

	for _, ra := range lst {
		key := string(ra)
		if _, ok := seen[key]; ok {
			continue
		}

		isKeyDelete := isDelete(key)
		if !wantDeleteOnly || !isKeyDelete {
			continue // Only include DELETE actions
		}

		seen[key] = struct{}{}
		reqs = append(reqs, key)
	}

	if len(reqs) == 0 {
		return defaultReq
	}
	return reqs
}

// pendingEnableLockRequestActions returns UPDATE, ENABLE, and DISABLE actions for blocking ENABLE requests.
// ENABLE is blocked by any pending UPDATE, ENABLE, or DISABLE action on the same resource.
func (ca *cpsActionService) pendingEnableLockRequestActions(actionName string, requestAction string) []string {
	normalize := func(s string) string {
		return strings.ToUpper(strings.TrimSpace(s))
	}
	isUpdateOrEnable := func(s string) bool {
		n := normalize(s)
		return strings.Contains(n, constants.UPDATE) || strings.Contains(n, constants.ENABLE) ||
			strings.Contains(n, constants.DISABLE)
	}

	defaultReq := []string{normalize(requestAction)}
	if actionName == "" {
		return defaultReq
	}

	lst, ok := RequestActionGroups[actionName]
	if !ok {
		return defaultReq
	}

	seen := map[string]struct{}{}
	reqs := make([]string, 0, len(lst))

	for _, ra := range lst {
		key := string(ra)
		if _, ok := seen[key]; ok {
			continue
		}
		if !isUpdateOrEnable(key) {
			continue
		}
		seen[key] = struct{}{}
		reqs = append(reqs, key)
	}

	if len(reqs) == 0 {
		return defaultReq
	}
	return reqs
}

// pendingDisableLockRequestActions returns only DISABLE actions for blocking DISABLE requests.
// DISABLE is blocked only by another pending DISABLE action on the same resource.
func (ca *cpsActionService) pendingDisableLockRequestActions(actionName string, requestAction string) []string {
	normalize := func(s string) string {
		return strings.ToUpper(strings.TrimSpace(s))
	}
	isDisable := func(s string) bool {
		return strings.Contains(normalize(s), constants.DISABLE)
	}

	defaultReq := []string{normalize(requestAction)}
	if actionName == "" {
		return defaultReq
	}

	lst, ok := RequestActionGroups[actionName]
	if !ok {
		return defaultReq
	}

	seen := map[string]struct{}{}
	reqs := make([]string, 0, len(lst))

	for _, ra := range lst {
		key := string(ra)
		if _, ok := seen[key]; ok {
			continue
		}
		if !isDisable(key) {
			continue
		}
		seen[key] = struct{}{}
		reqs = append(reqs, key)
	}

	if len(reqs) == 0 {
		return defaultReq
	}
	return reqs
}

func (ca *cpsActionService) ApproveCPSAction(ctx context.Context, action *model.CPSAction) error {
	log := local_util.LoggerFromCtx(ctx, ca.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "ApproveCPSAction", "CPSAction", "ApproveCPSAction")
	defer span.End()
	log.Infof("[CpsActionSvc][Approve] action: %s", action.ActionCode)

	mod, ok := ResolveModuleForRA(constants.RequestAction(action.RequestAction))
	if !ok {
		span.AddEvent("failed to resolve module for request action", trace.WithAttributes(attribute.String("error", "failed to resolve module for request action")))
		log.Errorf("[CpsActionSvc][Approve] failed to resolve module for request action: %s", action.RequestAction)
		return errors.New(localization.ErrorOperationNotAllowed.Code)
	}

	data, err := ca.repo.Update(ctx, action.ActionCode, *action, mod, RequestActionGroups)
	if err != nil {
		span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[CpsActionSvc][Approve] update err: %v", err)
		return err
	}

	checkerLevel := ""
	if idx, ok := ctx.Value(constants.ContextKey("role_checker_index")).(float64); ok {
		checkerLevel = fmt.Sprintf("%.0f", idx)
	}
	ca.logUserAction(ctx, action, imodel.CHECKER, action.ActionStatus, "", checkerLevel, "")

	if action.ActionStatus != string(constants.Approved) {
		return nil
	}

	approve, err := ca.dispatcher.Authorize(ctx, data)
	if err != nil && approve == nil {
		span.AddEvent("failed to authorize cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[CpsActionSvc][Approve] authorize err: %v", err)
		RollErr := ca.RollBack(ctx, data)
		if RollErr != nil {
			span.AddEvent("failed to roll back cps action", trace.WithAttributes(attribute.String("error", RollErr.Error())))
			log.Errorf("[CpsActionSvc][Approve] rollback err: %v", RollErr)
			return RollErr
		}
		if err.Error() == localization.ErrorTimeoutError.Code {
			return err
		}
		return err
	}

	if ca.actionLogRepo.ApproveUserActionsByActionCode(ctx, action.ActionCode) != nil {
		span.AddEvent("failed to approve user actions by action code", trace.WithAttributes(attribute.String("error", "failed to approve user actions by action code")))
		log.Errorf("[CpsActionSvc][Approve] failed to approve user actions by action code: %s", action.ActionCode)
		// Note: The main operation has succeeded at this point, so we don't return an error to avoid rolling back the main operation. Instead, we log the error for further investigation.
	}

	return nil

}
func (ca *cpsActionService) RejectCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "RejectCPSAction", "CPSAction", "RejectCPSAction")
	log := local_util.LoggerFromCtx(ctx, ca.logger)
	defer span.End()
	_, err := ca.repo.Update(ctx, action_code, *action, "", nil)
	if err != nil {
		span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	checkerLevel := ""
	if idx, ok := ctx.Value(constants.ContextKey("role_checker_index")).(float64); ok {
		checkerLevel = fmt.Sprintf("%.0f", idx)
	}
	ca.logUserAction(ctx, action, imodel.CHECKER, constants.Rejected, "", checkerLevel, "")

	if ca.actionLogRepo.RejectUserActionsByActionCode(ctx, action.ActionCode) != nil {
		span.AddEvent("failed to reject user actions by action code", trace.WithAttributes(attribute.String("error", "failed to reject user actions by action code")))
		log.Errorf("[CpsActionSvc][Reject] failed to reject user actions by action code: %s", action.ActionCode)
		// Note: The main operation has succeeded at this point, so we don't return an error to avoid rolling back the main operation. Instead, we log the error for further investigation.
	}
	return nil
}
func (ca *cpsActionService) CancelCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CancelCPSAction", "CPSAction", "CancelCPSAction")
	log := local_util.LoggerFromCtx(ctx, ca.logger)
	defer span.End()
	_, err := ca.repo.Update(ctx, action_code, *action, "", nil)
	if err != nil {
		span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[CpsActionSvc][Cancel] failed to update cps action: %v", err)
		return err
	}
	ca.logUserAction(ctx, action, imodel.MAKER, "CANCELED", "", "", "")
	if ca.actionLogRepo.CancelUserActionsByActionCode(ctx, action.ActionCode) != nil {
		span.AddEvent("failed to cancel user actions by action code", trace.WithAttributes(attribute.String("error", "failed to cancel user actions by action code")))
		log.Errorf("[CpsActionSvc][Cancel] failed to cancel user actions by action code: %s", action.ActionCode)
		// Note: The main operation has succeeded at this point, so we don't return an error to avoid rolling back the main operation. Instead, we log the error for further investigation.
	}
	return nil
}
func (ca *cpsActionService) GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionsByDepartment", "CPSAction", "GetCPSActionsByDepartment")
	defer span.End()
	if filterParams.Search != "" {
		if codes, err := ca.actionLogRepo.GetActionCodesBySearch(ctx, filterParams.Search); err == nil && len(codes) > 0 {
			filterParams.Filters["search_action_codes"] = codes
		}
	}
	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, department)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ca *cpsActionService) GetCPSActionsForApprover(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], string, error) {
	log := local_util.LoggerFromCtx(ctx, ca.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionsForApprover", "CPSAction", "GetCPSActionsForApprover")
	defer span.End()

	isExport := filterParams.Filters["action"] == "export"
	createdAtFrom, _ := filterParams.Filters["created_at_from"].(string)
	createdAtTo, _ := filterParams.Filters["created_at_to"].(string)

	levels := extractStringSlice(filterParams.Filters, "levels")
	checkerLevels := extractStringSlice(filterParams.Filters, "checker_levels")
	auditorLevels := extractStringSlice(filterParams.Filters, "auditor_levels")
	services := extractStringSlice(filterParams.Filters, "services")
	statuses := extractStringSlice(filterParams.Filters, "action_status")

	log.Infof("[CPSAction][GetCPSActionsForApprover] levels: %v, checkerLevels: %v, auditorLevels: %v, services: %v, statuses: %v", levels, checkerLevels, auditorLevels, services, statuses)
	checkerLevelStatuses := extractCheckerLevelStatusPairs(filterParams.Filters)

	makerUsernames := extractStringSlice(filterParams.Filters, "maker_usernames")
	checkerUsernames := extractStringSlice(filterParams.Filters, "checker_usernames")
	auditorUsernames := extractStringSlice(filterParams.Filters, "auditor_usernames")

	searchTerm := filterParams.Search

	logFilter := imodel.UserActionLogActionCodeFilter{
		RequestActions:       RAList,
		Levels:               levels,
		CheckerLevels:        checkerLevels,
		AuditorLevels:        auditorLevels,
		Services:             services,
		ActionStatuses:       statuses,
		CheckerLevelStatuses: checkerLevelStatuses,
		MakerUsernames:       makerUsernames,
		CheckerUsernames:     checkerUsernames,
		AuditorUsernames:     auditorUsernames,
		Search:               searchTerm,
	}

	// CHECKER log rows only exist after a checker has acted; scoping by CHECKER
	// responsibility hides freshly-created PENDING actions (MAKER row only).
	isPendingOnly := len(statuses) == 1 && statuses[0] == string(constants.Pending)
	if !isPendingOnly {
		logFilter.Responsibilities = []string{string(imodel.CHECKER)}
	}

	log.Infof("[CPSAction][GetCPSActionsForApprover] filter: %v", logFilter)

	actionCodes, err := ca.actionLogRepo.GetActionCodesByActionLogFilter(ctx, logFilter)
	if err != nil {
		log.Errorf("[CpsActionSvc][GetCPSActionsForApprover] log filter err: %v", err)
		return nil, "", err
	}

	// No matching codes in the log → empty page (unless this is an export where
	// the date range alone scopes the results).
	if len(actionCodes) == 0 && !isExport {
		return &types.PaginatedResponse[[]*model.CPSAction]{
			Data: []*model.CPSAction{},
			Meta: local_util.BuildPaginationMeta(0, filterParams.Page, filterParams.PerPage),
		}, "", nil
	}

	if len(actionCodes) > 0 {
		filterParams.Filters["action_code"] = actionCodes
	}
	if filterParams.Search != "" {
		if searchCodes, err := ca.actionLogRepo.GetActionCodesBySearch(ctx, filterParams.Search); err == nil && len(searchCodes) > 0 {
			filterParams.Filters["search_action_codes"] = searchCodes
		}
	}

	if len(statuses) > 0 {
		filterParams.Filters["action_status"] = statuses
	}

	if isExport {
		if createdAtFrom == "" || createdAtTo == "" {
			log.Errorf("[CpsActionSvc][Export][Approver] missing date filters: from=%q to=%q", createdAtFrom, createdAtTo)
			return nil, "", errors.New(localization.ErrorRequiredFieldMissing.Code)
		}
		filterParams.Filters["created_at_from"] = createdAtFrom
		filterParams.Filters["created_at_to"] = createdAtTo
		filterParams.Page = 1
		filterParams.PerPage = 100000
	}

	result, err := ca.repo.SanitizedFindAllWithPaginationForApprover(ctx, userID, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, "", err
	}

	if isExport {
		filterParams.Filters["created_at_from"] = createdAtFrom
		filterParams.Filters["created_at_to"] = createdAtTo
		url, err := lib.FileExporterForCPSAction(ctx, ca.cfg, ca.minioClient, ca.buckerName, filterParams, result.Data, CpsActionCSVHeader, ca.logger)
		if err != nil {
			span.AddEvent("failed to export CPS actions", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[CpsActionSvc][Export][Approver] err: %v", err)
			return nil, "", err
		}
		return result, url, nil
	}
	return result, "", nil
}

func (ca *cpsActionService) exportCPSActions(ctx context.Context, filterParams *types.Filter, result *types.PaginatedResponse[[]*model.CPSAction]) (string, error) {
	log := local_util.LoggerFromCtx(ctx, ca.logger)

	startDate, endDate, err := local_util.FormatDateRangeToUTCStrings(filterParams.Filters["created_at_from"].(string), filterParams.Filters["created_at_to"].(string))
	if err != nil {
		log.Errorf("[CpsActionSvc][Export] format date range to UTC strings err: %v", err)
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
		log.Errorf("[CpsActionSvc][Export] produce file from data err: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}
	return url, nil
}
func (ca *cpsActionService) GetCPSActionsForAuditor(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], string, error) {
	log := local_util.LoggerFromCtx(ctx, ca.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionsForAuditor", "CPSAction", "GetCPSActionsForAuditor")
	defer span.End()

	isExport := filterParams.Filters["action"] == "export"
	createdAtFrom, _ := filterParams.Filters["created_at_from"].(string)
	createdAtTo, _ := filterParams.Filters["created_at_to"].(string)

	levels := extractStringSlice(filterParams.Filters, "levels")
	checkerLevels := extractStringSlice(filterParams.Filters, "checker_levels")
	auditorLevels := extractStringSlice(filterParams.Filters, "auditor_levels")
	services := extractStringSlice(filterParams.Filters, "services")

	var auditorMarkStatuses []string  // MARKEDASRIGHT / MARKEDASWRONG (given_auditor_status)
	var auditorStateStatuses []string // NOTCHECKED / INPROGRESS / CHECKED (action_auditor_status)

	log.Infof("[CPSAction][GetCPSActionsForAuditor] levels=%v checker_levels=%v auditor_levels=%v services=%v",
		levels, checkerLevels, auditorLevels, services)

	auditStateSet := map[string]bool{
		string(constants.AUDITORNOTCHECKED): true,
		string(constants.AUDITORINPROGRESS): true,
		string(constants.AUDITORCHECKED):    true,
	}

	for _, src := range []string{"auditor_status", "auditor_statuses"} {
		for _, v := range extractStringSlice(filterParams.Filters, src) {
			if auditStateSet[strings.ToUpper(v)] {
				auditorStateStatuses = append(auditorStateStatuses, strings.ToUpper(v))
			} else {
				auditorMarkStatuses = append(auditorMarkStatuses, strings.ToUpper(v))
			}
		}
		delete(filterParams.Filters, src)
	}

	levelClaimPairs := extractLevelClaimPairs(filterParams.Filters)
	checkerLevelStatuses := extractCheckerLevelStatusPairs(filterParams.Filters)

	makerUsernames := extractStringSlice(filterParams.Filters, "maker_usernames")
	checkerUsernames := extractStringSlice(filterParams.Filters, "checker_usernames")
	auditorUsernames := extractStringSlice(filterParams.Filters, "auditor_usernames")

	searchTerm := filterParams.Search

	log.Infof("[CPSAction][GetCPSActionsForAuditor] levels=%v checker=%v auditor=%v services=%v markStatuses=%v stateStatuses=%v search=%q",
		levels, checkerLevels, auditorLevels, services, auditorMarkStatuses, auditorStateStatuses, searchTerm)

	logFilter := imodel.UserActionLogActionCodeFilter{
		RequestActions:        RAList,
		Levels:                levels,
		CheckerLevels:         checkerLevels,
		AuditorLevels:         auditorLevels,
		Services:              services,
		MakerUsernames:        makerUsernames,
		CheckerUsernames:      checkerUsernames,
		AuditorUsernames:      auditorUsernames,
		AuditorStatuses:       auditorMarkStatuses,
		ActionAuditorStatuses: auditorStateStatuses,
		LevelClaimPairs:       levelClaimPairs,
		CheckerLevelStatuses:  checkerLevelStatuses,
		Search:                searchTerm,
	}

	// Apply AUDITOR scope only when INPROGRESS or CHECKED is explicitly requested;
	// NOTCHECKED actions have no AUDITOR log rows yet.
	requiresAuditorScope := slices.Contains(auditorStateStatuses, string(constants.AUDITORINPROGRESS)) ||
		slices.Contains(auditorStateStatuses, string(constants.AUDITORCHECKED))
	if requiresAuditorScope {
		logFilter.Responsibilities = []string{string(imodel.AUDITOR)}
	}

	actionCodes, err := ca.actionLogRepo.GetActionCodesByActionLogFilter(ctx, logFilter)
	if err != nil {
		log.Errorf("[CpsActionSvc][GetCPSActionsForAuditor] log filter err: %v", err)
		return nil, "", err
	}
	log.Infof("[CpsActionSvc][GetCPSActionsForAuditor] log filter found %d codes (markStatuses=%v stateStatuses=%v levels=%v)",
		len(actionCodes), auditorMarkStatuses, auditorStateStatuses, levels)

	// No matching codes → empty page (unless export; date range scopes that path).
	if len(actionCodes) == 0 && !isExport {
		return &types.PaginatedResponse[[]*model.CPSAction]{
			Data: []*model.CPSAction{},
			Meta: local_util.BuildPaginationMeta(0, filterParams.Page, filterParams.PerPage),
		}, "", nil
	}

	if len(actionCodes) > 0 {
		filterParams.Filters["action_code"] = actionCodes
	}
	if filterParams.Search != "" {
		if searchCodes, err := ca.actionLogRepo.GetActionCodesBySearch(ctx, filterParams.Search); err == nil && len(searchCodes) > 0 {
			filterParams.Filters["search_action_codes"] = searchCodes
		}
	}

	if statuses := extractStringSlice(filterParams.Filters, "action_status"); len(statuses) > 0 {
		filterParams.Filters["action_status"] = statuses
	}

	if isExport {
		if createdAtFrom == "" || createdAtTo == "" {
			log.Errorf("[CpsActionSvc][Export][Auditor] missing date filters: from=%q to=%q", createdAtFrom, createdAtTo)
			return nil, "", errors.New(localization.ErrorRequiredFieldMissing.Code)
		}
		filterParams.Filters["created_at_from"] = createdAtFrom
		filterParams.Filters["created_at_to"] = createdAtTo
		filterParams.Page = 1
		filterParams.PerPage = 100000
	}

	result, err := ca.repo.SanitizedFindAllWithPaginationForAuditor(ctx, userID, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, "", err
	}

	if isExport {

		filterParams.Filters["created_at_from"] = createdAtFrom
		filterParams.Filters["created_at_to"] = createdAtTo
		url, err := lib.FileExporterForCPSAction(ctx, ca.cfg, ca.minioClient, ca.buckerName, filterParams, result.Data, CpsActionCSVHeader, ca.logger)
		if err != nil {
			span.AddEvent("failed to export CPS actions", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[CpsActionSvc][Export][Auditor] err: %v", err)
			return nil, "", err
		}
		return result, url, nil
	}
	return result, "", nil
}

func (ca *cpsActionService) GetCPSActions(ctx context.Context, userID, role string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActions", "CPSAction", "GetCPSActions")
	defer span.End()
	if filterParams.Search != "" {
		if codes, err := ca.actionLogRepo.GetActionCodesBySearch(ctx, filterParams.Search); err == nil && len(codes) > 0 {
			filterParams.Filters["search_action_codes"] = codes
		}
	}
	result, err := ca.repo.SanitizedFindAllWithPaginationCPSActions(ctx, userID, role, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}

	return result, nil
}

func (ca *cpsActionService) GetCPSActionByID(ctx context.Context, id, department string) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, ca.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionByID", "CPSAction", "GetCPSActionByID")
	defer span.End()
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		span.AddEvent("failed to parse the string to bson object", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[CpsActionSvc][GetByID] parse id err")
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
	log := local_util.LoggerFromCtx(ctx, ca.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetPendingCPSActionByRoleAndRequestActions", "CPSAction", "GetPendingCPSActionByRoleAndRequestActions")
	defer span.End()

	// If requestActions is empty, return not found immediately
	if len(requestActions) == 0 {
		log.Warnf("[CpsActionSvc][GetPendingByRoleAndRA] empty requestActions provided")
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}

	filter := bson.M{
		"unique_id":      uniqueId,
		"action_status":  string(constants.Pending),
		"request_action": bson.M{"$in": requestActions},
	}

	action, err := ca.repo.SanitizedFindOne(ctx, filter)
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[CpsActionSvc][GetPendingByRoleAndRA] find one err: %v", err)
		return nil, err
	}

	log.Infof("[CpsActionSvc][GetPendingByRoleAndRA] found pending action: %s for uniqueId: %s and requestActions: %v", action.ActionCode, uniqueId, requestActions)
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
	log := local_util.LoggerFromCtx(ctx, ca.logger)

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
		log.Infof("[CpsActionSvc][RollBack] soft-deleted: %s", data.ActionCode)
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
	log.Infof("[CpsActionSvc][RollBack] reverted to pending: %s", data.ActionCode)
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
	log := local_util.LoggerFromCtx(ctx, ca.logger)

	roleCode, _ := ctx.Value(constants.ContextKey("role_code")).(string)
	var approverData imodel.CPSActionApproveIndex

	log.Infof("[CpsActionSvc][GetAuthIdx] role: %s action: %s", roleCode, requestAction)
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
	log := local_util.LoggerFromCtx(ctx, ca.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetUserCreatedActions", "CPSAction", "GetUserCreatedActions")
	defer span.End()
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}

	isExport := filterParams.Filters["action"] == "export"
	createdAtFrom, _ := filterParams.Filters["created_at_from"].(string)
	createdAtTo, _ := filterParams.Filters["created_at_to"].(string)

	levels := extractStringSlice(filterParams.Filters, "levels")
	services := extractStringSlice(filterParams.Filters, "services")

	if len(levels) > 0 || len(services) > 0 {
		// Log-only filters: query user_action_log for matching codes.
		// Pre-log actions have no log metadata and won't appear.
		actionCodes, err := ca.actionLogRepo.GetActionCodesByActionLogFilter(ctx, imodel.UserActionLogActionCodeFilter{
			MakerUserIDs:   []string{userID},
			ActionStatuses: extractStringSlice(filterParams.Filters, "action_status"),
			Levels:         levels,
			Services:       services,
		})
		if err != nil {
			log.Errorf("[CpsActionSvc][GetUserCreatedActions] log filter err: %v", err)
			return nil, "", err
		}
		filterParams.Filters["action_codes"] = actionCodes
	} else {
		// Scope directly via maker_id — covers all data including pre-log actions.
		filterParams.Filters["maker_id"] = userID
		if statuses := extractStringSlice(filterParams.Filters, "action_status"); len(statuses) > 0 {
			filterParams.Filters["action_status"] = statuses
		}
	}

	if isExport {
		if createdAtFrom == "" || createdAtTo == "" {
			log.Errorf("[CpsActionSvc][Export][Maker] missing date filters: from=%q to=%q", createdAtFrom, createdAtTo)
			return nil, "", errors.New(localization.ErrorRequiredFieldMissing.Code)
		}
		filterParams.Filters["created_at_from"] = createdAtFrom
		filterParams.Filters["created_at_to"] = createdAtTo
		filterParams.Page = 1
		filterParams.PerPage = 100000
	}

	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")
	if err != nil {
		span.AddEvent("failed to find pending cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, "", err
	}

	if isExport {
		filterParams.Filters["created_at_from"] = createdAtFrom
		filterParams.Filters["created_at_to"] = createdAtTo
		log.Infof("[CpsActionSvc][Export][Maker] exporting %d records", len(result.Data))
		url, err := lib.FileExporterForCPSAction(ctx, ca.cfg, ca.minioClient, ca.buckerName, filterParams, result.Data, CpsActionCSVHeader, ca.logger)
		if err != nil {
			span.AddEvent("failed to export CPS actions", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[CpsActionSvc][Export][Maker] err: %v", err)
			return nil, "", err
		}
		return result, url, nil
	}
	return result, "", nil
}

func (ca *cpsActionService) GetUserCheckedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], string, error) {
	log := local_util.LoggerFromCtx(ctx, ca.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetUserCheckedActions", "CPSAction", "GetUserCheckedActions")
	defer span.End()
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}

	isExport := filterParams.Filters["action"] == "export"
	createdAtFrom, _ := filterParams.Filters["created_at_from"].(string)
	createdAtTo, _ := filterParams.Filters["created_at_to"].(string)

	levels := extractStringSlice(filterParams.Filters, "levels")
	services := extractStringSlice(filterParams.Filters, "services")

	if len(levels) > 0 || len(services) > 0 {
		// Log-only filters: query user_action_log for matching codes.
		// Pre-log actions have no log metadata and won't appear.
		actionCodes, err := ca.actionLogRepo.GetActionCodesByActionLogFilter(ctx, imodel.UserActionLogActionCodeFilter{
			CheckerUserIDs: []string{userID},
			ActionStatuses: extractStringSlice(filterParams.Filters, "action_status"),
			Levels:         levels,
			Services:       services,
		})
		if err != nil {
			log.Errorf("[CpsActionSvc][GetUserCheckedActions] log filter err: %v", err)
			return nil, "", err
		}
		filterParams.Filters["action_codes"] = actionCodes
	} else {
		// Scope directly via checker_users.checker_id — covers pre-log actions too.
		filterParams.Filters["checker_users.checker_id"] = userID
		if statuses := extractStringSlice(filterParams.Filters, "action_status"); len(statuses) > 0 {
			filterParams.Filters["action_status"] = statuses
		}
	}

	if isExport {
		if createdAtFrom == "" || createdAtTo == "" {
			log.Errorf("[CpsActionSvc][Export][Checker] missing date filters: from=%q to=%q", createdAtFrom, createdAtTo)
			return nil, "", errors.New(localization.ErrorRequiredFieldMissing.Code)
		}
		filterParams.Filters["created_at_from"] = createdAtFrom
		filterParams.Filters["created_at_to"] = createdAtTo
		filterParams.Page = 1
		filterParams.PerPage = 100000
	}

	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")
	if err != nil {
		span.AddEvent("failed to find approver cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, "", err
	}

	if isExport {
		filterParams.Filters["created_at_from"] = createdAtFrom
		filterParams.Filters["created_at_to"] = createdAtTo
		log.Infof("[CpsActionSvc][Export][Checker] exporting %d records", len(result.Data))
		url, err := lib.FileExporterForCPSAction(ctx, ca.cfg, ca.minioClient, ca.buckerName, filterParams, result.Data, CpsActionCSVHeader, ca.logger)
		if err != nil {
			span.AddEvent("failed to export CPS actions", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[CpsActionSvc][Export][Checker] err: %v", err)
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
	log := local_util.LoggerFromCtx(ctx, ca.logger)

	if filterMap == nil {
		filterMap = &types.Filter{Filters: map[string]interface{}{}}
	}
	if filterMap.Filters == nil {
		filterMap.Filters = map[string]interface{}{}
	}

	exportType = strings.TrimSpace(strings.ToLower(exportType))
	if exportType == "" {
		exportType = string(lib.FileTypeCSV)
	}
	filterMap.Filters["file_type"] = exportType

	if statuses := extractStringSlice(filterMap.Filters, "action_status"); len(statuses) > 0 {
		filterMap.Filters["action_status"] = statuses
	}

	actions, err := ca.repo.ActionByDateRange(ctx, *filterMap, req)
	if err != nil {
		log.Errorf("[CpsActionSvc][Export] find all with date range err: %v", err)
		return "", errors.New(localization.ErrorUnexpectedError.Code)
	}

	return lib.FileExporterForCPSAction(ctx, ca.cfg, ca.minioClient, ca.buckerName, filterMap, actions, CpsActionCSVHeader, ca.logger)
}

// CpsActionCSVHeader resolves the user-visible column labels for the given ?fields= keys.
// Unknown/missing keys fall back to the registry default. Delegates to lib so headers and
// row extractors stay in sync.
func extractStringSlice(filters map[string]interface{}, key string) []string {
	if filters == nil {
		return nil
	}
	return local_util.StringSliceFromFilterValue(filters[key])
}

var volatileChecksumKeys = map[string]bool{
	"id":          true,
	"action_code": true,
	"version":     true,
	"user_code":   true,

	"is_deleted": true,
	"is_enabled": true,

	"created_at":                    true,
	"create_at":                     true,
	"updated_at":                    true,
	"update_at":                     true,
	"last_modified_at":              true,
	"last_modified":                 true,
	"maker_action_time":             true,
	"approved_at":                   true,
	"reversed_at":                   true,
	"deleted_at":                    true,
	"date_joined":                   true,
	"issued_date":                   true,
	"sent_at":                       true,
	"published_at":                  true,
	"verified_at":                   true,
	"completed_at":                  true,
	"claimed_at":                    true,
	"linked_at":                     true,
	"initial_linked_at":             true,
	"initiated_linked_at":           true,
	"pin_changed_at":                true,
	"password_changed_at":           true,
	"application_installation_date": true,
	"last_login":                    true,
	"last_login_attempt":            true,
	"last_online_date":              true,
	"otp_last_tried_at":             true,
	"otp_last_verified_at":          true,
	"created_at_password_expiry":    true,
	"updated_at_password_expiry":    true,
	"created_at_block":              true,
	"updated_at_block":              true,
	"created_at_archive":            true,
	"updated_at_archive":            true,
	"created_at_total_cap":          true,
	"updated_at_total_cap":          true,
	"next_attempt_count":            true,

	// --- file / image / media URLs (Minio key varies per upload) ---
	"logo":            true,
	"icon":            true,
	"app_icon":        true,
	"company_logo":    true,
	"donation_icon":   true,
	"image":           true,
	"image_url":       true,
	"cover_image":     true,
	"cover_image_url": true,
	"banner_image":    true,
	"photo":           true,
	"selfie_photo":    true,
	"picture":         true,
	"thumbnail":       true,
	"avatar":          true,
	"document_front":  true,
	"document_back":   true,
	"signature":       true,
	"video_url":       true,
	"receipt_link":    true,
	"url":             true,
}

var uniqueFieldsRegistry = map[string][]string{
	string(constants.RequestCpsUserCreate): {"phone_number", "email", "username"},
	string(constants.RequestCreateBPSUser): {"phone_number", "email", "username"},

	string(constants.RequestCreateEcommerceMerchant): {"merchant_code", "bank_account_number"},
	// mini_app.MiniAppMerchant — code + account are the primary uniqueness keys;
	// phone/email added as extra guards per product requirement.
	string(constants.RequestCreateMiniAppMerchant):   {"merchant_code", "bank_account_number", "phone_number", "email", "account_number"},
	string(constants.RequestCreateEventMerchant):     {"merchant_id", "bank_account_number"},
	string(constants.RequestCreateLogisticsMerchant): {"merchant_id", "bank_account_number"},
	string(constants.RequestCreateUssdMerchant):      {"phone_number", "email", "account_number"},

	string(constants.RequestCreateDonation):         {"title"},
	string(constants.RequestCreateDonationCategory): {"category_name"},
	string(constants.RequestCreateDonationCompany):  {"company_name", "company_code"},

	string(constants.RequestCreateDepartment):      {"department", "department_code"},
	string(constants.RequestCreateCpsRole):         {"role_code", "name"},
	string(constants.RequestCreateJobRole):         {"code", "name"},
	string(constants.RequestCreateCpsActionRole):   {"action_name", "portal_card_name"},
	string(constants.RequestCreateActionRole):      {"action_name"},
	string(constants.RequestCreatePermissionGroup): {"group_name"},

	string(constants.RequestCreateBank):            {"bank_name", "bic_code"},
	string(constants.RequestCreateBankVault):       {"name"},
	string(constants.RequestCreateAmountBasedAuth): {"currency"},

	string(constants.RequestCreateWallet): {"name", "unique_code"},
	string(constants.RequestCreateTopup):  {"name", "code"},

	string(constants.RequestCreateAdvert): {"title"},
	string(constants.RequestCreateAvatar): {"label"},

	string(constants.RequestCreateServiceList): {"service_name", "service_key"},

	string(constants.RequestCreateBudgetCategory): {"name"},
	string(constants.RequestCreateVaultCategory):  {"name"},

	// shared model.Event has no json tags — Go field name used by json.Marshal
	string(constants.RequestCreateEvent): {"EventName"},

	// ── MiniApp ──────────────────────────────────────────────────────────────
	// mini_app.MiniApp : AppName→"app_name"
	string(constants.RequestCreateMiniApp): {"app_name"},
	// model.MiniAppCategory : Name→"name"
	string(constants.RequestCreateMiniAppCategory): {"name"},

	// ── Segmentation ─────────────────────────────────────────────────────────
	// CreateAccessListSegmentationRequest : SegmentationID→"segmentation_id" (the block/account id)
	string(constants.RequestCreateAccessListSegmentation): {"segmentation_id"},
	// MapCustomerSegmentationToMap : "customer_role"→nested; access_list_id used as product key
	string(constants.RequestCreateCustomerSegmentation): {"access_list_id"},

	// ── Customer segments ────────────────────────────────────────────────────
	// payload is MapSegmentToMap — keys are explicit strings
	string(constants.RequestCreateCustomerGroup): {"customer_group", "customer_segment", "customer_subsegment"},

	// ── News / Media ─────────────────────────────────────────────────────────
	// NewsTagCPSAction : TagName→"tag_name" (single-name field; tag_name_list slice is not supported)
	string(constants.RequestCreateNewsTag): {"tag_name"},
	// NewsCategoryCPSAction : CategoryName→"category_name"
	string(constants.RequestCreateNewsCategory): {"category_name"},
	// shared NewsArticle : Title→"title"
	string(constants.RequestCreateArticle): {"title"},
}

func extractUniqueTokens(requestAction string, currentAction interface{}) []string {
	fields, ok := uniqueFieldsRegistry[requestAction]
	if !ok || len(fields) == 0 {
		return nil
	}

	raw, _ := json.Marshal(currentAction)
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil
	}

	tokens := make([]string, 0, len(fields))
	for _, field := range fields {
		val, ok := m[field]
		if !ok || val == nil {
			continue
		}
		str, ok := val.(string)
		if !ok || strings.TrimSpace(str) == "" {
			continue
		}
		tokens = append(tokens, field+":"+strings.ToLower(strings.TrimSpace(str)))
	}
	return tokens
}

func computeActionChecksum(requestAction, uniqueID string, currentAction interface{}) string {
	raw, _ := json.Marshal(currentAction)

	var m map[string]interface{}
	if json.Unmarshal(raw, &m) == nil {
		for k := range volatileChecksumKeys {
			delete(m, k)
		}
		raw, _ = json.Marshal(m)
	}

	h := sha256.New()
	h.Write([]byte(requestAction))
	h.Write([]byte("|"))
	h.Write([]byte(uniqueID))
	h.Write([]byte("|"))
	h.Write(raw)
	return hex.EncodeToString(h.Sum(nil))
}

func extractLevelClaimPairs(filters map[string]interface{}) []imodel.LevelClaimPair {
	v, ok := filters["level_claim_pairs"]
	if !ok {
		return nil
	}
	if pairs, ok := v.([]imodel.LevelClaimPair); ok {
		return pairs
	}
	return nil
}

// extractCheckerLevelStatusPairs extracts checker level status pairs from filters.
// Looks for filters like "checker_level_1_status", "checker_level_2_status" etc.
// Also supports "checker_level_statuses" which should be []LevelClaimPair.
func extractCheckerLevelStatusPairs(filters map[string]interface{}) []imodel.LevelClaimPair {
	// First check if pre-parsed pairs exist
	if v, ok := filters["checker_level_statuses"]; ok {
		if pairs, ok := v.([]imodel.LevelClaimPair); ok {
			return pairs
		}
	}

	var pairs []imodel.LevelClaimPair

	// Look for checker_level_X_status pattern in filters
	for key, val := range filters {
		// Match pattern: checker_level_{number}_status
		var levelNum string
		if n, found := strings.CutPrefix(key, "checker_level_"); found {
			if n, found = strings.CutSuffix(n, "_status"); found {
				levelNum = n
			}
		}
		if levelNum == "" {
			continue
		}

		status, ok := val.(string)
		if !ok {
			continue
		}

		pairs = append(pairs, imodel.LevelClaimPair{
			Level: levelNum,
			Claim: strings.ToUpper(status),
		})
	}

	return pairs
}

func CpsActionCSVHeader(fields []string) []string {
	return lib.CPSActionHeadersFromFields(fields)
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
