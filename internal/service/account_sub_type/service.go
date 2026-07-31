package account_sub_type_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	account_sub_type_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/account_sub_type"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type accountSubTypeService struct {
	oracleRepo storage.AccountSubTypeOracleRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewAccountSubTypeService(oracleRepo storage.AccountSubTypeOracleRepository, cpsService service.CPSActionService, logger utils.Logger) service.AccountSubTypeService {
	return &accountSubTypeService{
		oracleRepo: oracleRepo,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (s *accountSubTypeService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "AccountSubType", "Authorize")
	defer span.End()

	var actionMap interface{}
	marshaled, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("failed to marshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		log.Errorf("[AccountSubTypeSvc][Authorize] marshal err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if err = json.Unmarshal(marshaled, &actionMap); err != nil {
		span.AddEvent("failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		log.Errorf("[AccountSubTypeSvc][Authorize] unmarshal err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	actionData := mapAccountSubTypeFromAction(actionMap.(map[string]interface{}))

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateAccountSubType):
		err = s.oracleRepo.Create(ctx, &actionData)
		if err != nil {
			span.AddEvent("[Authorize] create failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[AccountSubTypeSvc][Authorize] create err: %v", err)
			return nil, err
		}
		log.Infof("[AccountSubTypeSvc][Authorize] created")

	case string(constants.RequestUpdateAccountSubType):
		err = s.oracleRepo.Update(ctx, cpsAction.UniqueId, &actionData)
		if err != nil {
			span.AddEvent("[Authorize] update failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[AccountSubTypeSvc][Authorize] update err: %v", err)
			return nil, err
		}
		log.Infof("[AccountSubTypeSvc][Authorize] updated id: %s", cpsAction.UniqueId)

	case string(constants.RequestDeleteAccountSubType):
		err = s.oracleRepo.Delete(ctx, cpsAction.UniqueId)
		if err != nil {
			span.AddEvent("[Authorize] delete failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[AccountSubTypeSvc][Authorize] delete err: %v", err)
			return nil, err
		}
		log.Infof("[AccountSubTypeSvc][Authorize] deleted id: %s", cpsAction.UniqueId)

	case string(constants.RequestEnableAccountSubType):
		err = s.oracleRepo.EnableOrDisable(ctx, cpsAction.UniqueId, true)
		if err != nil {
			span.AddEvent("[Authorize] enable failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[AccountSubTypeSvc][Authorize] enable err: %v", err)
			return nil, err
		}
		log.Infof("[AccountSubTypeSvc][Authorize] enabled id: %s", cpsAction.UniqueId)

	case string(constants.RequestDisableAccountSubType):
		err = s.oracleRepo.EnableOrDisable(ctx, cpsAction.UniqueId, false)
		if err != nil {
			span.AddEvent("[Authorize] disable failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[AccountSubTypeSvc][Authorize] disable err: %v", err)
			return nil, err
		}
		log.Infof("[AccountSubTypeSvc][Authorize] disabled id: %s", cpsAction.UniqueId)

	default:
		span.AddEvent("[Authorize] unsupported action", trace.WithAttributes(attribute.String("action", cpsAction.RequestAction)))
		log.Errorf("[AccountSubTypeSvc][Authorize] unsupported: %s", cpsAction.RequestAction)
		return nil, fmt.Errorf("%s", localization.ErrorUnsupportedAction.Code)
	}

	log.Infof("[AccountSubTypeSvc][Authorize] done: %s", cpsAction.RequestAction)
	return cpsAction, nil
}

func (s *accountSubTypeService) CreateOneAccountSubType(ctx context.Context, req account_sub_type_dto.CreateAccountSubTypeRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateOneAccountSubType", "AccountSubType", "CreateOneAccountSubType")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[AccountSubTypeSvc][Create] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.oracleRepo.FindByCodeOrName(ctx, req.AccountSubTypeCode, "")
	if err != nil && !strings.Contains(err.Error(), localization.ErrorResourceNotFound.Code) {
		log.Errorf("[AccountSubTypeSvc][Create] dup check err: %v", err)
		return err
	}
	if existing != nil && existing.ID != "" {
		log.Errorf("[AccountSubTypeSvc][Create] code already exists")
		return fmt.Errorf("%s", localization.ErrorAccountSubTypeCodeAlreadyExists.Code)
	}

	existing, err = s.oracleRepo.FindByCodeOrName(ctx, "", req.AccountSubTypeName)
	if err != nil && !strings.Contains(err.Error(), localization.ErrorResourceNotFound.Code) {
		log.Errorf("[AccountSubTypeSvc][Create] dup name check err: %v", err)
		return err
	}
	if existing != nil && existing.ID != "" {
		log.Errorf("[AccountSubTypeSvc][Create] name already exists")
		return fmt.Errorf("%s", localization.ErrorAccountSubTypeNameAlreadyExists.Code)
	}

	payload := imodel.AccountSubType{
		AccountType:        req.AccountType,
		Gender:             req.Gender,
		AccountSubTypeName: req.AccountSubTypeName,
		AccountSubTypeCode: req.AccountSubTypeCode,
		IsEnabled:          0,
	}

	action := lib.CpsModelBuilder("", makerData, nil, payload, string(constants.RequestCreateAccountSubType), constants.CREATE)
	if err = s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("[Create] cps action err", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[AccountSubTypeSvc][Create] cps action err: %v", err)
		return err
	}

	log.Infof("[AccountSubTypeSvc][Create] request created")
	return nil
}

func (s *accountSubTypeService) UpdateOneAccountSubType(ctx context.Context, id string, req account_sub_type_dto.UpdateAccountSubTypeRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateOneAccountSubType", "AccountSubType", "UpdateOneAccountSubType")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[AccountSubTypeSvc][Update] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	current, err := s.oracleRepo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[AccountSubTypeSvc][Update] find err: %v", err)
		return err
	}

	if req.AccountSubTypeCode != "" && !strings.EqualFold(req.AccountSubTypeCode, current.AccountSubTypeCode) {
		existing, err := s.oracleRepo.FindByCodeOrName(ctx, req.AccountSubTypeCode, "")
		if err != nil && !strings.Contains(err.Error(), localization.ErrorResourceNotFound.Code) {
			log.Errorf("[AccountSubTypeSvc][Update] dup code check err: %v", err)
			return err
		}
		if existing != nil && existing.ID != "" && existing.ID != id {
			log.Errorf("[AccountSubTypeSvc][Update] code already exists")
			return fmt.Errorf("%s", localization.ErrorAccountSubTypeCodeAlreadyExists.Code)
		}
	}

	if req.AccountSubTypeName != "" && !strings.EqualFold(req.AccountSubTypeName, current.AccountSubTypeName) {
		existing, err := s.oracleRepo.FindByCodeOrName(ctx, "", req.AccountSubTypeName)
		if err != nil && !strings.Contains(err.Error(), localization.ErrorResourceNotFound.Code) {
			log.Errorf("[AccountSubTypeSvc][Update] dup name check err: %v", err)
			return err
		}
		if existing != nil && existing.ID != "" && existing.ID != id {
			log.Errorf("[AccountSubTypeSvc][Update] name already exists")
			return fmt.Errorf("%s", localization.ErrorAccountSubTypeNameAlreadyExists.Code)
		}
	}

	updated := *current
	if req.AccountType != "" {
		updated.AccountType = req.AccountType
	}
	if req.Gender != "" {
		updated.Gender = req.Gender
	}
	if req.AccountSubTypeName != "" {
		updated.AccountSubTypeName = req.AccountSubTypeName
	}
	if req.AccountSubTypeCode != "" {
		updated.AccountSubTypeCode = req.AccountSubTypeCode
	}

	action := lib.CpsModelBuilder(id, makerData, current, updated, string(constants.RequestUpdateAccountSubType), constants.UPDATE)
	if err = s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("[Update] cps action err", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[AccountSubTypeSvc][Update] cps action err: %v", err)
		return err
	}

	log.Infof("[AccountSubTypeSvc][Update] request created id: %s", id)
	return nil
}

func (s *accountSubTypeService) DeleteOneAccountSubType(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteOneAccountSubType", "AccountSubType", "DeleteOneAccountSubType")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[AccountSubTypeSvc][Delete] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	current, err := s.oracleRepo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[AccountSubTypeSvc][Delete] find err: %v", err)
		return err
	}

	deleted := *current
	deleted.IsDeleted = true

	action := lib.CpsModelBuilder(id, makerData, current, deleted, string(constants.RequestDeleteAccountSubType), constants.DELETE)
	if err = s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("[Delete] cps action err", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[AccountSubTypeSvc][Delete] cps action err: %v", err)
		return err
	}

	log.Infof("[AccountSubTypeSvc][Delete] request created id: %s", id)
	return nil
}

func (s *accountSubTypeService) EnableOrDisableAccountSubType(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableAccountSubType", "AccountSubType", "EnableOrDisableAccountSubType")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[AccountSubTypeSvc][EnableOrDisable] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	current, err := s.oracleRepo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[AccountSubTypeSvc][EnableOrDisable] find err: %v", err)
		return err
	}

	if current.IsEnabled == 1 && enable {
		log.Errorf("[AccountSubTypeSvc][EnableOrDisable] already enabled id: %s", id)
		return fmt.Errorf("%s", localization.ErrorAccountSubTypeAlreadyEnabled.Code)
	}
	if current.IsEnabled == 0 && !enable {
		log.Errorf("[AccountSubTypeSvc][EnableOrDisable] already disabled id: %s", id)
		return fmt.Errorf("%s", localization.ErrorAccountSubTypeAlreadyDisabled.Code)
	}

	updated := *current
	var requestAction string
	if enable {
		updated.IsEnabled = 1
		requestAction = string(constants.RequestEnableAccountSubType)
	} else {
		updated.IsEnabled = 0
		requestAction = string(constants.RequestDisableAccountSubType)
	}

	action := lib.CpsModelBuilder(id, makerData, current, updated, requestAction, constants.UPDATE)
	if err = s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("[EnableOrDisable] cps action err", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[AccountSubTypeSvc][EnableOrDisable] cps action err: %v", err)
		return err
	}

	log.Infof("[AccountSubTypeSvc][EnableOrDisable] request created id: %s, enable: %v", id, enable)
	return nil
}

func (s *accountSubTypeService) GetAllAccountSubTypes(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]imodel.AccountSubType], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllAccountSubTypes", "AccountSubType", "GetAllAccountSubTypes")
	defer span.End()

	result, err := s.oracleRepo.FindAllWithPagination(ctx, filterParams)
	if err != nil {
		span.AddEvent("[GetAll] fetch failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[AccountSubTypeSvc][GetAll] fetch err: %v", err)
		return nil, err
	}

	log.Infof("[AccountSubTypeSvc][GetAll] count: %d", len(result.Data))
	return result, nil
}

func (s *accountSubTypeService) GetOneAccountSubType(ctx context.Context, id string) (*imodel.AccountSubType, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "GetOneAccountSubType", "AccountSubType", "GetOneAccountSubType")
	defer span.End()

	result, err := s.oracleRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("[GetOne] fetch failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[AccountSubTypeSvc][GetOne] fetch err: %v", err)
		return nil, err
	}

	log.Infof("[AccountSubTypeSvc][GetOne] found id: %s", id)
	return result, nil
}

func mapAccountSubTypeFromAction(m map[string]interface{}) imodel.AccountSubType {
	ast := imodel.AccountSubType{}
	if v, ok := m["account_type"].(string); ok {
		ast.AccountType = v
	}
	if v, ok := m["gender"].(string); ok {
		ast.Gender = v
	}
	if v, ok := m["account_sub_type_name"].(string); ok {
		ast.AccountSubTypeName = v
	}
	if v, ok := m["account_sub_type_code"].(string); ok {
		ast.AccountSubTypeCode = v
	}
	if v, ok := m["is_enabled"].(float64); ok {
		ast.IsEnabled = int(v)
	}
	if v, ok := m["is_deleted"].(bool); ok {
		ast.IsDeleted = v
	}
	return ast
}
