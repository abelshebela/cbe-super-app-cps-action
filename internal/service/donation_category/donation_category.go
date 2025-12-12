package donation_category

import (
	"cbe-super-app-cps-action/internal/constants"
	dto "cbe-super-app-cps-action/internal/constants/dto/donation_category"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	core "cbe-super-app-cps-action/internal/service/donation_category/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"path"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DonationCategory struct {
	DonationCategoryRepo storage.DonationCategoryRepository
	cpsService           service.CPSActionService
	logger               utils.Logger
	minio                *s3.Client
	bucketName           string
	cfg                  *config.VaultConfig
	minioEndPoint        string
}

func NewDonationCategoryService(client *mongo.Client, DonationCategoryRepo storage.DonationCategoryRepository, cpsAction service.CPSActionService, logger utils.Logger, minio *s3.Client,
	bucketName string,
	cfg *config.VaultConfig,
	minioEndPoint string,
) service.DonationCategoryService {

	return &DonationCategory{
		DonationCategoryRepo: DonationCategoryRepo,
		cpsService:           cpsAction,
		logger:               logger,
		minio:                minio,
		bucketName:           bucketName,
		cfg:                  cfg,
		minioEndPoint:        minioEndPoint,
	}
}

func (d *DonationCategory) FetchDonationCategory(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]dto.DonationCategoryListResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchDonationCategory", "DonationCategory", "FetchDonationCategory")
	defer span.End()

	result, err := d.DonationCategoryRepo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch donation categories", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return result, nil
}

func (d *DonationCategory) FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchDonationCategoryByID", "DonationCategory", "FetchDonationCategoryByID")
	defer span.End()

	result, err := d.DonationCategoryRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to fetch donation category", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	return result, nil
}

func (d *DonationCategory) CreateDonationCategory(ctx context.Context, donationCategory dto.DonationCategoryRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateDonationCategory", "DonationCategory", "CreateDonationCategory")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
		))
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	ok, err := core.DonationNameExists(ctx, donationCategory.CategoryName, d.DonationCategoryRepo)
	if err != nil {
		span.AddEvent("Failed to check donation name existence", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("category_name", donationCategory.CategoryName),
		))
		return err
	}
	if ok {
		span.AddEvent("Donation category name duplicated", trace.WithAttributes(
			attribute.String("error", localization.ErrorDonationCategoryNameDuplicated.Code),
			attribute.String("category_name", donationCategory.CategoryName),
		))
		return errors.New(localization.ErrorDonationCategoryNameDuplicated.Code)
	}
	if donationCategory.Icon == nil {
		span.AddEvent("Icon required", trace.WithAttributes(
			attribute.String("error", localization.ErrorDonationCategoryIDRequired.Code),
		))
		return errors.New(localization.ErrorDonationCategoryIDRequired.Code)
	}
	url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donationCategory.Icon, string(constants.DonationIcon), *d.cfg, "", d.logger)
	if err != nil {
		d.logger.Errorf("failed to upload image to minio: %v", err)
		span.AddEvent("Failed to upload image to minio", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}
	result := model.DonationCategory{
		CategoryName:   donationCategory.CategoryName,
		Icon:           url,
		Enabled:        true,
		IsDeleted:      false,
		CreatedAt:      time.Now(),
		LastModifiedAt: time.Now(),
	}

	cpsAction := lib.CpsModelBuilder("", makerData, "", result, string(constants.RequestCreateDonationCategory), constants.CREATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}
	return nil
}

func (d *DonationCategory) UpdateDonationCategory(ctx context.Context, id string, donationCategory dto.DonationCategoryRequest) (dto.DonationCategoryRequest, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateDonationCategory", "DonationCategory", "UpdateDonationCategory")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
			attribute.String("id", id),
		))
		return donationCategory, errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	// Check if donation category exists
	existingCategory, err := d.DonationCategoryRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find donation category", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return donationCategory, err
	}
	if existingCategory == nil {
		span.AddEvent("Donation category not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
		))
		return donationCategory, errors.New(localization.ErrorFileNotFound.Code)
	}

	// Convert DTO to model for similarity check
	existingModel := &model.DonationCategory{
		CategoryName: existingCategory.CategoryName,
		Icon:         existingCategory.Icon,
		IsDeleted:    existingCategory.IsDeleted,
	}

	// Check if data is similar to existing data
	if core.IsDataSimilar(donationCategory, existingModel) {
		span.AddEvent("No changes to update", trace.WithAttributes(
			attribute.String("error", localization.ErrorNoChangesToUpdate.Code),
			attribute.String("id", id),
		))
		return donationCategory, errors.New(localization.ErrorNoChangesToUpdate.Code)
	}

	// Check if category name is being updated and if it already exists
	if donationCategory.CategoryName != "" && donationCategory.CategoryName != existingCategory.CategoryName {
		ok, err := core.DonationNameExists(ctx, donationCategory.CategoryName, d.DonationCategoryRepo)
		if err != nil {
			span.AddEvent("Failed to check donation name existence", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return donationCategory, err
		}
		if ok {
			span.AddEvent("Donation category name duplicated", trace.WithAttributes(
				attribute.String("error", localization.ErrorDonationCategoryNameDuplicated.Code),
				attribute.String("id", id),
				attribute.String("category_name", donationCategory.CategoryName),
			))
			return donationCategory, errors.New(localization.ErrorDonationCategoryNameDuplicated.Code)
		}
	}
	iconURL := existingCategory.Icon
	if donationCategory.Icon != nil {
		var objectkey string
		if existingCategory.Icon != "" {
			objectkey = path.Base(existingCategory.Icon)
		}

		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donationCategory.Icon, string(constants.DonationIcon), *d.cfg, objectkey, d.logger)
		if err != nil {
			span.AddEvent("Failed to upload image to minio", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return donationCategory, err
		}
		iconURL = url
	}
	donationCategoryName := existingCategory.CategoryName
	if donationCategory.CategoryName != "" {
		donationCategoryName = donationCategory.CategoryName
	}
	DonationCategory := core.MapToDonationCategory(donationCategoryName, iconURL, existingCategory.Enabled)

	cpsAction := lib.CpsModelBuilder(id, makerData, existingCategory, DonationCategory, string(constants.RequestUpdateDonationCategory), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return donationCategory, err
	}

	return dto.DonationCategoryRequest{
		CategoryName: DonationCategory.CategoryName,
	}, nil
}

func (d *DonationCategory) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "DonationCategory", "Authorize")
	defer span.End()

	// requestedAction := action.RequestAction
	if action.ActionStatus != constants.Approved {
		d.logger.Errorf("Tried to authorize service action without cps action approval")
		span.AddEvent("CPS action status invalid", trace.WithAttributes(
			attribute.String("error", localization.ErrorCPSActionStatusInvalid.Code),
			attribute.String("unique_id", action.UniqueId),
		))
		return nil, errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}
	var donationCPS *model.DonationCategory
	bindErr := core.BindAction(action.CurrentAction, &donationCPS)
	if bindErr != nil {
		d.logger.Errorf("failed to bind current action to donation Category: %v", bindErr)
		span.AddEvent("Failed to bind current action", trace.WithAttributes(
			attribute.String("error", bindErr.Error()),
			attribute.String("unique_id", action.UniqueId),
		))
		return nil, errors.New(localization.ErrorCPSActionFailed.Code)
	}

	switch action.RequestAction {
	case string(constants.RequestCreateDonationCategory):

		err := d.DonationCategoryRepo.Create(ctx, donationCPS)
		if err != nil {
			d.logger.Errorf("Failed to create donation category: %v", err)
			span.AddEvent("Failed to create donation category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestUpdateDonationCategory):
		err := d.DonationCategoryRepo.Update(ctx, action.UniqueId, donationCPS)
		if err != nil {
			d.logger.Errorf("Failed to update donation category: %v", err)
			span.AddEvent("Failed to update donation category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestEnableDonationCategory):
		err := d.DonationCategoryRepo.Update(ctx, action.UniqueId, donationCPS)
		if err != nil {
			d.logger.Errorf("Failed to update donation category: %v", err)
			span.AddEvent("Failed to enable donation category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDisableDonationCategory):
		err := d.DonationCategoryRepo.Update(ctx, action.UniqueId, donationCPS)
		if err != nil {
			d.logger.Errorf("Failed to update donation category: %v", err)
			span.AddEvent("Failed to disable donation category", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	default:
		d.logger.Errorf("Unsupported action requested: %s", action.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(action.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	d.logger.Infof("Service action authorization completed: %s", action.RequestAction)
	return action, nil

}

func (d *DonationCategory) EnableDonationCategory(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableDonationCategory", "DonationCategory", "EnableDonationCategory")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	existingDonationCategory, err := d.DonationCategoryRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find donation category", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if existingDonationCategory == nil {
		span.AddEvent("Donation category not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorFileNotFound.Code)
	}
	if existingDonationCategory.Enabled {
		span.AddEvent("Donation category already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorAlreadyEnabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorAlreadyEnabled.Code)
	}
	Enabled := true
	DonationCategory := core.MapToDonationCategory(existingDonationCategory.CategoryName, existingDonationCategory.Icon, Enabled)
	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonationCategory, DonationCategory, string(constants.RequestEnableDonationCategory), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil

}

func (d *DonationCategory) DisableDonationCategory(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableDonationCategory", "DonationCategory", "DisableDonationCategory")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	existingDonationCategory, err := d.DonationCategoryRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find donation category", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if existingDonationCategory == nil {
		span.AddEvent("Donation category not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorFileNotFound.Code)
	}
	if !existingDonationCategory.Enabled {
		span.AddEvent("Donation category already disabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorAlreadyDisabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorAlreadyDisabled.Code)
	}
	Enabled := false
	DonationCategory := core.MapToDonationCategory(existingDonationCategory.CategoryName, existingDonationCategory.Icon, Enabled)
	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonationCategory, DonationCategory, string(constants.RequestDisableDonationCategory), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil

}
