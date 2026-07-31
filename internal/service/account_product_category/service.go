package account_product_category_service

import (
	"context"
	"errors"
	"strings"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	apc_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/account_product_category"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	apc_core "github.com/abelshebela/cbe-super-app-cps-action/internal/service/account_product_category/core"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type accountProductCategoryService struct {
	repo       storage.AccountProductCategoryRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

var _ service.AccountProductCategoryService = (*accountProductCategoryService)(nil)

func NewAccountProductCategoryService(
	repo storage.AccountProductCategoryRepository,
	cpsService service.CPSActionService,
	logger utils.Logger,
) service.AccountProductCategoryService {
	return &accountProductCategoryService{
		repo:       repo,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (s *accountProductCategoryService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "AccountProductCategory", "Authorize")
	defer span.End()

	apc, err := apc_core.MapFromAction(cpsAction.CurrentAction)
	if err != nil {
		log.Errorf("[APCSvc][Authorize] map action err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	switch constants.RequestAction(cpsAction.RequestAction) {
	case constants.RequestCreateAccountProductCategory:
		if err := s.repo.Create(ctx, &apc); err != nil {
			span.AddEvent("create failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[APCSvc][Authorize] create err: %v", err)
			return nil, err
		}
		log.Infof("[APCSvc][Authorize] created")

	case constants.RequestUpdateAccountProductCategory:
		if err := s.repo.Update(ctx, cpsAction.UniqueId, &apc); err != nil {
			span.AddEvent("update failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[APCSvc][Authorize] update err: %v", err)
			return nil, err
		}
		log.Infof("[APCSvc][Authorize] updated id=%s", cpsAction.UniqueId)

	case constants.RequestDeleteAccountProductCategory:
		if err := s.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			span.AddEvent("delete failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[APCSvc][Authorize] delete err: %v", err)
			return nil, err
		}
		log.Infof("[APCSvc][Authorize] deleted id=%s", cpsAction.UniqueId)

	case constants.RequestEnableAccountProductCategory:
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			log.Errorf("[APCSvc][Authorize] enable err: %v", err)
			return nil, err
		}
	case constants.RequestDisableAccountProductCategory:
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			log.Errorf("[APCSvc][Authorize] disable err: %v", err)
			return nil, err
		}

	default:
		log.Errorf("[APCSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	log.Infof("[APCSvc][Authorize] done action=%s", cpsAction.RequestAction)
	return cpsAction, nil
}

func (s *accountProductCategoryService) GetAll(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]imodel.AccountProductCategory], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAll", "AccountProductCategory", "GetAll")
	defer span.End()

	result, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("fetch failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[APCSvc][GetAll] err: %v", err)
		return nil, err
	}
	log.Infof("[APCSvc][GetAll] count=%d", len(result.Data))
	return result, nil
}

func (s *accountProductCategoryService) GetByID(ctx context.Context, id string) (*imodel.AccountProductCategory, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "GetByID", "AccountProductCategory", "GetByID")
	defer span.End()

	result, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("fetch failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[APCSvc][GetByID] id=%s err=%v", id, err)
		return nil, err
	}
	log.Infof("[APCSvc][GetByID] found id=%s", id)
	return result, nil
}

func (s *accountProductCategoryService) Create(ctx context.Context, req apc_dto.CreateAPCRequest) (*imodel.AccountProductCategory, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "Create", "AccountProductCategory", "Create")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[APCSvc][Create] incomplete user")
		return nil, errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByCBSCode(ctx, req.CBSCategoryCode)
	if err != nil && !strings.Contains(err.Error(), localization.ErrorResourceNotFound.Code) {
		log.Errorf("[APCSvc][Create] dup check err: %v", err)
		return nil, err
	}
	if existing != nil {
		log.Errorf("[APCSvc][Create] duplicate cbs_category_code=%s", req.CBSCategoryCode)
		return nil, errors.New(localization.ErrorAPCAlreadyExists.Code)
	}

	existingByName, err := s.repo.FindByCategoryName(ctx, req.CategoryName)
	if err != nil && !strings.Contains(err.Error(), localization.ErrorResourceNotFound.Code) {
		log.Errorf("[APCSvc][Create] name dup check err: %v", err)
		return nil, err
	}
	if existingByName != nil {
		log.Errorf("[APCSvc][Create] duplicate category_name=%s", req.CategoryName)
		return nil, errors.New(localization.ErrorAPCAlreadyExists.Code)
	}

	payload := imodel.AccountProductCategory{
		AccountType:     strings.ToUpper(req.ProductLine),
		CategoryName:    req.CategoryName,
		CBSCategoryCode: req.CBSCategoryCode,
		Description:     req.Description,
		IsEnabled:       true,
	}

	action := lib.CpsModelBuilder("", makerData, nil, payload,
		string(constants.RequestCreateAccountProductCategory), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("cps action failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[APCSvc][Create] cps err: %v", err)
		return nil, err
	}
	log.Infof("[APCSvc][Create] request created code=%s", req.CBSCategoryCode)
	return &payload, nil
}

func (s *accountProductCategoryService) Update(ctx context.Context, id string, req apc_dto.UpdateAPCRequest) (*imodel.AccountProductCategory, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "Update", "AccountProductCategory", "Update")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[APCSvc][Update] incomplete user")
		return nil, errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("fetch failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[APCSvc][Update] find err: %v", err)
		return nil, err
	}

	updated := *existing
	if req.ProductLine != "" {
		updated.AccountType = strings.ToUpper(req.ProductLine)
	}
	if req.CategoryName != "" {
		if !strings.EqualFold(req.CategoryName, existing.CategoryName) {
			dupByName, _ := s.repo.FindByCategoryName(ctx, req.CategoryName)
			if dupByName != nil {
				return nil, errors.New(localization.ErrorAPCAlreadyExists.Code)
			}
		}
		updated.CategoryName = req.CategoryName
	}
	if req.CBSCategoryCode != "" {
		if req.CBSCategoryCode != existing.CBSCategoryCode {
			dup, _ := s.repo.FindByCBSCode(ctx, req.CBSCategoryCode)
			if dup != nil {
				return nil, errors.New(localization.ErrorAPCAlreadyExists.Code)
			}
		}
		updated.CBSCategoryCode = req.CBSCategoryCode
	}
	if req.Description != "" {
		updated.Description = req.Description
	}

	action := lib.CpsModelBuilder(id, makerData, existing, updated,
		string(constants.RequestUpdateAccountProductCategory), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("cps action failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[APCSvc][Update] cps err: %v", err)
		return nil, err
	}
	log.Infof("[APCSvc][Update] request created id=%s", id)
	return &updated, nil
}

func (s *accountProductCategoryService) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "Delete", "AccountProductCategory", "Delete")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[APCSvc][Delete] find err: %v", err)
		return err
	}

	updated := *existing
	updated.IsDeleted = true

	action := lib.CpsModelBuilder(id, makerData, existing, updated,
		string(constants.RequestDeleteAccountProductCategory), constants.DELETE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("cps action failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[APCSvc][Delete] cps err: %v", err)
		return err
	}
	log.Infof("[APCSvc][Delete] request created id=%s", id)
	return nil
}

func (s *accountProductCategoryService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisable", "AccountProductCategory", "EnableOrDisable")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[APCSvc][EnableOrDisable] find err: %v", err)
		return err
	}

	if existing.IsEnabled == enable {
		if enable {
			return errors.New(localization.ErrorAPCAlreadyEnabled.Code)
		}
		return errors.New(localization.ErrorAPCAlreadyDisabled.Code)
	}

	updated := *existing
	updated.IsEnabled = enable

	requestAction := constants.RequestEnableAccountProductCategory
	if !enable {
		requestAction = constants.RequestDisableAccountProductCategory
	}

	action := lib.CpsModelBuilder(id, makerData, existing, updated, string(requestAction), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("cps action failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[APCSvc][EnableOrDisable] cps err: %v", err)
		return err
	}
	log.Infof("[APCSvc][EnableOrDisable] request created id=%s enable=%v", id, enable)
	return nil
}
