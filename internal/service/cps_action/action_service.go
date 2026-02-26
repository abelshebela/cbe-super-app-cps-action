package cpsaction

import (
	"bytes"
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	actionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/minio/minio-go/v7"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type cpsActionService struct {
	repo       storage.CPSActionRepository
	roles      storage.CPSActionRoleRepository
	logger     utils.Logger
	dispatcher Dispatcher
	minioClient *s3.Client
	buckerName string
}

// IsMakerOnlyForRequest returns true if the module mapped from requestAction is configured as maker-only in CPSActionRole.
func (ca *cpsActionService) IsMakerOnlyForRequest(ctx context.Context, requestAction string) (bool, error) {
	if mod, ok := ResolveModuleForRA(RequestAction(requestAction)); ok && ca.roles != nil {
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
	_, err = ca.repo.Update(ctx, actionCode, upd)
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

func NewCPSActionService(roles storage.CPSActionRoleRepository, repo storage.CPSActionRepository, logger utils.Logger, dispatcher Dispatcher,minioClient *s3.Client,bucketName string) service.CPSActionService {
	return &cpsActionService{
		repo:       repo,
		logger:     logger,
		roles:      roles,
		dispatcher: dispatcher,
		minioClient: minioClient,
		buckerName: bucketName,

	}
}

func (ca *cpsActionService) CreateCPSAction(ctx context.Context, cpsAction *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateCPSAction", "CPSAction", "CreateCPSAction")
	defer span.End()
	var existing *model.CPSAction
	var err error
	ca.logger.Infof("[CpsActionSvc][Create] action: %s", cpsAction.RequestAction)

	roleCode := ctx.Value(constants.ContextKey("role_code")).(string)
	actionName, _ := ctx.Value(constants.ContextKey("action_name")).(string)

	reqs := ca.pendingLockRequestActions(actionName, cpsAction.RequestAction)
	if len(reqs) > 0 {
		existing, err = ca.GetPendingCPSActionByRoleAndRequestActions(ctx, roleCode, reqs)
		if err != nil && err.Error() != localization.ErrorActionNotFound.Code {
			span.AddEvent("failed to get cps action by role and request actions", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
	}

	if existing != nil {
		span.AddEvent("pending cps action exists", trace.WithAttributes(attribute.String("error", "pending cps action exists")))
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

	defaultReq := []string{normalize(requestAction)}
	if actionName == "" {
		return defaultReq
	}

	lst, ok := RequestActionGroups[actionName]
	if !ok {
		return defaultReq
	}

	wantCreateOnly := isCreate(requestAction)
	seen := map[string]struct{}{}
	reqs := make([]string, 0, len(lst))

	for _, ra := range lst {
		key := string(ra)
		if _, ok := seen[key]; ok {
			continue
		}

		isKeyCreate := isCreate(key)
		if wantCreateOnly != isKeyCreate {
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
	ctx, span := local_util.TraceLogger(ctx, "service", "ApproveCPSAction", "CPSAction", "ApproveCPSAction")
	defer span.End()
	ca.logger.Infof("[CpsActionSvc][Approve] action: %s", action.ActionCode)

	data, err := ca.repo.Update(ctx, action.ActionCode, *action)
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
	_, err := ca.repo.Update(ctx, action_code, *action)
	if err != nil {
		span.AddEvent("failed to update cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}
func (ca *cpsActionService) CancelCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CancelCPSAction", "CPSAction", "CancelCPSAction")
	defer span.End()
	_, err := ca.repo.Update(ctx, action_code, *action)
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

func (ca *cpsActionService) GetCPSActionsForApprover(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionsForApprover", "CPSAction", "GetCPSActionsForApprover")
	defer span.End()
	result, err := ca.repo.SanitizedFindAllWithPaginationForApprover(ctx, userID, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ca *cpsActionService) GetCPSActionsForAuditor(ctx context.Context, userID string, RAList []string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCPSActionsForAuditor", "CPSAction", "GetCPSActionsForAuditor")
	defer span.End()
	result, err := ca.repo.SanitizedFindAllWithPaginationForAuditor(ctx, userID, *filterParams, RAList)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
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

func (ca *cpsActionService) GetPendingCPSActionByRoleAndRequestActions(ctx context.Context, roleCode string, requestActions []string) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPendingCPSActionByRoleAndRequestActions", "CPSAction", "GetPendingCPSActionByRoleAndRequestActions")
	defer span.End()

	filter := bson.M{
		"role_code":      roleCode,
		"action_status":  string(constants.Pending),
		"request_action": bson.M{"$in": requestActions},
	}

	action, err := ca.repo.SanitizedFindOne(ctx, filter)
	if err != nil {
		span.AddEvent("failed to find one", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
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
	if mod, ok := ResolveModuleForRA(RequestAction(requestAction)); ok && ca.roles != nil {
		if approver, err := ca.roles.FindApproverByActionName(ctx, strings.ToUpper(mod), roleCode); err == nil {
			approverData = approver
		}
	} else {
		return imodel.CPSActionApproveIndex{}, errors.New(localization.ErrorOperationNotAllowed.Message)
	}

	return approverData, nil
}

func (ca *cpsActionService) GetUserCreatedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetUserCreatedActions", "CPSAction", "GetUserCreatedActions")
	defer span.End()
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	filterParams.Filters["maker_id"] = userID

	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")
	if err != nil {
		span.AddEvent("failed to find pending cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ca *cpsActionService) GetUserCheckedActions(ctx context.Context, userID string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetUserCheckedActions", "CPSAction", "GetUserCheckedActions")
	defer span.End()
	if filterParams == nil {
		filterParams = &types.Filter{}
	}
	if filterParams.Filters == nil {
		filterParams.Filters = map[string]interface{}{}
	}
	filterParams.Filters["maker_id"] = userID

	result, err := ca.repo.SanitizedFindAllWithPagination(ctx, *filterParams, "")
	if err != nil {
		span.AddEvent("failed to find approver cps actions by user", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	return result, nil
}

func (ca *cpsActionService) ExportCpsActionData(
    ctx context.Context,
    startDate, endDate time.Time,export_type string,
) (string, error) {

    if endDate.Before(startDate) {
        return "", errors.New("end_date cannot be before start_date")
    }

    // 1 Create temp file
    tmpFile, err := os.CreateTemp("", "cps_actions_*.csv")
    if err != nil {
        return "", fmt.Errorf("create temp file: %w", err)
    }
    defer os.Remove(tmpFile.Name())
    defer tmpFile.Close()

    writer := csv.NewWriter(tmpFile)

    // 2️ Write Header
    if err := writer.Write(CpsActionCSVHeader()); err != nil {
        return "", fmt.Errorf("write header: %w", err)
    }

    // 3️Stream from repository
    err = ca.repo.StreamByDateRange(ctx, startDate, endDate,
        func(action *model.CPSAction) error {
            return ca.processCPSAction(writer, action)
        },
    )
    if err != nil {
        return "", fmt.Errorf("stream data: %w", err)
    }

    writer.Flush()
    if err := writer.Error(); err != nil {
        return "", fmt.Errorf("flush csv: %w", err)
    }

    // 4️Upload to MinIO
    objectName := fmt.Sprintf(
        "exports/cps-actions/cps_actions_%s_to_%s_%d.csv",
        startDate.Format("20060102"),
        endDate.Format("20060102"),
        time.Now().Unix(),
    )

    link,err := UploadFileToMinio(ctx, tmpFile.Name(), objectName)
	if err != nil {
        return "", err
    }

    // 5️⃣ Generate presigned URL
    presignedURL, err := ca.minioClient.PresignedGetObject(
        ctx,
        ca.minioBucket,
        objectName,
        15*time.Minute,
        nil,
    )
    if err != nil {
        return "", fmt.Errorf("generate presigned url: %w", err)
    }

    return presignedURL.String(), nil
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

    checkerJSON, _ := json.Marshal(a.CheckerUsers)
    auditorJSON, _ := json.Marshal(a.AuditorUsers)
    prevJSON, _ := json.Marshal(a.PreviousAction)
    currJSON, _ := json.Marshal(a.CurrentAction)

    return []string{
        a.ID.Hex(),
        a.ActionCode,
        a.UniqueId,
        a.MakerID,
        a.MakerName,
        a.MakerPhoneNumber,
        string(checkerJSON),
        string(auditorJSON),
        strconv.Itoa(int(a.AuditorCount)),
        string(a.AuditorStatus),
        fmt.Sprintf("%f", a.CurrentAuditorIndex),
        strconv.Itoa(int(a.CheckerCount)),
        fmt.Sprintf("%f", a.CurrentCheckerIndex),
        a.RoleCode,
        a.RejectionReason,
        a.CanceledReason,
        string(prevJSON),
        string(currJSON),
        a.ActionStatus,
        a.ActionType,
        strconv.FormatBool(a.IsDeleted),
        a.RequestAction,
        strconv.FormatInt(a.Version, 10),
        a.ReversedByRoleID,
        a.ReversedByID,
        a.ReversedByName,
        // formatTime(a.ReversedAt),
        // formatTime(a.CreatedAt),
        // formatTime(a.LastModifiedAt),
        // formatTime(a.MakerActionTime),
    }, nil
}

func CpsActionCSVHeader() []string {
    return []string{
        "ID",
        "ActionCode",
        "UniqueId",
        "MakerID",
        "MakerName",
        "MakerPhoneNumber",
        "CheckerUsers",
        "AuditorUsers",
        "AuditorCount",
        "AuditorStatus",
        "CurrentAuditorIndex",
        "CheckerCount",
        "CurrentCheckerIndex",
        "RoleCode",
        "RejectionReason",
        "CanceledReason",
        "PreviousAction",
        "CurrentAction",
        "ActionStatus",
        "ActionType",
        "IsDeleted",
        "RequestAction",
        "Version",
        "ReversedByRoleID",
        "ReversedByID",
        "ReversedByName",
        "ReversedAt",
        "CreatedAt",
        "LastModifiedAt",
        "MakerActionTime",
    }
}

func  UploadFileToMinio(
    ctx context.Context,
	s3Client *s3.Client,
	bucketName string,
	fileHeader *multipart.FileHeader,
	env config.VaultConfig,
    filePath string,
    objectKey string,
) (string,error) {

    file, err := os.Open(filePath)
    if err != nil {
        return "",fmt.Errorf("open file: %w", err)
    }
    defer file.Close()

    stat, err := file.Stat()
    if err != nil {
        return "",err
    }
	size := stat.Size()
	contentType := "text/csv"
	putInput := &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		Body:          bytes.NewReader(finalBytes),
		ContentType:   aws.String(contentType),
		ContentLength: &size,
	}
    _, err = s3Client.PutObject(
        ctx,
        putInput,
    )
    if err != nil {
        return "",fmt.Errorf("upload to minio: %w", err)
    }
	 baseURL := env.MinioPublicEndPoint
	url := fmt.Sprintf("%s/%s", baseURL, strings.TrimPrefix(objectKey, "/")) 
	   return url,nil
}


func  UploadCSVFile(
    ctx context.Context,
	s3Client *s3.Client,
    filePath string,
    objectKey string,
	bucketName string,
) (string, error) {

    // 1. Open file (NO memory loading)
    file, err := os.Open(filePath)
    if err != nil {
        return "", fmt.Errorf("open csv file: %w", err)
    }
    defer file.Close()

    // 2. Get file size (required by MinIO)
    stat, err := file.Stat()
    if err != nil {
        return "", fmt.Errorf("stat file: %w", err)
    }

	byteFile := []byte(file)

	minioInputObject := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:aws.String(objectKey),
		Body: bytes.NewReader(file),
	}
    // 3. Upload using STREAM (file reader)
    _, err = s3Client.PutObject(
        ctx,
        bucketName,
        objectKey,
        file,              //  STREAM directly (important)
        stat.Size(),       // file size
        minio.PutObjectOptions{
            ContentType: "text/csv",
        },
    )
    if err != nil {
        return "", fmt.Errorf("upload csv to minio: %w", err)
    }

    // 4. Generate presigned URL (download link)
    presignedURL, err := m.client.PresignedGetObject(
        ctx,
        m.bucketName,
        objectKey,
        15*time.Minute,
        nil,
    )
    if err != nil {
        return "", fmt.Errorf("generate presigned url: %w", err)
    }

    return presignedURL.String(), nil
}