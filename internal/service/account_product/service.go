package account_product_service

import (
	"context"
	"errors"
	"path"
	"strings"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	ap_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/account_product"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	ap_core "github.com/abelshebela/cbe-super-app-cps-action/internal/service/account_product/core"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type accountProductService struct {
	repo         storage.AccountProductRepository
	categoryRepo storage.AccountProductCategoryRepository
	cpsService   service.CPSActionService
	logger       utils.Logger
	minio        *s3.Client
	bucketName   string
	cfg          *config.VaultConfig
}

var _ service.AccountProductService = (*accountProductService)(nil)

func NewAccountProductService(
	repo storage.AccountProductRepository,
	categoryRepo storage.AccountProductCategoryRepository,
	cpsService service.CPSActionService,
	logger utils.Logger,
	minio *s3.Client,
	bucketName string,
	cfg *config.VaultConfig,
) service.AccountProductService {
	return &accountProductService{
		repo:         repo,
		categoryRepo: categoryRepo,
		cpsService:   cpsService,
		logger:       logger,
		minio:        minio,
		bucketName:   bucketName,
		cfg:          cfg,
	}
}

func (s *accountProductService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "AccountProduct", "Authorize")
	defer span.End()

	ap, err := ap_core.MapFromAction(cpsAction.CurrentAction)
	if err != nil {
		log.Errorf("[APSvc][Authorize] map err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	switch constants.RequestAction(cpsAction.RequestAction) {
	case constants.RequestCreateAccountProduct:

		if err := s.repo.Create(ctx, &ap); err != nil {
			span.AddEvent("create failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[APSvc][Authorize] create err: %v", err)
			return nil, err
		}
		log.Infof("[APSvc][Authorize] created")

	case constants.RequestUpdateAccountProduct:
		if err := s.repo.Update(ctx, cpsAction.UniqueId, &ap); err != nil {
			span.AddEvent("update failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[APSvc][Authorize] update err: %v", err)
			return nil, err
		}
		log.Infof("[APSvc][Authorize] updated id=%s", cpsAction.UniqueId)

	case constants.RequestDeleteAccountProduct:
		if err := s.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			span.AddEvent("delete failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[APSvc][Authorize] delete err: %v", err)
			return nil, err
		}
		log.Infof("[APSvc][Authorize] deleted id=%s", cpsAction.UniqueId)

	case constants.RequestEnableAccountProduct:
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			log.Errorf("[APSvc][Authorize] enable err: %v", err)
			return nil, err
		}

	case constants.RequestDisableAccountProduct:
		if err := s.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			log.Errorf("[APSvc][Authorize] disable err: %v", err)
			return nil, err
		}

	default:
		log.Errorf("[APSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	log.Infof("[APSvc][Authorize] done action=%s", cpsAction.RequestAction)
	return cpsAction, nil
}

func (s *accountProductService) GetAll(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]imodel.AccountProduct], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAll", "AccountProduct", "GetAll")
	defer span.End()

	result, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("fetch failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[APSvc][GetAll] err: %v", err)
		return nil, err
	}
	log.Infof("[APSvc][GetAll] count=%d", len(result.Data))
	return result, nil
}

func (s *accountProductService) GetByID(ctx context.Context, id string) (*imodel.AccountProduct, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "GetByID", "AccountProduct", "GetByID")
	defer span.End()

	result, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("fetch failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[APSvc][GetByID] id=%s err=%v", id, err)
		return nil, err
	}
	log.Infof("[APSvc][GetByID] found id=%s", id)
	return result, nil
}

func (s *accountProductService) Create(ctx context.Context, req ap_dto.CreateAPRequest) (*imodel.AccountProduct, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "Create", "AccountProduct", "Create")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[APSvc][Create] incomplete user")
		return nil, errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	//   bypass for test
	existing, err := s.repo.FindByCBSCode(ctx, req.CBSProductCode)
	if err != nil && !strings.Contains(err.Error(), localization.ErrorResourceNotFound.Code) {
		log.Errorf("[APSvc][Create] dup check err: %v", err)
		return nil, err
	}
	if existing != nil {
		log.Errorf("[APSvc][Create] duplicate cbs_product_code=%s", req.CBSProductCode)
		return nil, errors.New(localization.ErrorAPAlreadyExists.Code)
	}

	// iconURL, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Icon,
	// 	string(constants.AccountProductFolderName), *s.cfg, "", s.logger)
	// if err != nil {
	// 	span.AddEvent("icon upload failed", trace.WithAttributes(attribute.String("error", err.Error())))
	// 	log.Errorf("[APSvc][Create] upload icon err: %v", err)
	// 	return nil, errors.New(localization.ErrorUnhandledServer.Code)
	// }

	coverImageURL, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.CoverImage,
		string(constants.AccountProductFolderName), *s.cfg, "", s.logger)
	if err != nil {
		span.AddEvent("cover image upload failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[APSvc][Create] upload cover image err: %v", err)
		return nil, errors.New(localization.ErrorUnhandledServer.Code)
	}

	log.Infof("[AccountProduct][Create] Cover Image: %v", coverImageURL)

	payload := imodel.AccountProduct{
		CBSProductCode:          req.CBSProductCode,
		ProductName:             req.ProductName,
		ProductTagLine:          req.ProductTagLine,
		AccountCategoryID:       req.AccountCategoryID,
		AccountCurrency:         strings.ToUpper(req.AccountCurrency),
		MinimumOpeningBalance:   req.MinimumOpeningBalance,
		MinimumMaintenanceFee:   req.MinimumMaintenanceFee,
		InterestFee:             req.InterestRate,
		FaqURL:                  req.FaqURL,
		ProductFeatures:         req.ProductFeatures,
		HasPhysicalCard:         req.HasPhysicalCard,
		HasVirtualCard:          req.HasVirtualCard,
		InterestRate:            req.InterestRate,
		IsAvailableForOnbording: req.IsAvailableForOnbording,
		// ProductIcon:           iconURL,
		ProductCoverImage: coverImageURL,
		IsEnabled:         true,
		Eligibility: req.Eligibility,
		ProductDescription: req.ProductDescription,
	}
	//bypass 2
	if s.categoryRepo != nil {
		if cat, err := s.categoryRepo.FindByID(ctx, req.AccountCategoryID); err == nil {
			payload.CategoryName = cat.CategoryName
			payload.CBSCategoryCode = cat.CBSCategoryCode
			payload.AccountType = cat.AccountType
		} else {
			log.Errorf("[APSvc][Create] category lookup err: %v", err)
		}
	}

	action := lib.CpsModelBuilder("", makerData, nil, payload,
		string(constants.RequestCreateAccountProduct), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("cps action failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[APSvc][Create] cps err: %v", err)
		return nil, err
	}
	log.Infof("[APSvc][Create] request created code=%s", req.CBSProductCode)
	return &payload, nil
}

func (s *accountProductService) Update(ctx context.Context, id string, req ap_dto.UpdateAPRequest) (*imodel.AccountProduct, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "Update", "AccountProduct", "Update")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return nil, errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[APSvc][Update] find err: %v", err)
		return nil, err
	}

	updated := *existing

	if req.CBSProductCode != "" && req.CBSProductCode != existing.CBSProductCode {
		dup, _ := s.repo.FindByCBSCode(ctx, req.CBSProductCode)
		if dup != nil {
			return nil, errors.New(localization.ErrorAPAlreadyExists.Code)
		}
		updated.CBSProductCode = req.CBSProductCode
	}
	if req.ProductName != "" {
		updated.ProductName = req.ProductName
	}
	if req.ProductTagLine != "" {
		updated.ProductTagLine = req.ProductTagLine
	}
	if req.Eligibility != ""{
		updated.Eligibility =req.Eligibility
	}
	if req.AccountCategoryID != "" {
		updated.AccountCategoryID = req.AccountCategoryID
		if s.categoryRepo != nil {
			if cat, err := s.categoryRepo.FindByID(ctx, req.AccountCategoryID); err == nil {
				updated.CategoryName = cat.CategoryName
				updated.CBSCategoryCode = cat.CBSCategoryCode
				updated.AccountType = cat.AccountType
			} else {
				log.Errorf("[APSvc][Update] category lookup err: %v", err)
			}
		}
	}
	if req.AccountCurrency != "" {
		updated.AccountCurrency = strings.ToUpper(req.AccountCurrency)
	}
	if req.MinimumOpeningBalance > 0 {
		updated.MinimumOpeningBalance = req.MinimumOpeningBalance
	}
	if req.MinimumMaintenanceFee > 0 {
		updated.MinimumMaintenanceFee = req.MinimumMaintenanceFee
	}
	if req.InterestRate > 0 {
		updated.InterestFee = req.InterestRate
	}
	if req.FaqURL != "" {
		updated.FaqURL = req.FaqURL
	}
	if req.ProductFeatures[0] != "" {
		updated.ProductFeatures = req.ProductFeatures
	}
	if req.HasPhysicalCard != nil {
		updated.HasPhysicalCard = *req.HasPhysicalCard
	}
	if req.HasVirtualCard != nil {
		updated.HasVirtualCard = *req.HasVirtualCard
	}
	if req.IsAvailableForOnbording != nil {
		updated.IsAvailableForOnbording = *req.IsAvailableForOnbording
	}

	if req.CoverImage != nil {
		var objectKey string
		if existing.ProductCoverImage != "" {
			objectKey = path.Base(existing.ProductCoverImage)
		}
		coverImageURL, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.CoverImage,
			string(constants.AccountProductFolderName), *s.cfg, objectKey, s.logger)
		if err != nil {
			span.AddEvent("cover image upload failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[APSvc][Update] upload cover image err: %v", err)
			return nil, errors.New(localization.ErrorUnhandledServer.Code)
		}
		updated.ProductCoverImage = coverImageURL
	}
	action := lib.CpsModelBuilder(id, makerData, existing, updated,
		string(constants.RequestUpdateAccountProduct), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("cps action failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[APSvc][Update] cps err: %v", err)
		return nil, err
	}
	log.Infof("[APSvc][Update] request created id=%s", id)
	return &updated, nil
}

func (s *accountProductService) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "Delete", "AccountProduct", "Delete")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[APSvc][Delete] find err: %v", err)
		return err
	}

	updated := *existing
	updated.IsDeleted = true

	action := lib.CpsModelBuilder(id, makerData, existing, updated,
		string(constants.RequestDeleteAccountProduct), constants.DELETE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("cps action failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[APSvc][Delete] cps err: %v", err)
		return err
	}
	log.Infof("[APSvc][Delete] request created id=%s", id)
	return nil
}

func (s *accountProductService) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisable", "AccountProduct", "EnableOrDisable")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[APSvc][EnableOrDisable] find err: %v", err)
		return err
	}

	if existing.IsEnabled == enable {
		if enable {
			return errors.New(localization.ErrorAPAlreadyEnabled.Code)
		}
		return errors.New(localization.ErrorAPAlreadyDisabled.Code)
	}

	updated := *existing
	updated.IsEnabled = enable

	requestAction := constants.RequestEnableAccountProduct
	if !enable {
		requestAction = constants.RequestDisableAccountProduct
	}

	action := lib.CpsModelBuilder(id, makerData, existing, updated, string(requestAction), constants.UPDATE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("cps action failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[APSvc][EnableOrDisable] cps err: %v", err)
		return err
	}
	log.Infof("[APSvc][EnableOrDisable] request created id=%s enable=%v", id, enable)
	return nil
}
