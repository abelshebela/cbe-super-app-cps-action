package term_and_condition_service

import (
	"context"
	"errors"
	"strings"

	"cbe-super-app-cps-action/internal/constants"
	tac_dto "cbe-super-app-cps-action/internal/constants/dto/term_and_condition"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	tac_core "cbe-super-app-cps-action/internal/service/term_and_condition/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type accountOpeningTermsService struct {
	repo       storage.AccountOpeningTermsRepository
	cpsService service.CPSActionService
	logger     utils.Logger
	minio      *s3.Client
	bucketName string
	cfg        *config.VaultConfig
}

var _ service.AccountOpeningTermsService = (*accountOpeningTermsService)(nil)

func NewAccountOpeningTermsService(
	repo storage.AccountOpeningTermsRepository,
	cpsService service.CPSActionService,
	logger utils.Logger,
	minio *s3.Client,
	bucketName string,
	cfg *config.VaultConfig,
) service.AccountOpeningTermsService {
	return &accountOpeningTermsService{
		repo:       repo,
		cpsService: cpsService,
		logger:     logger,
		minio:      minio,
		bucketName: bucketName,
		cfg:        cfg,
	}
}

func (s *accountOpeningTermsService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "AccountOpeningTerms", "Authorize")
	defer span.End()

	tac, err := tac_core.MapFromAction(cpsAction.CurrentAction)
	if err != nil {
		log.Errorf("[TACSvc][Authorize] map err: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	switch constants.RequestAction(cpsAction.RequestAction) {
	case constants.RequestCreateTermAndCondition:
		if err := s.repo.Create(ctx, &tac); err != nil {
			span.AddEvent("create failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[TACSvc][Authorize] create err: %v", err)
			return nil, err
		}
		log.Infof("[TACSvc][Authorize] created")

	case constants.RequestDeleteTermAndCondition:
		if err := s.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			span.AddEvent("delete failed", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[TACSvc][Authorize] delete err: %v", err)
			return nil, err
		}
		log.Infof("[TACSvc][Authorize] deleted id=%s", cpsAction.UniqueId)

	default:
		log.Errorf("[TACSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	log.Infof("[TACSvc][Authorize] done action=%s", cpsAction.RequestAction)
	return cpsAction, nil
}

func (s *accountOpeningTermsService) GetAll(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]imodel.AccountOpeningTerms], error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAll", "AccountOpeningTerms", "GetAll")
	defer span.End()

	result, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("fetch failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[TACSvc][GetAll] err: %v", err)
		return nil, err
	}
	log.Infof("[TACSvc][GetAll] count=%d", len(result.Data))
	return result, nil
}

func (s *accountOpeningTermsService) GetByID(ctx context.Context, id string) (*imodel.AccountOpeningTerms, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "GetByID", "AccountOpeningTerms", "GetByID")
	defer span.End()

	result, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("fetch failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		log.Errorf("[TACSvc][GetByID] id=%s err=%v", id, err)
		return nil, err
	}
	log.Infof("[TACSvc][GetByID] found id=%s", id)
	return result, nil
}

func (s *accountOpeningTermsService) Upload(ctx context.Context, req tac_dto.CreateTACRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "Upload", "AccountOpeningTerms", "Upload")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		log.Errorf("[TACSvc][Upload] incomplete user")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByProductAndVersion(ctx, req.AccountProductID, req.VersionLabel)
	if err != nil && !strings.Contains(err.Error(), localization.ErrorResourceNotFound.Code) {
		log.Errorf("[TACSvc][Upload] dup check err: %v", err)
		return err
	}
	if existing != nil {
		log.Errorf("[TACSvc][Upload] duplicate product_id=%s version=%s", req.AccountProductID, req.VersionLabel)
		return errors.New(localization.ErrorTACAlreadyExists.Code)
	}

	pdfURL, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.TermAndCondition,
		string(constants.TermAndConditionFolderName), *s.cfg, "", s.logger)
	if err != nil {
		span.AddEvent("pdf upload failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[TACSvc][Upload] upload pdf err: %v", err)
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	payload := imodel.AccountOpeningTerms{
		AccountProductID:       req.AccountProductID,
		ActivationTime:         req.ActivationTime,
		VersionLabel:           req.VersionLabel,
		TermsAndConditionsPath: pdfURL,
		IsEnabled:              true,
	}

	action := lib.CpsModelBuilder("", makerData, nil, payload,
		string(constants.RequestCreateTermAndCondition), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("cps action failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[TACSvc][Upload] cps err: %v", err)
		return err
	}
	log.Infof("[TACSvc][Upload] request created product=%s version=%s", req.AccountProductID, req.VersionLabel)
	return nil
}

func (s *accountOpeningTermsService) Delete(ctx context.Context, id string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)
	ctx, span := local_util.TraceLogger(ctx, "service", "Delete", "AccountOpeningTerms", "Delete")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[TACSvc][Delete] find err: %v", err)
		return err
	}

	updated := *existing
	updated.IsDeleted = true

	action := lib.CpsModelBuilder(id, makerData, existing, updated,
		string(constants.RequestDeleteTermAndCondition), constants.DELETE)

	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		span.AddEvent("cps action failed", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[TACSvc][Delete] cps err: %v", err)
		return err
	}
	log.Infof("[TACSvc][Delete] request created id=%s", id)
	return nil
}
